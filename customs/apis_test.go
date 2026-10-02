package customs

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/kenshin579/opendata-go"
)

// fixtures 는 오퍼레이션 경로 첫 세그먼트 → testdata/customs 파일 (1단계 실측 응답 축약본: 일반 행 3개 + 총계 행).
var fixtures = map[string]string{
	"nitemtrade": "01_nitemtrade", "Itemtrade": "02_itemtrade", "nationtrade": "03_nationtrade",
	"continenttradet": "04_continent", "economytrade": "05_economy", "Idfytempertrade": "06_idfytemper",
	"ntempertrade": "07_ntemper", "newtempertrade": "08_newtemper", "nnewtempertrade": "09_nnewtemper",
	"kindtrade": "10_kind", "customstrade": "11_customs", "porttrade": "12_port", "sidotrade": "13_sido",
	"sidoitemtrade": "14_sidoitem", "sidotempertrade": "15_sidotemper", "sigunguperimexacrs": "16_sigungu",
	"sigunguperprlstperacrs": "17_sigunguitem",
}

func newTestService(t *testing.T, last *url.URL) *Service {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		*last = *r.URL
		seg := strings.Split(strings.TrimPrefix(r.URL.Path, basePath), "/")[0]
		b, err := os.ReadFile("../testdata/customs/" + fixtures[seg] + ".xml")
		if err != nil {
			t.Errorf("no fixture for %s", r.URL.Path)
		}
		w.Write(b)
	}))
	t.Cleanup(srv.Close)
	c, err := opendata.NewClient("k", opendata.WithBaseURL(srv.URL))
	if err != nil {
		t.Fatal(err)
	}
	return New(c)
}

// check 는 행 수·총계 수와, 첫 일반 행의 모든 필드가 채워졌는지(= xml 태그가 실응답과 맞는지) 본다.
func check[R totaler](t *testing.T, name string, rows []R, err error, wantRows, wantTotals int) {
	t.Helper()
	if err != nil {
		t.Fatalf("%s: %v", name, err)
	}
	totals := len(rows) - len(WithoutTotal(rows))
	if len(rows) != wantRows || totals != wantTotals {
		t.Fatalf("%s: rows=%d totals=%d; want %d, %d", name, len(rows), totals, wantRows, wantTotals)
	}
	first := WithoutTotal(rows)[0]
	v := reflect.ValueOf(first)
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).String() == "" {
			t.Errorf("%s: field %s is empty in first row — xml tag mismatch?", name, v.Type().Field(i).Name)
		}
	}
}

func TestAPIs(t *testing.T) {
	ctx := context.Background()
	var u url.URL
	s := newTestService(t, &u)
	const a, b = "202601", "202606"

	r1, err := s.ItemCountryTrade(ctx, ItemCountryTradeParams{Start: a, End: b, HsSgn: "8542", CntyCd: "CN"})
	check(t, "ItemCountryTrade", r1, err, 4, 1)
	if q := u.Query(); u.Path != "/1220000/nitemtrade/getNitemtradeList" || q.Get("hsSgn") != "8542" || q.Get("cntyCd") != "CN" || q.Get("strtYymm") != a || q.Get("endYymm") != b {
		t.Errorf("ItemCountryTrade request = %s?%s", u.Path, u.RawQuery)
	}
	r2, err := s.ItemTrade(ctx, ItemTradeParams{Start: a, End: b})
	check(t, "ItemTrade", r2, err, 4, 1)
	if _, ok := u.Query()["hsSgn"]; ok {
		t.Error("empty HsSgn must not be sent")
	}
	r3, err := s.CountryTrade(ctx, CountryTradeParams{Start: a, End: b, CntyCd: "US"})
	check(t, "CountryTrade", r3, err, 3, 0)
	r4, err := s.ContinentTrade(ctx, ContinentTradeParams{Start: a, End: b})
	check(t, "ContinentTrade", r4, err, 3, 0)
	r5, err := s.EconomicBlocTrade(ctx, EconomicBlocTradeParams{Start: a, End: b})
	check(t, "EconomicBlocTrade", r5, err, 3, 0)
	r6, err := s.PropertyTrade(ctx, PropertyTradeParams{Start: a, End: b, ImexTpcd: "1"})
	check(t, "PropertyTrade", r6, err, 3, 0)
	r7, err := s.PropertyCountryTrade(ctx, PropertyCountryTradeParams{Start: a, End: b, ImexTpcd: "1", CntyCd: "US"})
	check(t, "PropertyCountryTrade", r7, err, 3, 0)
	r8, err := s.NewPropertyTrade(ctx, NewPropertyTradeParams{Start: a, End: b, ImexTpcd: "1"})
	check(t, "NewPropertyTrade", r8, err, 3, 0)
	r9, err := s.NewPropertyCountryTrade(ctx, NewPropertyCountryTradeParams{Start: a, End: b, ImexTpcd: "1", ImexTmprUnfcClsfCd: "13050102", CntyCd: "US"})
	check(t, "NewPropertyCountryTrade", r9, err, 3, 0)
	r10, err := s.KindTrade(ctx, KindTradeParams{Start: a, End: b, ImexTpcd: "1"})
	check(t, "KindTrade", r10, err, 4, 1)
	r11, err := s.CustomsOfficeTrade(ctx, CustomsOfficeTradeParams{Start: a, End: b})
	check(t, "CustomsOfficeTrade", r11, err, 3, 0)
	r12, err := s.PortTrade(ctx, PortTradeParams{Start: a, End: b})
	check(t, "PortTrade", r12, err, 4, 1)
	r13, err := s.SidoTrade(ctx, SidoTradeParams{Start: a, End: b})
	check(t, "SidoTrade", r13, err, 4, 1)
	r14, err := s.SidoItemTrade(ctx, SidoItemTradeParams{Start: a, End: b, SidoCd: "41"})
	check(t, "SidoItemTrade", r14, err, 4, 1)
	r15, err := s.SidoPropertyTrade(ctx, SidoPropertyTradeParams{Start: a, End: b, SidoCd: "41", ImexTpcd: "1"})
	check(t, "SidoPropertyTrade", r15, err, 4, 1)
	r16, err := s.SigunguTrade(ctx, SigunguTradeParams{Start: a, End: b, SidoCd: "41"})
	check(t, "SigunguTrade", r16, err, 3, 0)
	r17, err := s.SigunguItemTrade(ctx, SigunguItemTradeParams{Start: a, End: b, SidoCd: "41", HsSgn: "854232"})
	check(t, "SigunguItemTrade", r17, err, 3, 0)
	if q := u.Query(); q.Get("HsSgn") != "854232" {
		t.Errorf("SigunguItemTrade must send uppercase HsSgn: %s", u.RawQuery)
	}

	// 천 달러·공백 패딩 숫자 파싱
	v, err := r13[1].ExpUsdAmt.Int64()
	if err != nil || v <= 0 {
		t.Errorf("SidoTrade ExpUsdAmt %q → %d, %v", r13[1].ExpUsdAmt, v, err)
	}
}

func TestAPIRejectsLongPeriod(t *testing.T) {
	var u url.URL
	s := newTestService(t, &u)
	if _, err := s.CountryTrade(context.Background(), CountryTradeParams{Start: "202506", End: "202606"}); err == nil {
		t.Fatal("want error for 13 months")
	}
	if u.Path != "" {
		t.Fatal("must not call the server")
	}
}
