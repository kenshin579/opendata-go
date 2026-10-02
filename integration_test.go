//go:build integration

package opendata_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/kenshin579/opendata-go"
	"github.com/kenshin579/opendata-go/customs"
)

// 실 API 호출: OPENDATA_API_KEY 필요, 17개 API 모두 활용신청돼 있어야 한다.
//
//	go test -tags integration ./...
func TestCustomsLive(t *testing.T) {
	c, err := opendata.NewClientFromEnv(opendata.WithTimeout(60 * time.Second))
	if err != nil {
		t.Skip("OPENDATA_API_KEY not set")
	}
	s := customs.New(c)
	ctx := context.Background()
	const m = "202601"
	calls := []struct {
		name    string
		allow0  bool
		rowsErr func() (int, error)
	}{
		{"ItemCountryTrade", false, func() (int, error) {
			r, e := s.ItemCountryTrade(ctx, customs.ItemCountryTradeParams{Start: m, End: m, HsSgn: "8542", CntyCd: "CN"})
			return len(r), e
		}},
		{"ItemTrade", false, func() (int, error) {
			r, e := s.ItemTrade(ctx, customs.ItemTradeParams{Start: m, End: m, HsSgn: "8542"})
			return len(r), e
		}},
		{"CountryTrade", false, func() (int, error) {
			r, e := s.CountryTrade(ctx, customs.CountryTradeParams{Start: m, End: m, CntyCd: "US"})
			return len(r), e
		}},
		{"ContinentTrade", false, func() (int, error) {
			r, e := s.ContinentTrade(ctx, customs.ContinentTradeParams{Start: m, End: m})
			return len(r), e
		}},
		{"EconomicBlocTrade", false, func() (int, error) {
			r, e := s.EconomicBlocTrade(ctx, customs.EconomicBlocTradeParams{Start: m, End: m})
			return len(r), e
		}},
		{"PropertyTrade", false, func() (int, error) {
			r, e := s.PropertyTrade(ctx, customs.PropertyTradeParams{Start: m, End: m, ImexTpcd: "1"})
			return len(r), e
		}},
		{"PropertyCountryTrade", false, func() (int, error) {
			r, e := s.PropertyCountryTrade(ctx, customs.PropertyCountryTradeParams{Start: m, End: m, ImexTpcd: "1", CntyCd: "US"})
			return len(r), e
		}},
		{"NewPropertyTrade", false, func() (int, error) {
			r, e := s.NewPropertyTrade(ctx, customs.NewPropertyTradeParams{Start: m, End: m, ImexTpcd: "1", ImexTmprUnfcClsfCd: "13050102"})
			return len(r), e
		}},
		{"NewPropertyCountryTrade", true, func() (int, error) {
			r, e := s.NewPropertyCountryTrade(ctx, customs.NewPropertyCountryTradeParams{Start: m, End: m, ImexTpcd: "1", ImexTmprUnfcClsfCd: "13050102", CntyCd: "US"})
			return len(r), e
		}},
		{"KindTrade", false, func() (int, error) {
			r, e := s.KindTrade(ctx, customs.KindTradeParams{Start: m, End: m, ImexTpcd: "1"})
			return len(r), e
		}},
		{"CustomsOfficeTrade", false, func() (int, error) {
			r, e := s.CustomsOfficeTrade(ctx, customs.CustomsOfficeTradeParams{Start: m, End: m})
			return len(r), e
		}},
		{"PortTrade", false, func() (int, error) {
			r, e := s.PortTrade(ctx, customs.PortTradeParams{Start: m, End: m, PortAirptRegnCd: "010"})
			return len(r), e
		}},
		{"SidoTrade", false, func() (int, error) {
			r, e := s.SidoTrade(ctx, customs.SidoTradeParams{Start: m, End: m})
			return len(r), e
		}},
		{"SidoItemTrade", false, func() (int, error) {
			r, e := s.SidoItemTrade(ctx, customs.SidoItemTradeParams{Start: m, End: m, SidoCd: "41"})
			return len(r), e
		}},
		{"SidoPropertyTrade", false, func() (int, error) {
			r, e := s.SidoPropertyTrade(ctx, customs.SidoPropertyTradeParams{Start: m, End: m, SidoCd: "41", ImexTpcd: "1"})
			return len(r), e
		}},
		{"SigunguTrade", false, func() (int, error) {
			r, e := s.SigunguTrade(ctx, customs.SigunguTradeParams{Start: m, End: m, SidoCd: "41"})
			return len(r), e
		}},
		{"SigunguItemTrade", false, func() (int, error) {
			r, e := s.SigunguItemTrade(ctx, customs.SigunguItemTradeParams{Start: m, End: m, SidoCd: "41", HsSgn: "854232"})
			return len(r), e
		}},
	}
	for _, c := range calls {
		n, err := c.rowsErr()
		if err != nil {
			t.Errorf("%s: %v", c.name, err)
			continue
		}
		if n == 0 && !c.allow0 {
			t.Errorf("%s: 0 rows", c.name)
		}
		t.Logf("%s: %d rows", c.name, n)
	}
}

func TestGatewayErrorLive(t *testing.T) {
	c, _ := opendata.NewClient("not-a-registered-key")
	_, err := customs.New(c).CountryTrade(context.Background(), customs.CountryTradeParams{Start: "202601", End: "202601"})
	var ge *opendata.GatewayError
	if !errors.As(err, &ge) || ge.Code != "30" {
		t.Fatalf("want gateway 30, got %v", err)
	}
}
