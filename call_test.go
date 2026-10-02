package opendata

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"
)

const testKey = "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"

type row struct {
	Year   string `xml:"year"`
	ExpDlr Num    `xml:"expDlr"`
}

// serve 는 status 와 body 로 응답하고 마지막 요청을 *got 에 담는 서버를 띄운다.
func serve(t *testing.T, status int, body string, got **http.Request) *Client {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got != nil {
			*got = r
		}
		w.WriteHeader(status)
		w.Write([]byte(body))
	}))
	t.Cleanup(srv.Close)
	c, err := NewClient(testKey, WithBaseURL(srv.URL))
	if err != nil {
		t.Fatal(err)
	}
	return c
}

func fixture(t *testing.T, name string) string {
	t.Helper()
	b, err := os.ReadFile("testdata/portal/" + name)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func TestFetchQuery(t *testing.T) {
	var r *http.Request
	c := serve(t, 200, fixture(t, "single_item.xml"), &r)
	q := url.Values{"strtYymm": {"202601"}, "endYymm": {"202601"}}
	if _, err := Fetch[row](context.Background(), c, "/1220000/nationtrade/getNationtradeList", q); err != nil {
		t.Fatal(err)
	}
	if r.URL.Path != "/1220000/nationtrade/getNationtradeList" {
		t.Errorf("path = %s", r.URL.Path)
	}
	got := r.URL.Query()
	if got.Get("serviceKey") != testKey || got.Get("strtYymm") != "202601" {
		t.Errorf("query = %v", got)
	}
	if _, ok := got["_type"]; ok {
		t.Error("_type must never be sent (gateway error 04)")
	}
}

func TestFetchItems(t *testing.T) {
	c := serve(t, 200, fixture(t, "single_item.xml"), nil)
	rows, err := Fetch[row](context.Background(), c, "/x", nil)
	if err != nil || len(rows) != 1 || rows[0].Year == "" {
		t.Fatalf("rows=%+v err=%v", rows, err)
	}
}

func TestFetchEmpty(t *testing.T) {
	bodies := map[string]string{
		"items self-closing": fixture(t, "empty_items.xml"),
		"body self-closing":  `<response><header><resultCode>00</resultCode><resultMsg>OK</resultMsg></header><body/></response>`,
		"no body":            `<response><header><resultCode>00</resultCode><resultMsg>OK</resultMsg></header></response>`,
		"totalCount only":    `<response><header><resultCode>00</resultCode><resultMsg>OK</resultMsg></header><body><totalCount>0</totalCount></body></response>`,
	}
	for name, body := range bodies {
		c := serve(t, 200, body, nil)
		rows, err := Fetch[row](context.Background(), c, "/x", nil)
		if err != nil || rows == nil || len(rows) != 0 {
			t.Errorf("%s: rows=%v (nil=%v) err=%v; want empty non-nil slice", name, rows, rows == nil, err)
		}
	}
}

func TestFetchGatewayError(t *testing.T) {
	cases := []struct {
		file   string
		status int
		code   string
	}{
		{"gateway_30_403.xml", 403, "30"},
		{"gateway_20_401.xml", 401, "20"},
		{"gateway_12_400.xml", 400, "12"},
	}
	for _, tc := range cases {
		c := serve(t, tc.status, fixture(t, tc.file), nil)
		_, err := Fetch[row](context.Background(), c, "/x", nil)
		var ge *GatewayError
		if !errors.As(err, &ge) || ge.Code != tc.code || ge.HTTPStatus != tc.status || ge.ErrMsg == "" {
			t.Errorf("%s: err=%v", tc.file, err)
		}
	}
}

func TestFetchAPIError(t *testing.T) {
	for _, f := range []string{"api_99_nobody.xml", "api_99_emptybody.xml", "api_99_totalcount.xml"} {
		c := serve(t, 200, fixture(t, f), nil)
		_, err := Fetch[row](context.Background(), c, "/x", nil)
		var ae *APIError
		if !errors.As(err, &ae) || ae.Code != "99" || ae.Message == "" {
			t.Errorf("%s: err=%v", f, err)
		}
	}
}

func TestFetchUnexpected(t *testing.T) {
	c := serve(t, 502, "<html><body>Bad Gateway</body></html>", nil)
	_, err := Fetch[row](context.Background(), c, "/x", nil)
	if err == nil || !strings.Contains(err.Error(), "502") {
		t.Fatalf("err=%v", err)
	}
	c = serve(t, 200, "", nil)
	if _, err := Fetch[row](context.Background(), c, "/x", nil); err == nil {
		t.Fatal("want error for empty body")
	}
}

func TestFetchMasksKey(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
	}))
	defer srv.Close()
	c, _ := NewClient(testKey, WithBaseURL(srv.URL), WithTimeout(20*time.Millisecond))
	_, err := Fetch[row](context.Background(), c, "/x", nil)
	if err == nil {
		t.Fatal("want timeout error")
	}
	if strings.Contains(err.Error(), testKey) || !strings.Contains(err.Error(), "{SERVICE_KEY}") {
		t.Fatalf("key not masked: %v", err)
	}
}
