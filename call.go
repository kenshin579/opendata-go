package opendata

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// Fetch 는 게이트웨이에 GET 요청을 보내고 response/body/items/item 을 []T 로 디코드한다.
// path 는 "/1220000/nationtrade/getNationtradeList" 처럼 기관코드부터 쓴다. q 에 serviceKey 는 넣지 않는다.
//
// 결과가 0건이면(<items/>, <body/>, body 없음 등) 빈 슬라이스와 nil 에러를 돌려준다.
// 게이트웨이 차단은 *GatewayError, resultCode ≠ "00" 은 *APIError 다.
// 반환하는 어떤 에러에도 serviceKey 는 담기지 않는다.
func Fetch[T any](ctx context.Context, c *Client, path string, q url.Values) ([]T, error) {
	body, status, err := c.get(ctx, path, q)
	if err != nil {
		return nil, err
	}
	return decode[T](body, status, path)
}

func (c *Client) get(ctx context.Context, path string, q url.Values) ([]byte, int, error) {
	// serviceKey 를 먼저 두고 나머지는 Encode 로 붙인다 (_type 은 붙이지 않는다 — 게이트웨이 에러 04).
	raw := "serviceKey=" + url.QueryEscape(c.key)
	if enc := q.Encode(); enc != "" {
		raw += "&" + enc
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+path+"?"+raw, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("opendata: %s: %s", path, c.mask(err))
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("opendata: GET %s: %s", path, c.mask(err))
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("opendata: GET %s: %s", path, c.mask(err))
	}
	return body, resp.StatusCode, nil
}

// mask 는 에러 문자열에서 serviceKey(원문·URL 인코딩)를 지운다.
func (c *Client) mask(err error) string {
	s := err.Error()
	for _, k := range []string{c.key, url.QueryEscape(c.key)} {
		if k != "" {
			s = strings.ReplaceAll(s, k, "{SERVICE_KEY}")
		}
	}
	return s
}

type gatewayEnvelope struct {
	XMLName xml.Name `xml:"OpenAPI_ServiceResponse"`
	Header  struct {
		ErrMsg     string `xml:"errMsg"`
		AuthMsg    string `xml:"returnAuthMsg"`
		ReasonCode string `xml:"returnReasonCode"`
	} `xml:"cmmMsgHeader"`
}

type responseEnvelope[T any] struct {
	XMLName xml.Name `xml:"response"`
	Header  struct {
		ResultCode string `xml:"resultCode"`
		ResultMsg  string `xml:"resultMsg"`
	} `xml:"header"`
	Items []T `xml:"body>items>item"`
}

func decode[T any](body []byte, status int, path string) ([]T, error) {
	switch rootName(body) {
	case "OpenAPI_ServiceResponse":
		var g gatewayEnvelope
		if err := xml.Unmarshal(body, &g); err != nil {
			return nil, fmt.Errorf("opendata: decode gateway error %s: %w", path, err)
		}
		return nil, &GatewayError{HTTPStatus: status, Code: g.Header.ReasonCode, ErrMsg: g.Header.ErrMsg, AuthMsg: g.Header.AuthMsg}
	case "response":
		var r responseEnvelope[T]
		if err := xml.Unmarshal(body, &r); err != nil {
			return nil, fmt.Errorf("opendata: decode %s: %w", path, err)
		}
		if r.Header.ResultCode != "00" {
			return nil, &APIError{Code: r.Header.ResultCode, Message: r.Header.ResultMsg}
		}
		if r.Items == nil {
			r.Items = []T{}
		}
		return r.Items, nil
	default:
		return nil, fmt.Errorf("opendata: GET %s: unexpected response (http %d): %q", path, status, head(body, 200))
	}
}

// rootName 은 XML 문서의 첫 요소 이름을 돌려준다. XML 이 아니면 "".
func rootName(body []byte) string {
	d := xml.NewDecoder(bytes.NewReader(body))
	for {
		tok, err := d.Token()
		if err != nil {
			return ""
		}
		if se, ok := tok.(xml.StartElement); ok {
			return se.Name.Local
		}
	}
}

func head(b []byte, n int) string {
	if len(b) > n {
		b = b[:n]
	}
	return string(b)
}
