// Package customs 는 관세청(기관코드 1220000) 수출입실적 OpenAPI 17개를 제공한다.
// 문서: docs/api/customs/README.md. 금액 단위는 API 마다 다르다 — 시도·시군구 계열(Sido*, Sigungu*)은 천 달러.
package customs

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/kenshin579/opendata-go"
)

const (
	basePath    = "/1220000/"
	totalPeriod = "총계"
	maxMonths   = 12
)

// Service 는 관세청 수출입실적 API 묶음이다.
type Service struct{ c *opendata.Client }

// New 는 opendata.Client 로 Service 를 만든다.
func New(c *opendata.Client) *Service { return &Service{c: c} }

// Period 는 조회 기간(YYYYMM~YYYYMM)이다.
type Period struct{ Start, End string }

// totaler 는 IsTotal 을 가진 Row 타입이다.
type totaler interface{ IsTotal() bool }

// WithoutTotal 은 총계 행을 뺀 새 슬라이스를 돌려준다.
func WithoutTotal[R totaler](rows []R) []R {
	out := make([]R, 0, len(rows))
	for _, r := range rows {
		if !r.IsTotal() {
			out = append(out, r)
		}
	}
	return out
}

// Years 는 [start, end] 를 달력 연도 경계로 나눈다 (예: 202407~202606 → 3개 창).
// 한 번 호출로 12개월까지만 조회되므로 긴 기간은 창마다 호출해 이어 붙인다.
// 연 단위로 합산하는 Sido* API 가 연도 중간에서 잘리지 않게 연도 경계로 자른다.
func Years(start, end string) ([]Period, error) {
	sy, sm, err := parseYM(start)
	if err != nil {
		return nil, err
	}
	ey, em, err := parseYM(end)
	if err != nil {
		return nil, err
	}
	if sy*12+sm > ey*12+em {
		return nil, fmt.Errorf("customs: start %s is after end %s", start, end)
	}
	var out []Period
	for y := sy; y <= ey; y++ {
		from, to := 1, 12
		if y == sy {
			from = sm
		}
		if y == ey {
			to = em
		}
		out = append(out, Period{Start: fmt.Sprintf("%04d%02d", y, from), End: fmt.Sprintf("%04d%02d", y, to)})
	}
	return out, nil
}

// periodQuery 는 기간을 검증하고 strtYymm/endYymm 쿼리를 만든다.
func periodQuery(start, end string) (url.Values, error) {
	sy, sm, err := parseYM(start)
	if err != nil {
		return nil, err
	}
	ey, em, err := parseYM(end)
	if err != nil {
		return nil, err
	}
	n := (ey*12 + em) - (sy*12 + sm) + 1
	if n < 1 {
		return nil, fmt.Errorf("customs: start %s is after end %s", start, end)
	}
	if n > maxMonths {
		return nil, fmt.Errorf("customs: period %s~%s is %d months; at most %d (use customs.Years)", start, end, n, maxMonths)
	}
	return url.Values{"strtYymm": {start}, "endYymm": {end}}, nil
}

func parseYM(s string) (year, month int, err error) {
	if len(s) != 6 {
		return 0, 0, fmt.Errorf("customs: %q is not YYYYMM", s)
	}
	y, err1 := strconv.Atoi(s[:4])
	m, err2 := strconv.Atoi(s[4:])
	if err1 != nil || err2 != nil || m < 1 || m > 12 {
		return 0, 0, fmt.Errorf("customs: %q is not YYYYMM", s)
	}
	return y, m, nil
}

// set 은 빈 값이 아니면 쿼리에 넣는다.
func set(q url.Values, key, v string) {
	if v != "" {
		q.Set(key, v)
	}
}
