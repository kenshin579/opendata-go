// Package opendata 는 공공데이터포털(data.go.kr) OpenAPI 의 Go 클라이언트다.
// 포털 공통(서비스키·게이트웨이 호출·XML 봉투·에러)을 맡고, 기관별 API 는 서브패키지(customs 등)가 제공한다.
package opendata

import (
	"errors"
	"net/http"
	"os"
	"time"
)

// DefaultBaseURL 은 공공데이터포털 API 게이트웨이 주소다.
const DefaultBaseURL = "https://apis.data.go.kr"

// EnvServiceKey 는 NewClientFromEnv 가 읽는 환경변수 이름이다.
const EnvServiceKey = "OPENDATA_API_KEY"

// Client 는 포털 게이트웨이 호출 통로다. 동시에 여러 고루틴에서 써도 된다.
type Client struct {
	key     string
	baseURL string
	http    *http.Client
}

// NewClient 는 서비스키로 Client 를 만든다.
func NewClient(serviceKey string, opts ...Option) (*Client, error) {
	if serviceKey == "" {
		return nil, errors.New("opendata: serviceKey is required")
	}
	cfg := clientOptions{baseURL: DefaultBaseURL, timeout: 30 * time.Second}
	for _, opt := range opts {
		opt(&cfg)
	}
	hc := cfg.httpClient
	if hc == nil {
		hc = &http.Client{Timeout: cfg.timeout}
	}
	return &Client{key: serviceKey, baseURL: cfg.baseURL, http: hc}, nil
}

// NewClientFromEnv 는 OPENDATA_API_KEY 환경변수의 서비스키로 Client 를 만든다.
func NewClientFromEnv(opts ...Option) (*Client, error) {
	return NewClient(os.Getenv(EnvServiceKey), opts...)
}
