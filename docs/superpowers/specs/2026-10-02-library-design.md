# opendata-go 설계 — 2단계: 라이브러리 구현 (루트 `opendata` + `customs`)

- 날짜: 2026-10-02
- 상태: 승인됨 (사용자 "ok" — 1단계 스펙의 2단계 방향을 바탕으로 진행)
- 근거 문서: `docs/api/README.md`(포털 공통), `docs/api/customs/`(17개 API, 실측 반영). 1단계 스펙
  `2026-10-02-customs-api-docs-design.md` 의 "(참고) 2단계" 절을 여기서 확정한다.

## 목표

관세청 수출입실적 API 17개를 Go 로 호출하는 클라이언트를 만들고 v0.1.0 으로 릴리스한다.
moneyflow backend 가 `go.mod` 태그 버전으로 가져다 쓴다. 외부 의존성 0, Go 1.25.

## 패키지 구조

```
opendata-go/                 module github.com/kenshin579/opendata-go
├── client.go  config.go     package opendata — Client, Option, NewClient/NewClientFromEnv
├── call.go                  Fetch[T] — GET + XML 봉투 파싱 + 에러 매핑 (기관 패키지가 쓰는 유일한 통로)
├── errors.go                GatewayError, APIError
├── num.go                   Num — 원본 문자열 보존 숫자 타입
├── customs/
│   ├── customs.go           Service, New(c), Period·검증, Years(창 분할), IsTotal
│   ├── apis.go              17개 API 의 Params·Row 타입과 메서드 (문서 표에서 생성)
│   └── *_test.go
├── testdata/customs/        실 응답 fixture (축약본, 17 + 에러 케이스)
├── examples/basic/main.go
├── integration_test.go      //go:build integration
└── scripts/release.sh       ecos-go 것 재사용
```

`internal/httpclient` 를 따로 두지 않는다 — 기관 서브패키지가 루트 패키지의 공개 함수 `Fetch` 를
불러야 하므로 HTTP 계층이 루트에 있어야 한다. 공개 범위는 `Fetch` 하나로 좁힌다.

## 루트 패키지 `opendata`

- `NewClient(serviceKey string, opts ...Option) (*Client, error)` — 빈 키는 에러.
  `NewClientFromEnv()` 는 `OPENDATA_API_KEY`.
- 옵션: `WithBaseURL`(기본 `https://apis.data.go.kr`), `WithTimeout`(기본 30s), `WithHTTPClient`.
- `Fetch[T any](ctx, c *Client, path string, q url.Values) ([]T, error)`
  - `GET {baseURL}{path}?serviceKey=...&{q}`. `_type` 은 절대 붙이지 않는다(붙이면 게이트웨이 에러 04).
  - 응답 본문의 **루트 요소로 가른다**(HTTP 상태가 아니라):
    - `OpenAPI_ServiceResponse` → `*GatewayError{HTTPStatus, Code(returnReasonCode), ErrMsg, AuthMsg}`
    - `response` 이고 `header/resultCode` ≠ `00` → `*APIError{Code, Message}` (관세청은 99 하나뿐, 원인은 Message)
    - `response` + `00` → `body/items/item` 을 `[]T` 로. `<items/>`·body 없음·`<body/>`·`<body><totalCount>0</totalCount></body>` 모두 **빈 슬라이스, nil 에러**
    - 그 밖(빈 본문, HTML, 디코드 실패) → 상태 코드와 본문 앞 200바이트를 담은 일반 에러
  - 반환하는 모든 에러 문자열에서 serviceKey 를 `{SERVICE_KEY}` 로 마스킹(`*url.Error` 는 URL 전체를 담는다).
  - 페이지네이션 없음(실측) → 페이지 인자·반복 헬퍼 없음.
- `Num` (string 기반 타입): XML 값 원문 보존. `Int64() (int64, error)` 는 앞뒤 공백·천 단위 콤마를
  지우고 파싱(`"          -2,218"` → -2218). `-`·빈 값은 `ErrNotNumber`. `String()` 은 원문.

## `customs` 패키지

- `New(c *opendata.Client) *Service`.
- 공통 기간: 모든 Params 는 `Start, End string`(YYYYMM). 메서드 진입 시 형식(`^\d{6}$`, 월 01~12),
  `Start ≤ End`, **시작·종료 포함 12개월 이하**를 검증해 서버 호출 없이 에러를 낸다(서버도 99 로 거부).
- 17개 메서드 — 이름은 아래 표. 각 메서드는 `func (s *Service) X(ctx, p XParams) ([]XRow, error)`.
- Row 타입 필드는 XML 태그 그대로 대응(Go 이름은 첫 글자 대문자). 단위가 있는 필드(건·kg·달러·천 달러·원)는
  `opendata.Num`, 나머지는 `string`. 필드 주석에 단위를 적는다(13~17 금액은 **천 달러**).
- 명세에만 있고 실응답엔 없는 필드(성질별 `statCd`·`statCdCntnKor1`)는 넣지 않는다.
- 총계 행: Row 마다 `IsTotal() bool` — 기간 필드(`Year`/`PriodTitle`) == `총계`. 라이브러리는 총계 행을
  **걸러내지 않는다**(원 응답 보존). 걸러내기는 호출자 몫이며 `customs.WithoutTotal(rows)` 제네릭 헬퍼를 둔다.
- 긴 기간: `customs.Years(start, end string) ([]Period, error)` — `[start,end]` 를 **달력 연도 경계**로
  나눈다(예: 202407~202606 → 202407~202412, 202501~202512, 202601~202606). 연 단위 합산 API(13~15)가
  연도 중간에서 잘리지 않게 하려는 것. 호출자가 창마다 메서드를 부르고 이어 붙인다. 17개 각각의 `...All`
  메서드는 두지 않는다(YAGNI — 반복문 세 줄이면 된다).

| # | 메서드 | 오퍼레이션 |
|---|---|---|
| 1 | `ItemCountryTrade` | `nitemtrade/getNitemtradeList` |
| 2 | `ItemTrade` | `Itemtrade/getItemtradeList` |
| 3 | `CountryTrade` | `nationtrade/getNationtradeList` |
| 4 | `ContinentTrade` | `continenttradet/getContinenttradeList` |
| 5 | `EconomicBlocTrade` | `economytrade/getEconomytradeList` |
| 6 | `PropertyTrade` | `Idfytempertrade/getIdfytempertradeList` |
| 7 | `PropertyCountryTrade` | `ntempertrade/getNtempertradeList` |
| 8 | `NewPropertyTrade` | `newtempertrade/getNewtempertradeList` |
| 9 | `NewPropertyCountryTrade` | `nnewtempertrade/getNnewtempertradeList` |
| 10 | `KindTrade` | `kindtrade/getKindtradeList` |
| 11 | `CustomsOfficeTrade` | `customstrade/getCustomstradeList` |
| 12 | `PortTrade` | `porttrade/getPorttradeList` |
| 13 | `SidoTrade` | `sidotrade/getSidotradeList` |
| 14 | `SidoItemTrade` | `sidoitemtrade/getSidoitemtradeList` |
| 15 | `SidoPropertyTrade` | `sidotempertrade/getSidotempertradeList` |
| 16 | `SigunguTrade` | `sigunguperimexacrs/getSigunguPerImexAcrs` |
| 17 | `SigunguItemTrade` | `sigunguperprlstperacrs/getSigunguPerPrlstPerAcrs` |

요청 인자: 필수 여부는 클라이언트가 검사하지 않는다(실측상 명세와 다른 경우가 많다 — 서버 99 에 맡긴다).
빈 문자열 인자는 쿼리에서 뺀다. 17번 `HsSgn` 은 명세대로 대문자로 보낸다.

Params·Row 의 필드 목록은 `docs/api/customs/<API>.md` 의 "요청 인자"·"응답 필드" 표를 그대로 옮긴다.
손으로 17개를 옮기면 오타가 나므로 **표에서 Go 코드를 생성하는 일회성 스크립트**로 만들고, 생성물을 커밋한다
(스크립트는 커밋하지 않는다 — 문서가 바뀔 일이 드물고, 생성물은 리뷰 대상이다).

## 테스트

- 루트: `httptest` 로 정상·`<items/>`·body 변형 3종·게이트웨이 401/403/400·서비스 99·HTML 응답·
  타임아웃 시 키 마스킹·`_type` 미포함·빈 인자 제외를 검증. `Num` 표 테스트.
- customs: 기간 검증 표 테스트, `Years` 표 테스트, 17개 API 각각 fixture 로 디코드해 행 수·첫 행 대표 필드·
  총계 판별을 확인(요청 경로와 쿼리도 검사).
- fixture: 1단계 실측 응답을 item 3개 + 총계 행으로 줄인 축약본(`testdata/customs/NN_*.xml`). 서비스키는 없다.
- `integration_test.go`(`-tags integration`): `OPENDATA_API_KEY` 로 17개를 202601~202601 한 달씩 호출해
  에러 없음·행 1개 이상(9번은 데이터가 드물어 0행 허용)을 확인. 8번(1개월 23,650행)은 `imexTmprUnfcClsfCd` 를 줘서 줄인다.

## 릴리스·연동

- `scripts/release.sh`: ecos-go 것을 복사해 이름만 바꾼다(모듈 zip 검증 포함 — 한글 파일명이 모듈 zip 규칙을
  통과하는지 이 단계에서 확인된다).
- README 에 사용 예·커버리지 표·단위/총계 주의. 워크스페이스 `CLAUDE.md` 의 opendata-go 행 갱신.
- PR → 사용자 머지 → `scripts/release.sh v0.1.0`. moneyflow 연동은 별도 스펙.

## 완료 기준

- `go build ./... && go vet ./... && go test ./...` 통과, `go test -tags integration ./...` 실 호출 통과
- 17개 메서드가 문서 표의 필드를 모두 갖는다(생성 스크립트가 보장, 테스트가 fixture 로 확인)
- 어떤 에러 문자열에도 serviceKey 가 없다(테스트)
