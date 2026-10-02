# opendata-go

공공데이터포털([data.go.kr](https://www.data.go.kr)) OpenAPI 의 Go 클라이언트 라이브러리.
포털 공통(서비스키·게이트웨이·응답 봉투·에러)을 루트 패키지 `opendata` 가 맡고, 기관별 API 는 서브패키지로 붙인다.

## 설치

```bash
go get github.com/kenshin579/opendata-go@latest
```

Go 1.25+, 외부 의존성 없음(표준 라이브러리만).

## 사용

```go
c, _ := opendata.NewClientFromEnv() // OPENDATA_API_KEY
s := customs.New(c)
ctx := context.Background()

// 국가별 수출입실적 (미국, 2026년 1~6월)
rows, err := s.CountryTrade(ctx, customs.CountryTradeParams{Start: "202601", End: "202606", CntyCd: "US"})
for _, r := range customs.WithoutTotal(rows) {
    exp, _ := r.ExpDlr.Int64() // 달러
    fmt.Println(r.Year, r.StatCdCntnKor1, exp)
}

// 한 번에 12개월까지만 조회된다 — 긴 기간은 연도별 창으로 나눠 부른다
windows, _ := customs.Years("202301", "202606")
for _, w := range windows {
    part, err := s.SidoTrade(ctx, customs.SidoTradeParams{Start: w.Start, End: w.End})
    ...
}
```

예제: [`examples/basic`](examples/basic/main.go) — 경기도 시군구별 메모리 반도체 수출.

## 주의

- **금액 단위가 API 마다 다르다.** `Sido*`·`Sigungu*` 는 천 달러, 나머지는 달러(종류별 `Won` 은 원).
  필드 주석과 [관세청 문서](docs/api/customs/README.md)를 확인할 것.
- **총계 행이 섞여 온다.** 위치가 API 마다 다르므로 `IsTotal()`/`customs.WithoutTotal` 로 거른다.
- 숫자는 `opendata.Num`(원문 보존). `Int64()` 가 공백 패딩·천 단위 콤마를 처리한다. 총계 행의 `-` 는 `ErrNotNumber`.
- 페이지네이션이 없다 — 응답 크기는 기간·필터 인자로만 줄일 수 있다(예: 신성질별 1개월 ≈ 2만 행).
- 결과 0건은 에러가 아니라 빈 슬라이스다.

## 에러

```go
var ge *opendata.GatewayError // 서비스키 누락·미등록·활용신청 안 한 API·트래픽 초과 (HTTP 401/403/400)
var ae *opendata.APIError     // 기관 서비스 실패 resultCode ≠ 00 (관세청은 99, 원인은 Message)
```

에러 문자열에 서비스키는 담기지 않는다(`{SERVICE_KEY}` 로 마스킹).

## 옵션

```go
c, _ := opendata.NewClient(key,
    opendata.WithTimeout(10*time.Second), // HTTP 타임아웃 (기본 30s)
    opendata.WithBaseURL("https://..."),  // 게이트웨이 주소 교체 (테스트/프록시)
    opendata.WithHTTPClient(custom),      // *http.Client 주입
)
```

## 지원 기관

| 기관 | 서브패키지 | API | 문서 |
| --- | --- | --- | --- |
| 관세청 (1220000) | `customs` | 수출입실적 17개 | [docs/api/customs](docs/api/customs/README.md) |

관세청 메서드: `ItemCountryTrade` `ItemTrade` `CountryTrade` `ContinentTrade` `EconomicBlocTrade`
`PropertyTrade` `PropertyCountryTrade` `NewPropertyTrade` `NewPropertyCountryTrade` `KindTrade`
`CustomsOfficeTrade` `PortTrade` `SidoTrade` `SidoItemTrade` `SidoPropertyTrade` `SigunguTrade` `SigunguItemTrade`

## 문서

- [포털 공통 규약](docs/api/README.md) — 서비스키, 게이트웨이, 응답 봉투, 에러
- 포털 명세 원본·코드표 다시 받기: `python3 scripts/portal-spec/harvest.py`

## 인증

data.go.kr 에서 발급받은 서비스키를 `OPENDATA_API_KEY` 환경변수로 두거나 `opendata.NewClient(key)` 로 넘긴다.
API 마다 포털에서 활용신청이 필요하다(안 하면 `GatewayError` 코드 30).

## 테스트

```bash
go test ./...                       # 단위 테스트 (fixture)
go test -tags integration ./...     # 실 API 호출 (OPENDATA_API_KEY 필요)
```

## 릴리스

```bash
./scripts/release.sh vX.Y.Z
```
