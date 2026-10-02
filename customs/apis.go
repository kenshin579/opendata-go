// Code generated from docs/api/customs/*.md tables; review before editing. DO NOT EDIT by hand without updating docs.

package customs

import (
	"context"

	"github.com/kenshin579/opendata-go"
)

// ItemCountryTradeParams 는 품목별 국가별 수출입실적 요청 인자다. 빈 문자열 인자는 보내지 않는다.
type ItemCountryTradeParams struct {
	Start  string // 시작년월 YYYYMM
	End    string // 종료년월 YYYYMM (Start 포함 최대 12개월)
	HsSgn  string // hsSgn: 품목코드 (N). HS 2·4·6·10자리. 그 밖의 자릿수는 99 에러. 자릿수가 출력 행 단위를 정한다(함정). 코드는 codes/품목코드.csv
	CntyCd string // cntyCd: 국가코드 (명세 Y / 실측 생략 가능). 대문자 2자리(CN). 생략하면 전 국가. 코드는 codes/국가코드.csv
}

// ItemCountryTradeRow 는 품목별 국가별 수출입실적 응답 행이다. 문서: docs/api/customs/품목별국가별.md
type ItemCountryTradeRow struct {
	Year           string       `xml:"year"`           // 기간. 2026.01 형식. 총계 행은 총계
	StatCd         string       `xml:"statCd"`         // 국가코드. CN. 총계 행은 -
	StatCdCntnKor1 string       `xml:"statCdCntnKor1"` // 국가명. 중국. 총계 행은 -
	HsCd           string       `xml:"hsCd"`           // HS코드. 문자열. 자릿수는 요청 hsSgn 에 따라 4·6·10자리. 총계 행은 -
	StatKor        string       `xml:"statKor"`        // 품목명. 총계 행은 -
	ExpWgt         opendata.Num `xml:"expWgt"`         // 수출중량 (kg). 정수
	ExpDlr         opendata.Num `xml:"expDlr"`         // 수출금액 (달러). 정수
	ImpWgt         opendata.Num `xml:"impWgt"`         // 수입중량 (kg). 정수
	ImpDlr         opendata.Num `xml:"impDlr"`         // 수입금액 (달러). 정수
	BalPayments    opendata.Num `xml:"balPayments"`    // 무역수지 (달러). 정수, 음수 가능
}

// IsTotal 은 기간 필드가 "총계" 인 합계 행인지 알려 준다.
func (r ItemCountryTradeRow) IsTotal() bool { return r.Year == totalPeriod }

// ItemCountryTrade 는 품목별 국가별 수출입실적을 조회한다 (nitemtrade/getNitemtradeList).
func (s *Service) ItemCountryTrade(ctx context.Context, p ItemCountryTradeParams) ([]ItemCountryTradeRow, error) {
	q, err := periodQuery(p.Start, p.End)
	if err != nil {
		return nil, err
	}
	set(q, "hsSgn", p.HsSgn)
	set(q, "cntyCd", p.CntyCd)
	return opendata.Fetch[ItemCountryTradeRow](ctx, s.c, basePath+"nitemtrade/getNitemtradeList", q)
}

// ItemTradeParams 는 품목별 수출입실적 요청 인자다. 빈 문자열 인자는 보내지 않는다.
type ItemTradeParams struct {
	Start string // 시작년월 YYYYMM
	End   string // 종료년월 YYYYMM (Start 포함 최대 12개월)
	HsSgn string // hsSgn: 품목코드 (N). HS 코드 앞자리(2·4·6·10자리 실측). 출력은 항상 10자리 행. 명세 타입은 number 지만 0101 처럼 0 으로 시작하므로 문자열로 보낸다. 코드는 codes/품목코드.csv
}

// ItemTradeRow 는 품목별 수출입실적 응답 행이다. 문서: docs/api/customs/품목별.md
type ItemTradeRow struct {
	Year        string       `xml:"year"`        // 기간. 2026.01 형식. 총계 행은 총계
	HsCode      string       `xml:"hsCode"`      // HS코드. 10자리 문자열(명세 타입은 number). 총계 행은 -
	StatKor     string       `xml:"statKor"`     // 품목명. 10자리 품목명. 총계 행은 -
	ExpWgt      opendata.Num `xml:"expWgt"`      // 수출중량 (kg). 정수
	ExpDlr      opendata.Num `xml:"expDlr"`      // 수출금액 (달러). 정수
	ImpWgt      opendata.Num `xml:"impWgt"`      // 수입중량 (kg). 정수
	ImpDlr      opendata.Num `xml:"impDlr"`      // 수입금액 (달러). 정수
	BalPayments opendata.Num `xml:"balPayments"` // 무역수지 (달러). 정수, 음수 가능
}

// IsTotal 은 기간 필드가 "총계" 인 합계 행인지 알려 준다.
func (r ItemTradeRow) IsTotal() bool { return r.Year == totalPeriod }

// ItemTrade 는 품목별 수출입실적을 조회한다 (Itemtrade/getItemtradeList).
func (s *Service) ItemTrade(ctx context.Context, p ItemTradeParams) ([]ItemTradeRow, error) {
	q, err := periodQuery(p.Start, p.End)
	if err != nil {
		return nil, err
	}
	set(q, "hsSgn", p.HsSgn)
	return opendata.Fetch[ItemTradeRow](ctx, s.c, basePath+"Itemtrade/getItemtradeList", q)
}

// CountryTradeParams 는 국가별 수출입실적 요청 인자다. 빈 문자열 인자는 보내지 않는다.
type CountryTradeParams struct {
	Start  string // 시작년월 YYYYMM
	End    string // 종료년월 YYYYMM (Start 포함 최대 12개월)
	CntyCd string // cntyCd: 국가코드 (N). 대문자 2자리(US). 생략하면 전 국가. 코드는 codes/국가코드.csv
}

// CountryTradeRow 는 국가별 수출입실적 응답 행이다. 문서: docs/api/customs/국가별.md
type CountryTradeRow struct {
	Year           string       `xml:"year"`           // 기간. 2026.01 형식
	StatCd         string       `xml:"statCd"`         // 국가코드. US
	StatCdCntnKor1 string       `xml:"statCdCntnKor1"` // 국가명. 미국
	ExpCnt         opendata.Num `xml:"expCnt"`         // 수출건수 (건). 정수
	ExpDlr         opendata.Num `xml:"expDlr"`         // 수출금액 (달러). 정수
	ImpCnt         opendata.Num `xml:"impCnt"`         // 수입건수 (건). 정수
	ImpDlr         opendata.Num `xml:"impDlr"`         // 수입금액 (달러). 정수
	BalPayments    opendata.Num `xml:"balPayments"`    // 무역수지 (달러). 정수, 음수 가능
}

// IsTotal 은 기간 필드가 "총계" 인 합계 행인지 알려 준다.
func (r CountryTradeRow) IsTotal() bool { return r.Year == totalPeriod }

// CountryTrade 는 국가별 수출입실적을 조회한다 (nationtrade/getNationtradeList).
func (s *Service) CountryTrade(ctx context.Context, p CountryTradeParams) ([]CountryTradeRow, error) {
	q, err := periodQuery(p.Start, p.End)
	if err != nil {
		return nil, err
	}
	set(q, "cntyCd", p.CntyCd)
	return opendata.Fetch[CountryTradeRow](ctx, s.c, basePath+"nationtrade/getNationtradeList", q)
}

// ContinentTradeParams 는 대륙별 수출입실적 요청 인자다. 빈 문자열 인자는 보내지 않는다.
type ContinentTradeParams struct {
	Start             string // 시작년월 YYYYMM
	End               string // 종료년월 YYYYMM (Start 포함 최대 12개월)
	CntnEbkUnfcClsfCd string // cntnEbkUnfcClsfCd: 대륙경제권통합분류코드 (N). 대륙코드 2자리(10 아시아). 생략하면 전 대륙. 코드는 codes/대륙코드.csv
}

// ContinentTradeRow 는 대륙별 수출입실적 응답 행이다. 문서: docs/api/customs/대륙별.md
type ContinentTradeRow struct {
	Year           string       `xml:"year"`           // 기간. 2026.01 형식
	StatCd         string       `xml:"statCd"`         // 대륙코드. 10 같은 2자리 문자열. 명세는 "국가부호코드"(number)라고 적었지만 실제 값은 대륙코드다
	StatCdCntnKor1 string       `xml:"statCdCntnKor1"` // 대륙명. 아시아. 명세 명칭은 "국가부호명"
	ExpCnt         opendata.Num `xml:"expCnt"`         // 수출건수 (건). 정수
	ExpDlr         opendata.Num `xml:"expDlr"`         // 수출금액 (달러). 정수
	ImpCnt         opendata.Num `xml:"impCnt"`         // 수입건수 (건). 정수
	ImpDlr         opendata.Num `xml:"impDlr"`         // 수입금액 (달러). 정수
	BalPayments    opendata.Num `xml:"balPayments"`    // 무역수지 (달러). 정수, 음수 가능
}

// IsTotal 은 기간 필드가 "총계" 인 합계 행인지 알려 준다.
func (r ContinentTradeRow) IsTotal() bool { return r.Year == totalPeriod }

// ContinentTrade 는 대륙별 수출입실적을 조회한다 (continenttradet/getContinenttradeList).
func (s *Service) ContinentTrade(ctx context.Context, p ContinentTradeParams) ([]ContinentTradeRow, error) {
	q, err := periodQuery(p.Start, p.End)
	if err != nil {
		return nil, err
	}
	set(q, "cntnEbkUnfcClsfCd", p.CntnEbkUnfcClsfCd)
	return opendata.Fetch[ContinentTradeRow](ctx, s.c, basePath+"continenttradet/getContinenttradeList", q)
}

// EconomicBlocTradeParams 는 경제권별 수출입실적 요청 인자다. 빈 문자열 인자는 보내지 않는다.
type EconomicBlocTradeParams struct {
	Start             string // 시작년월 YYYYMM
	End               string // 종료년월 YYYYMM (Start 포함 최대 12개월)
	CntnEbkUnfcClsfCd string // cntnEbkUnfcClsfCd: 대륙경제권통합분류코드 (N). 경제권코드 2자리(10 EU). 생략하면 전 경제권. 코드는 codes/경제권코드.csv
}

// EconomicBlocTradeRow 는 경제권별 수출입실적 응답 행이다. 문서: docs/api/customs/경제권별.md
type EconomicBlocTradeRow struct {
	Year           string       `xml:"year"`           // 기간. 2026.01 형식
	StatCd         string       `xml:"statCd"`         // 경제권코드. 10 같은 2자리 문자열(명세 타입은 number)
	StatCdCntnKor1 string       `xml:"statCdCntnKor1"` // 경제권명. EU
	ExpCnt         opendata.Num `xml:"expCnt"`         // 수출건수 (건). 정수
	ExpDlr         opendata.Num `xml:"expDlr"`         // 수출금액 (달러). 정수
	ImpCnt         opendata.Num `xml:"impCnt"`         // 수입건수 (건). 정수
	ImpDlr         opendata.Num `xml:"impDlr"`         // 수입금액 (달러). 정수
	BalPayments    opendata.Num `xml:"balPayments"`    // 무역수지 (달러). 정수, 음수 가능
}

// IsTotal 은 기간 필드가 "총계" 인 합계 행인지 알려 준다.
func (r EconomicBlocTradeRow) IsTotal() bool { return r.Year == totalPeriod }

// EconomicBlocTrade 는 경제권별 수출입실적을 조회한다 (economytrade/getEconomytradeList).
func (s *Service) EconomicBlocTrade(ctx context.Context, p EconomicBlocTradeParams) ([]EconomicBlocTradeRow, error) {
	q, err := periodQuery(p.Start, p.End)
	if err != nil {
		return nil, err
	}
	set(q, "cntnEbkUnfcClsfCd", p.CntnEbkUnfcClsfCd)
	return opendata.Fetch[EconomicBlocTradeRow](ctx, s.c, basePath+"economytrade/getEconomytradeList", q)
}

// PropertyTradeParams 는 성질별 수출입실적 요청 인자다. 빈 문자열 인자는 보내지 않는다.
type PropertyTradeParams struct {
	Start          string // 시작년월 YYYYMM
	End            string // 종료년월 YYYYMM (Start 포함 최대 12개월)
	ImexTpcd       string // imexTpcd: 수출입구분코드 (Y). 1 수출, 2 수입(codes/수출수입코드.csv). 생략하면 99
	ImexTmprClsfCd string // imexTmprClsfCd: 수출입성질분류코드 (N). 5자리 문자열(11101, 4Z001). 수출은 코드표 왼쪽 표, 수입은 오른쪽 표(codes/성질분류코드.csv). 명세 타입은 number 지만 영문이 섞인다. 생략하면 전 성질
}

// PropertyTradeRow 는 성질별 수출입실적 응답 행이다. 문서: docs/api/customs/성질별.md
type PropertyTradeRow struct {
	Year    string       `xml:"year"`    // 기간. 2026.01 형식
	Impexp  string       `xml:"impexp"`  // 수출입구분. 수출 또는 수입
	GodsCd  string       `xml:"godsCd"`  // 성질코드. 5자리 문자열. 4Z001·26A01 처럼 영문이 섞인다
	GodsKor string       `xml:"godsKor"` // 성질명. - 돼지고기, (참  치) 처럼 코드표 표기 그대로(앞의 - ·괄호·내부 공백 포함)
	Wgt     opendata.Num `xml:"wgt"`     // 중량 (kg). 정수
	Dlr     opendata.Num `xml:"dlr"`     // 금액 (달러). 정수
}

// IsTotal 은 기간 필드가 "총계" 인 합계 행인지 알려 준다.
func (r PropertyTradeRow) IsTotal() bool { return r.Year == totalPeriod }

// PropertyTrade 는 성질별 수출입실적을 조회한다 (Idfytempertrade/getIdfytempertradeList).
func (s *Service) PropertyTrade(ctx context.Context, p PropertyTradeParams) ([]PropertyTradeRow, error) {
	q, err := periodQuery(p.Start, p.End)
	if err != nil {
		return nil, err
	}
	set(q, "imexTpcd", p.ImexTpcd)
	set(q, "imexTmprClsfCd", p.ImexTmprClsfCd)
	return opendata.Fetch[PropertyTradeRow](ctx, s.c, basePath+"Idfytempertrade/getIdfytempertradeList", q)
}

// PropertyCountryTradeParams 는 성질별 국가별 수출입실적 요청 인자다. 빈 문자열 인자는 보내지 않는다.
type PropertyCountryTradeParams struct {
	Start          string // 시작년월 YYYYMM
	End            string // 종료년월 YYYYMM (Start 포함 최대 12개월)
	ImexTpcd       string // imexTpcd: 수출입구분코드 (Y). 1 수출, 2 수입(codes/수출수입코드.csv)
	ImexTmprClsfCd string // imexTmprClsfCd: 수출입성질분류코드 (N). 5자리 문자열(11201). 수출·수입 코드표가 다르다(codes/성질분류코드.csv). 생략하면 전 성질
	CntyCd         string // cntyCd: 국가코드 (Y). 대문자 2자리(US). 코드는 codes/국가코드.csv
}

// PropertyCountryTradeRow 는 성질별 국가별 수출입실적 응답 행이다. 문서: docs/api/customs/성질별국가별.md
type PropertyCountryTradeRow struct {
	Year           string       `xml:"year"`           // 기간. 2026.01 형식
	Impexp         string       `xml:"impexp"`         // 수출입구분. 수출 또는 수입
	StatCd         string       `xml:"statCd"`         // 국가코드. US
	StatCdCntnKor1 string       `xml:"statCdCntnKor1"` // 국가명. 미국
	GodsCd         string       `xml:"godsCd"`         // 성질코드. 5자리 문자열. 4Z001 처럼 영문이 섞인다
	GodsKor        string       `xml:"godsKor"`        // 성질명. - 기타 육류 및 조제품, (참  치) 처럼 코드표 표기 그대로
	Wgt            opendata.Num `xml:"wgt"`            // 중량 (kg). 정수
	Dlr            opendata.Num `xml:"dlr"`            // 금액 (달러). 정수
}

// IsTotal 은 기간 필드가 "총계" 인 합계 행인지 알려 준다.
func (r PropertyCountryTradeRow) IsTotal() bool { return r.Year == totalPeriod }

// PropertyCountryTrade 는 성질별 국가별 수출입실적을 조회한다 (ntempertrade/getNtempertradeList).
func (s *Service) PropertyCountryTrade(ctx context.Context, p PropertyCountryTradeParams) ([]PropertyCountryTradeRow, error) {
	q, err := periodQuery(p.Start, p.End)
	if err != nil {
		return nil, err
	}
	set(q, "imexTpcd", p.ImexTpcd)
	set(q, "imexTmprClsfCd", p.ImexTmprClsfCd)
	set(q, "cntyCd", p.CntyCd)
	return opendata.Fetch[PropertyCountryTradeRow](ctx, s.c, basePath+"ntempertrade/getNtempertradeList", q)
}

// NewPropertyTradeParams 는 신성질별 수출입실적 요청 인자다. 빈 문자열 인자는 보내지 않는다.
type NewPropertyTradeParams struct {
	Start              string // 시작년월 YYYYMM
	End                string // 종료년월 YYYYMM (Start 포함 최대 12개월)
	ImexTpcd           string // imexTpcd: 수출입구분코드 (Y). 1 수출, 2 수입(codes/수출수입코드.csv)
	ImexTmprUnfcClsfCd string // imexTmprUnfcClsfCd: 수출입성질통합분류코드 (N). 8자리 세세분류코드(13050102) — codes/성질통합분류코드.csv 15번째 열 관세청 신성질별 분류세세분류코드. 상위 분류코드는 0건(함정). 생략하면 전 성질
}

// NewPropertyTradeRow 는 신성질별 수출입실적 응답 행이다. 문서: docs/api/customs/신성질별.md
type NewPropertyTradeRow struct {
	Year           string       `xml:"year"`           // 기간. 2026.01 형식
	Impexp         string       `xml:"impexp"`         // 수출입구분. 수출 또는 수입
	StatCd         string       `xml:"statCd"`         // 국가코드. US. 국가 필터 인자는 없지만 행은 국가별로 나뉜다. ZZ 는 기타국
	StatCdCntnKor1 string       `xml:"statCdCntnKor1"` // 국가명. 미국. & 는 &amp; 로 온다
	GodsCd         string       `xml:"godsCd"`         // 성질코드. 8자리 세세분류코드 문자열(명세 타입은 number)
	GodsKor        string       `xml:"godsKor"`        // 성질명. (기초화장품) 처럼 괄호로 감싼 세세분류명
	Wgt            opendata.Num `xml:"wgt"`            // 중량 (kg). 정수
	Dlr            opendata.Num `xml:"dlr"`            // 금액 (달러). 정수
}

// IsTotal 은 기간 필드가 "총계" 인 합계 행인지 알려 준다.
func (r NewPropertyTradeRow) IsTotal() bool { return r.Year == totalPeriod }

// NewPropertyTrade 는 신성질별 수출입실적을 조회한다 (newtempertrade/getNewtempertradeList).
func (s *Service) NewPropertyTrade(ctx context.Context, p NewPropertyTradeParams) ([]NewPropertyTradeRow, error) {
	q, err := periodQuery(p.Start, p.End)
	if err != nil {
		return nil, err
	}
	set(q, "imexTpcd", p.ImexTpcd)
	set(q, "imexTmprUnfcClsfCd", p.ImexTmprUnfcClsfCd)
	return opendata.Fetch[NewPropertyTradeRow](ctx, s.c, basePath+"newtempertrade/getNewtempertradeList", q)
}

// NewPropertyCountryTradeParams 는 신성질별 국가별 수출입실적 요청 인자다. 빈 문자열 인자는 보내지 않는다.
type NewPropertyCountryTradeParams struct {
	Start              string // 시작년월 YYYYMM
	End                string // 종료년월 YYYYMM (Start 포함 최대 12개월)
	ImexTpcd           string // imexTpcd: 수출입구분코드 (Y). 1 수출, 2 수입(codes/수출수입코드.csv)
	ImexTmprUnfcClsfCd string // imexTmprUnfcClsfCd: 수출입성질통합분류코드 (명세 Y / 실측 생략 가능). 8자리 세세분류코드(13050102) — codes/성질통합분류코드.csv 15번째 열. 생략하면 그 국가의 전 성질
	CntyCd             string // cntyCd: 국가코드 (Y). 대문자 2자리(US). 생략하면 99. 코드는 codes/국가코드.csv
}

// NewPropertyCountryTradeRow 는 신성질별 국가별 수출입실적 응답 행이다. 문서: docs/api/customs/신성질별국가별.md
type NewPropertyCountryTradeRow struct {
	Year           string       `xml:"year"`           // 기간. 2026.01 형식
	Impexp         string       `xml:"impexp"`         // 수출입구분. 수출 또는 수입
	StatCd         string       `xml:"statCd"`         // 국가코드. US
	StatCdCntnKor1 string       `xml:"statCdCntnKor1"` // 국가명. 미국
	GodsCd         string       `xml:"godsCd"`         // 성질코드. 8자리 세세분류코드 문자열(명세 타입은 number)
	GodsKor        string       `xml:"godsKor"`        // 성질명. (기초화장품) 처럼 괄호로 감싼 세세분류명
	Wgt            opendata.Num `xml:"wgt"`            // 중량 (kg). 정수
	Dlr            opendata.Num `xml:"dlr"`            // 금액 (달러). 정수
}

// IsTotal 은 기간 필드가 "총계" 인 합계 행인지 알려 준다.
func (r NewPropertyCountryTradeRow) IsTotal() bool { return r.Year == totalPeriod }

// NewPropertyCountryTrade 는 신성질별 국가별 수출입실적을 조회한다 (nnewtempertrade/getNnewtempertradeList).
func (s *Service) NewPropertyCountryTrade(ctx context.Context, p NewPropertyCountryTradeParams) ([]NewPropertyCountryTradeRow, error) {
	q, err := periodQuery(p.Start, p.End)
	if err != nil {
		return nil, err
	}
	set(q, "imexTpcd", p.ImexTpcd)
	set(q, "imexTmprUnfcClsfCd", p.ImexTmprUnfcClsfCd)
	set(q, "cntyCd", p.CntyCd)
	return opendata.Fetch[NewPropertyCountryTradeRow](ctx, s.c, basePath+"nnewtempertrade/getNnewtempertradeList", q)
}

// KindTradeParams 는 종류별 수출입실적 요청 인자다. 빈 문자열 인자는 보내지 않는다.
type KindTradeParams struct {
	Start    string // 시작년월 YYYYMM
	End      string // 종료년월 YYYYMM (Start 포함 최대 12개월)
	ImexKcd  string // imexKcd: 수출입종류코드 (N). 수출은 영문(A 일반수출), 수입은 숫자(11 일반수입(외화획득용)). 코드는 codes/수출입종류코드.csv — 수출·수입 표가 좌우로 나뉜다. 명세 타입은 number 지만 영문이다. 생략하면 전 종류
	ImexTpcd string // imexTpcd: 수출입구분코드 (Y). 1 수출, 2 수입(codes/수출수입코드.csv)
}

// KindTradeRow 는 종류별 수출입실적 응답 행이다. 문서: docs/api/customs/종류별.md
type KindTradeRow struct {
	Year           string       `xml:"year"`           // 기간. 202601 형식(점 없음). 총계 행은 총계
	Impexp         string       `xml:"impexp"`         // 수출입구분. 수출 또는 수입
	StatCd         string       `xml:"statCd"`         // 종류코드. 수출입종류코드(A). 국가코드가 아니다. 총계 행은 -
	StatCdCntnKor1 string       `xml:"statCdCntnKor1"` // 종류명. 일반수출. 총계 행은 -
	Cnt            opendata.Num `xml:"cnt"`            // 건수 (건). 정수
	Wgt            opendata.Num `xml:"wgt"`            // 중량 (kg). 정수
	Dlr            opendata.Num `xml:"dlr"`            // 금액 (달러). 정수
	Won            opendata.Num `xml:"won"`            // 금액 (원). 정수. 원화 금액이 있는 유일한 API
}

// IsTotal 은 기간 필드가 "총계" 인 합계 행인지 알려 준다.
func (r KindTradeRow) IsTotal() bool { return r.Year == totalPeriod }

// KindTrade 는 종류별 수출입실적을 조회한다 (kindtrade/getKindtradeList).
func (s *Service) KindTrade(ctx context.Context, p KindTradeParams) ([]KindTradeRow, error) {
	q, err := periodQuery(p.Start, p.End)
	if err != nil {
		return nil, err
	}
	set(q, "imexKcd", p.ImexKcd)
	set(q, "imexTpcd", p.ImexTpcd)
	return opendata.Fetch[KindTradeRow](ctx, s.c, basePath+"kindtrade/getKindtradeList", q)
}

// CustomsOfficeTradeParams 는 세관별 수출입실적 요청 인자다. 빈 문자열 인자는 보내지 않는다.
type CustomsOfficeTradeParams struct {
	Start     string // 시작년월 YYYYMM
	End       string // 종료년월 YYYYMM (Start 포함 최대 12개월)
	CstmSgnYn string // cstmSgnYn: 세관구분여부 (N). Y 본부세관 단위(월 7행), N 개별 세관 단위(월 45행). 생략하면 Y 와 같다. codes/세관구분코드.csv
}

// CustomsOfficeTradeRow 는 세관별 수출입실적 응답 행이다. 문서: docs/api/customs/세관별.md
type CustomsOfficeTradeRow struct {
	Year          string       `xml:"year"`          // 기간. 202601 형식(점 없음)
	Center        string       `xml:"center"`        // 본부세관. 서울본부, 인천공항본부, 직할세관 등 7개
	Cstm          string       `xml:"cstm"`          // 세관코드. 3자리 문자열(070). Y/생략이면 -
	StatCdCntnKor string       `xml:"statCdCntnKor"` // 세관명. 목포세관. Y/생략이면 -
	ExpCnt        opendata.Num `xml:"expCnt"`        // 수출건수 (건). 정수
	ExpDlr        opendata.Num `xml:"expDlr"`        // 수출금액 (달러). 정수
	ImpCnt        opendata.Num `xml:"impCnt"`        // 수입건수 (건). 정수
	ImpDlr        opendata.Num `xml:"impDlr"`        // 수입금액 (달러). 정수
	BalPayments   opendata.Num `xml:"balPayments"`   // 무역수지 (달러). 정수, 음수 가능
}

// IsTotal 은 기간 필드가 "총계" 인 합계 행인지 알려 준다.
func (r CustomsOfficeTradeRow) IsTotal() bool { return r.Year == totalPeriod }

// CustomsOfficeTrade 는 세관별 수출입실적을 조회한다 (customstrade/getCustomstradeList).
func (s *Service) CustomsOfficeTrade(ctx context.Context, p CustomsOfficeTradeParams) ([]CustomsOfficeTradeRow, error) {
	q, err := periodQuery(p.Start, p.End)
	if err != nil {
		return nil, err
	}
	set(q, "cstmSgnYn", p.CstmSgnYn)
	return opendata.Fetch[CustomsOfficeTradeRow](ctx, s.c, basePath+"customstrade/getCustomstradeList", q)
}

// PortTradeParams 는 항구·공항별 수출입실적 요청 인자다. 빈 문자열 인자는 보내지 않는다.
type PortTradeParams struct {
	Start           string // 시작년월 YYYYMM
	End             string // 종료년월 YYYYMM (Start 포함 최대 12개월)
	PortAirptRegnCd string // portAirptRegnCd: 항구공항지역코드 (N). KRPUS(부산항), ICN(인천국제공항), 010(서울세관) 등. 코드는 codes/항구공항코드.csv. 생략하면 전 항구
}

// PortTradeRow 는 항구·공항별 수출입실적 응답 행이다. 문서: docs/api/customs/항구공항별.md
type PortTradeRow struct {
	Year        string       `xml:"year"`        // 기간. 202601 형식(점 없음). 총계 행은 총계
	PortCd      string       `xml:"portCd"`      // 항구공항지역코드. 문자열. 5자리 항만코드(KRPUS), 3자리 공항코드(ICN), 3자리 세관코드(010), ZZZZZ(기타항). 총계 행은 -
	StatKor     string       `xml:"statKor"`     // 항구공항지역명. 부산항, 인천국제공항, 서울세관. 총계 행은 -
	CstmSgn     string       `xml:"cstmSgn"`     // 세관코드. 신고 세관 3자리 문자열(010). 총계 행에는 태그가 없다
	ExpCnt      opendata.Num `xml:"expCnt"`      // 수출건수 (건). 정수
	ExpDlr      opendata.Num `xml:"expDlr"`      // 수출금액 (달러). 정수
	ImpCnt      opendata.Num `xml:"impCnt"`      // 수입건수 (건). 정수
	ImpDlr      opendata.Num `xml:"impDlr"`      // 수입금액 (달러). 정수
	BalPayments opendata.Num `xml:"balPayments"` // 무역수지 (달러). 정수, 음수 가능
}

// IsTotal 은 기간 필드가 "총계" 인 합계 행인지 알려 준다.
func (r PortTradeRow) IsTotal() bool { return r.Year == totalPeriod }

// PortTrade 는 항구·공항별 수출입실적을 조회한다 (porttrade/getPorttradeList).
func (s *Service) PortTrade(ctx context.Context, p PortTradeParams) ([]PortTradeRow, error) {
	q, err := periodQuery(p.Start, p.End)
	if err != nil {
		return nil, err
	}
	set(q, "portAirptRegnCd", p.PortAirptRegnCd)
	return opendata.Fetch[PortTradeRow](ctx, s.c, basePath+"porttrade/getPorttradeList", q)
}

// SidoTradeParams 는 시도별 수출입실적 요청 인자다. 빈 문자열 인자는 보내지 않는다.
type SidoTradeParams struct {
	Start  string // 시작년월 YYYYMM
	End    string // 종료년월 YYYYMM (Start 포함 최대 12개월)
	SidoCd string // sidoCd: 시도코드 (N). 2자리(41 경기도). 코드는 codes/시도코드.csv — 광주·전남은 2026.7.1. 전후로 코드가 다르다(코드표 함정). 생략하면 전 시도
}

// SidoTradeRow 는 시도별 수출입실적 응답 행이다. 문서: docs/api/customs/시도별.md
type SidoTradeRow struct {
	PriodTitle  string       `xml:"priodTitle"`  // 기간. 연도 2026(월이 아님). 총계 행은 총계
	SidoNm      string       `xml:"sidoNm"`      // 시도명. 서울특별시. 총계 행에는 태그가 없다
	ExpCnt      opendata.Num `xml:"expCnt"`      // 수출건수 (건)
	ExpUsdAmt   opendata.Num `xml:"expUsdAmt"`   // 수출금액 (천 달러). 명세는 "달러"
	ImpCnt      opendata.Num `xml:"impCnt"`      // 수입건수 (건)
	ImpUsdAmt   opendata.Num `xml:"impUsdAmt"`   // 수입금액 (천 달러). 명세는 "달러"
	CmtrBlncAmt opendata.Num `xml:"cmtrBlncAmt"` // 무역수지 (천 달러). 음수 가능(     -68,733,092)
}

// IsTotal 은 기간 필드가 "총계" 인 합계 행인지 알려 준다.
func (r SidoTradeRow) IsTotal() bool { return r.PriodTitle == totalPeriod }

// SidoTrade 는 시도별 수출입실적을 조회한다 (sidotrade/getSidotradeList).
func (s *Service) SidoTrade(ctx context.Context, p SidoTradeParams) ([]SidoTradeRow, error) {
	q, err := periodQuery(p.Start, p.End)
	if err != nil {
		return nil, err
	}
	set(q, "sidoCd", p.SidoCd)
	return opendata.Fetch[SidoTradeRow](ctx, s.c, basePath+"sidotrade/getSidotradeList", q)
}

// SidoItemTradeParams 는 시도별 품목별 수출입실적 요청 인자다. 빈 문자열 인자는 보내지 않는다.
type SidoItemTradeParams struct {
	Start  string // 시작년월 YYYYMM
	End    string // 종료년월 YYYYMM (Start 포함 최대 12개월)
	SidoCd string // sidoCd: 시도코드 (Y). 2자리(41 경기도). 코드는 codes/시도코드.csv — 광주·전남은 2026.7.1. 전후로 코드가 다르다(코드표 함정)
}

// SidoItemTradeRow 는 시도별 품목별 수출입실적 응답 행이다. 문서: docs/api/customs/시도별품목별.md
type SidoItemTradeRow struct {
	PriodTitle  string       `xml:"priodTitle"`  // 기간. 연도 2026(월이 아님). 총계 행은 총계. 명세 타입은 number
	HsSgn       string       `xml:"hsSgn"`       // 품목코드. HS 2자리 문자열(01). 총계 행에는 태그가 없다
	KorePrlstNm string       `xml:"korePrlstNm"` // 품목명. 류 명칭(살아 있는 동물). 총계 행에는 태그가 없다
	ExpLnCnt    opendata.Num `xml:"expLnCnt"`    // 수출 란수 (건). 명세는 "수출품목건수". API 13 의 expCnt(신고 건수)와 값이 다르다
	ExpUsdAmt   opendata.Num `xml:"expUsdAmt"`   // 수출금액 (천 달러). 명세는 "달러"
	ImpLnCnt    opendata.Num `xml:"impLnCnt"`    // 수입 란수 (건). 명세는 "수입품목건수"
	ImpUsdAmt   opendata.Num `xml:"impUsdAmt"`   // 수입금액 (천 달러). 명세는 "달러"
	CmtrBlncAmt opendata.Num `xml:"cmtrBlncAmt"` // 무역수지 (천 달러). 음수 가능
}

// IsTotal 은 기간 필드가 "총계" 인 합계 행인지 알려 준다.
func (r SidoItemTradeRow) IsTotal() bool { return r.PriodTitle == totalPeriod }

// SidoItemTrade 는 시도별 품목별 수출입실적을 조회한다 (sidoitemtrade/getSidoitemtradeList).
func (s *Service) SidoItemTrade(ctx context.Context, p SidoItemTradeParams) ([]SidoItemTradeRow, error) {
	q, err := periodQuery(p.Start, p.End)
	if err != nil {
		return nil, err
	}
	set(q, "sidoCd", p.SidoCd)
	return opendata.Fetch[SidoItemTradeRow](ctx, s.c, basePath+"sidoitemtrade/getSidoitemtradeList", q)
}

// SidoPropertyTradeParams 는 시도별 성질별 수출입실적 요청 인자다. 빈 문자열 인자는 보내지 않는다.
type SidoPropertyTradeParams struct {
	Start          string // 시작년월 YYYYMM
	End            string // 종료년월 YYYYMM (Start 포함 최대 12개월)
	DtlTmprYn      string // dtlTmprYn: 세부성질여부 (N). Y 면 세분류(5자리 말단)까지 포함(경기도 수출 6개월 175행). 생략하면 대·중분류만(22행). 행 수는 총계 포함
	SidoCd         string // sidoCd: 시도코드 (Y). 2자리(41 경기도). 코드는 codes/시도코드.csv — 광주·전남은 2026.7.1. 전후로 코드가 다르다(코드표 함정)
	ImexTpcd       string // imexTpcd: 수출입구분코드 (Y). 1 수출, 2 수입(codes/수출수입코드.csv)
	ImexTmprClsfCd string // imexTmprClsfCd: 수출입성질분류코드 (N). 5자리 문자열(11101, 4Z000). codes/성질분류코드.csv — 수출·수입 표가 다르다. 미실측. 생략하면 전 성질
}

// SidoPropertyTradeRow 는 시도별 성질별 수출입실적 응답 행이다. 문서: docs/api/customs/시도별성질별.md
type SidoPropertyTradeRow struct {
	PriodTitle  string       `xml:"priodTitle"`  // 기간. 연도 2026(월이 아님). 총계 행은 총계. 명세 타입은 number
	TmprTpcd    string       `xml:"tmprTpcd"`    // 성질코드. 5자리 문자열(10000, 4Z000). 총계 행은 00000
	CdValtValNm string       `xml:"cdValtValNm"` // 성질명. 1. 식료 및 직접소비재, 가. 섬유원료, - 돼지고기. 총계 행에는 태그가 없다
	ImexLnCnt   opendata.Num `xml:"imexLnCnt"`   // 수출입 란수 (건). 명세는 "수출품목건수"
	ImexUsdAmt  opendata.Num `xml:"imexUsdAmt"`  // 수출입금액 (천 달러). 명세는 "수출금액(달러)"
}

// IsTotal 은 기간 필드가 "총계" 인 합계 행인지 알려 준다.
func (r SidoPropertyTradeRow) IsTotal() bool { return r.PriodTitle == totalPeriod }

// SidoPropertyTrade 는 시도별 성질별 수출입실적을 조회한다 (sidotempertrade/getSidotempertradeList).
func (s *Service) SidoPropertyTrade(ctx context.Context, p SidoPropertyTradeParams) ([]SidoPropertyTradeRow, error) {
	q, err := periodQuery(p.Start, p.End)
	if err != nil {
		return nil, err
	}
	set(q, "dtlTmprYn", p.DtlTmprYn)
	set(q, "sidoCd", p.SidoCd)
	set(q, "imexTpcd", p.ImexTpcd)
	set(q, "imexTmprClsfCd", p.ImexTmprClsfCd)
	return opendata.Fetch[SidoPropertyTradeRow](ctx, s.c, basePath+"sidotempertrade/getSidotempertradeList", q)
}

// SigunguTradeParams 는 시군구별 수출입실적 요청 인자다. 빈 문자열 인자는 보내지 않는다.
type SigunguTradeParams struct {
	Start  string // 시작년월 YYYYMM
	End    string // 종료년월 YYYYMM (Start 포함 최대 12개월)
	SidoCd string // sidoCd: 시도코드 (Y). 2자리(41 경기도). 코드는 codes/시도코드.csv — 광주·전남은 2026.7.1. 전후로 코드가 다르다(코드표 함정). 생략하면 99
}

// SigunguTradeRow 는 시군구별 수출입실적 응답 행이다. 문서: docs/api/customs/시군구별.md
type SigunguTradeRow struct {
	PriodTitle  string       `xml:"priodTitle"`  // 기간. 2026.01 형식(월별). 시도 계열(13~15)과 다르다
	SidoSggNm   string       `xml:"sidoSggNm"`   // 시군구명. 시도명 포함(경기도 가평군)
	ExpCnt      opendata.Num `xml:"expCnt"`      // 수출건수 (건)
	ExpUsdAmt   opendata.Num `xml:"expUsdAmt"`   // 수출금액 (천 달러). 명세는 "수출미화금액"(단위 없음)
	ImpCnt      opendata.Num `xml:"impCnt"`      // 수입건수 (건)
	ImpUsdAmt   opendata.Num `xml:"impUsdAmt"`   // 수입금액 (천 달러)
	CmtrBlncAmt opendata.Num `xml:"cmtrBlncAmt"` // 무역수지 (천 달러). 음수 가능
}

// IsTotal 은 기간 필드가 "총계" 인 합계 행인지 알려 준다.
func (r SigunguTradeRow) IsTotal() bool { return r.PriodTitle == totalPeriod }

// SigunguTrade 는 시군구별 수출입실적을 조회한다 (sigunguperimexacrs/getSigunguPerImexAcrs).
func (s *Service) SigunguTrade(ctx context.Context, p SigunguTradeParams) ([]SigunguTradeRow, error) {
	q, err := periodQuery(p.Start, p.End)
	if err != nil {
		return nil, err
	}
	set(q, "sidoCd", p.SidoCd)
	return opendata.Fetch[SigunguTradeRow](ctx, s.c, basePath+"sigunguperimexacrs/getSigunguPerImexAcrs", q)
}

// SigunguItemTradeParams 는 시군구별 품목별 수출입실적 요청 인자다. 빈 문자열 인자는 보내지 않는다.
type SigunguItemTradeParams struct {
	Start  string // 시작년월 YYYYMM
	End    string // 종료년월 YYYYMM (Start 포함 최대 12개월)
	HsSgn  string // HsSgn: HS부호 6단위 (Y). 대문자 H(명세). HS 6자리(854232 메모리). 소문자 hsSgn 도 같은 결과(실측). 코드는 codes/품목코드.csv
	SidoCd string // sidoCd: 시도코드 (Y). 2자리(41 경기도). 코드는 codes/시도코드.csv — 광주·전남은 2026.7.1. 전후로 코드가 다르다(코드표 함정)
}

// SigunguItemTradeRow 는 시군구별 품목별 수출입실적 응답 행이다. 문서: docs/api/customs/시군구별품목별.md
type SigunguItemTradeRow struct {
	PriodTitle  string       `xml:"priodTitle"`  // 기간. 2026.01 형식(월별)
	SggNm       string       `xml:"sggNm"`       // 시군구명. 시도명 포함(경기도 고양시). 시군구별은 같은 값이 sidoSggNm 이다
	HsSgn       string       `xml:"hsSgn"`       // HS부호. 요청한 6자리 그대로(854232). 소문자 h
	KorePrlstNm string       `xml:"korePrlstNm"` // 품목명. 메모리
	ExpCnt      opendata.Num `xml:"expCnt"`      // 수출건수 (건)
	ExpUsdAmt   opendata.Num `xml:"expUsdAmt"`   // 수출금액 (천 달러). 명세는 "수출미화금액"(단위 없음)
	ImpCnt      opendata.Num `xml:"impCnt"`      // 수입건수 (건)
	ImpUsdAmt   opendata.Num `xml:"impUsdAmt"`   // 수입금액 (천 달러)
	CmtrBlncAmt opendata.Num `xml:"cmtrBlncAmt"` // 무역수지 (천 달러). 음수 가능
}

// IsTotal 은 기간 필드가 "총계" 인 합계 행인지 알려 준다.
func (r SigunguItemTradeRow) IsTotal() bool { return r.PriodTitle == totalPeriod }

// SigunguItemTrade 는 시군구별 품목별 수출입실적을 조회한다 (sigunguperprlstperacrs/getSigunguPerPrlstPerAcrs).
func (s *Service) SigunguItemTrade(ctx context.Context, p SigunguItemTradeParams) ([]SigunguItemTradeRow, error) {
	q, err := periodQuery(p.Start, p.End)
	if err != nil {
		return nil, err
	}
	set(q, "HsSgn", p.HsSgn)
	set(q, "sidoCd", p.SidoCd)
	return opendata.Fetch[SigunguItemTradeRow](ctx, s.c, basePath+"sigunguperprlstperacrs/getSigunguPerPrlstPerAcrs", q)
}
