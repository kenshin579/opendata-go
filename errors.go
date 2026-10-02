package opendata

import "fmt"

// GatewayError 는 포털 게이트웨이가 요청을 막았을 때의 에러다 (OpenAPI_ServiceResponse).
// 예: 서비스키 누락(20, HTTP 401), 미등록 키·활용신청 안 한 API(30, HTTP 403), 없는 오퍼레이션(12, HTTP 400),
// 트래픽 초과(22). 코드표는 docs/api/README.md.
type GatewayError struct {
	HTTPStatus int
	Code       string // returnReasonCode
	ErrMsg     string // errMsg (예: SERVICE_KEY_IS_NOT_REGISTERED_ERROR)
	AuthMsg    string // returnAuthMsg (예: 등록되지 않은 서비스키)
}

func (e *GatewayError) Error() string {
	return fmt.Sprintf("opendata: gateway error %s %s (%s, http %d)", e.Code, e.ErrMsg, e.AuthMsg, e.HTTPStatus)
}

// APIError 는 기관 서비스가 resultCode ≠ "00" 으로 응답했을 때의 에러다.
// 관세청은 모든 실패가 99 이고 원인은 Message 로만 구분된다 (HTTP 200).
type APIError struct {
	Code    string // resultCode
	Message string // resultMsg
}

func (e *APIError) Error() string {
	return fmt.Sprintf("opendata: api error %s: %s", e.Code, e.Message)
}
