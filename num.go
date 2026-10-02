package opendata

import (
	"errors"
	"strconv"
	"strings"
)

// ErrNotNumber 는 Num 이 숫자가 아닐 때(총계 행의 "-", 빈 값 등) 반환된다.
var ErrNotNumber = errors.New("opendata: not a number")

// Num 은 응답의 숫자 값을 원문 그대로 담는다.
// 관세청 시도·시군구 계열은 "          -2,218" 처럼 공백 패딩과 천 단위 콤마가 붙어 온다.
type Num string

// String 은 원문을 돌려준다.
func (n Num) String() string { return string(n) }

// Int64 는 앞뒤 공백과 천 단위 콤마를 지우고 정수로 파싱한다.
func (n Num) Int64() (int64, error) {
	s := strings.ReplaceAll(strings.TrimSpace(string(n)), ",", "")
	if s == "" || s == "-" {
		return 0, ErrNotNumber
	}
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, ErrNotNumber
	}
	return v, nil
}
