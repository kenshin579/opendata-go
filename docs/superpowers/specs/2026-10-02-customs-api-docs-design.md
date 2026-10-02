# opendata-go 설계 — 1단계: 공통 규약 + 관세청 수출입실적 API 명세 문서

- 날짜: 2026-10-02
- 상태: 승인됨 (저장소 이름·범용 구조 사용자 확인)
- 이번 작업 범위: **API 명세 문서(md) 작성**. 라이브러리 구현(2단계)은 문서 검수 후 별도 플랜.

## 배경

moneyflow 에 국내 수출 데이터 페이지를 만들고 싶다. 수출 통계를 REST 로 주는 곳은
공공데이터포털(data.go.kr)의 관세청 수출입실적 API 다. 관세청만 감싸는 SDK 대신,
포털 공통 부분(서비스키·게이트웨이·응답 봉투·에러·페이지네이션)을 한 번 만들고
기관별 서브패키지를 붙이는 **범용 SDK** `github.com/kenshin579/opendata-go` 로 간다
(사용자 결정). 첫 기관은 관세청이다.

자매 프로젝트 `ecos-go`·`opendart-go` 처럼 **명세를 md 로 먼저 정리하고, 그 문서를 근거로
구현**한다. 공식 명세는 포털의 Swagger·기술문서(docx)·`관세청조회코드_v1.3.xlsx` 에 흩어져
있어 크롤러 없이 **공식 명세 + 실 API 호출 응답을 대조해 직접 작성**한다.

## 패키지 구조 (2단계 기준, 문서 구조도 이를 따른다)

```
opendata/            # 루트 패키지: Client, 서비스키, 게이트웨이 호출, 봉투 파싱, 에러, 페이지 반복
opendata/customs/    # 관세청(기관코드 1220000) 수출입실적 17개
opendata/<기관>/      # 이후 필요할 때 추가
```

## 대상 API — 관세청 수출입실적(GW) 17개

포털 검색("관세청 수출입실적", 오픈API) 결과 17건 (2026-10-02 기준).

| # | 데이터명 | 데이터 ID | 문서 파일 |
|---|---------|-----------|-----------|
| 1 | 품목별 국가별 수출입실적 | 15100475 | `docs/api/customs/품목별국가별.md` |
| 2 | 품목별 수출입실적 | 15101609 | `docs/api/customs/품목별.md` |
| 3 | 국가별 수출입실적 | 15101612 | `docs/api/customs/국가별.md` |
| 4 | 대륙별 수출입실적 | 15101630 | `docs/api/customs/대륙별.md` |
| 5 | 경제권별 수출입실적 | 15101632 | `docs/api/customs/경제권별.md` |
| 6 | 성질별 수출입실적 | 15102109 | `docs/api/customs/성질별.md` |
| 7 | 성질별 국가별 수출입실적 | 15100476 | `docs/api/customs/성질별국가별.md` |
| 8 | 신성질별 수출입실적 | 15101616 | `docs/api/customs/신성질별.md` |
| 9 | 신성질별 국가별 수출입실적 | 15101607 | `docs/api/customs/신성질별국가별.md` |
| 10 | 종류별 수출입실적 | 15101634 | `docs/api/customs/종류별.md` |
| 11 | 세관별 수출입실적 | 15101633 | `docs/api/customs/세관별.md` |
| 12 | 항구 공항별 수출입실적 | 15101636 | `docs/api/customs/항구공항별.md` |
| 13 | 시도별 수출입실적 | 15101643 | `docs/api/customs/시도별.md` |
| 14 | 시도별 품목별 수출입실적 | 15101641 | `docs/api/customs/시도별품목별.md` |
| 15 | 시도별 성질별 수출입실적 | 15101639 | `docs/api/customs/시도별성질별.md` |
| 16 | 시군구별 수출입실적 | 15134344 | `docs/api/customs/시군구별.md` |
| 17 | 시군구별 품목별 수출입실적 | 15134343 | `docs/api/customs/시군구별품목별.md` |

### 포털 공식 명세 (2026-10-02 수집, 로그인 불필요)

API 상세 페이지는 두 형식이다 — 15개는 페이지에 **Swagger 2.0 JSON 이 인라인**(`const swaggerJson`)되어
있고, 2개(15100475·15100476)는 Swagger 없이 `POST /tcs/dss/selectApiDetailFunction.do` 가 HTML 표로
명세를 준다. 둘 다 로그인 없이 받을 수 있어 **수집 스크립트로 원본을 저장소에 남긴다**(재현 가능).
공통 코드표 `관세청조회코드_v1.3.xlsx`(시트 11개: 수출수입·품목·국가·성질분류·성질통합분류·대륙·
경제권·세관구분·수출입종류·항구공항·시도)도 로그인 없이 받아진다.

| # | 오퍼레이션 (`/1220000/...`) | 요청 인자 (`*` 필수, serviceKey 제외) |
|---|---|---|
| 1 | `nitemtrade/getNitemtradeList` | strtYymm* endYymm* hsSgn cntyCd* |
| 2 | `Itemtrade/getItemtradeList` | strtYymm* endYymm* hsSgn |
| 3 | `nationtrade/getNationtradeList` | strtYymm* endYymm* cntyCd |
| 4 | `continenttradet/getContinenttradeList` | strtYymm* endYymm* cntnEbkUnfcClsfCd |
| 5 | `economytrade/getEconomytradeList` | strtYymm* endYymm* cntnEbkUnfcClsfCd |
| 6 | `Idfytempertrade/getIdfytempertradeList` | strtYymm* endYymm* imexTpcd* imexTmprClsfCd |
| 7 | `ntempertrade/getNtempertradeList` | strtYymm* endYymm* imexTpcd* imexTmprClsfCd cntyCd* |
| 8 | `newtempertrade/getNewtempertradeList` | strtYymm* endYymm* imexTpcd* imexTmprUnfcClsfCd |
| 9 | `nnewtempertrade/getNnewtempertradeList` | strtYymm* endYymm* imexTpcd* imexTmprUnfcClsfCd* cntyCd* |
| 10 | `kindtrade/getKindtradeList` | strtYymm* endYymm* imexKcd imexTpcd* |
| 11 | `customstrade/getCustomstradeList` | strtYymm* endYymm* cstmSgnYn |
| 12 | `porttrade/getPorttradeList` | strtYymm* endYymm* portAirptRegnCd |
| 13 | `sidotrade/getSidotradeList` | strtYymm* endYymm* sidoCd |
| 14 | `sidoitemtrade/getSidoitemtradeList` | strtYymm* endYymm* sidoCd* |
| 15 | `sidotempertrade/getSidotempertradeList` | strtYymm* endYymm* dtlTmprYn sidoCd* imexTpcd* imexTmprClsfCd |
| 16 | `sigunguperimexacrs/getSigunguPerImexAcrs` | strtYymm* endYymm* sidoCd* |
| 17 | `sigunguperprlstperacrs/getSigunguPerPrlstPerAcrs` | strtYymm* endYymm* HsSgn* sidoCd* |

(경로 대소문자·오타(`continenttradet`)·`HsSgn` 대문자는 명세 원문 그대로다. 실측으로 확인한다.)

응답 필드는 **이름 체계가 두 계열**로 갈린다 — 2단계 타입 설계에 영향:

- 구 계열(1~12): `year`, `expDlr`/`impDlr`/`balPayments`, `expWgt`/`impWgt`, `expCnt`/`impCnt`,
  성질·종류 계열은 `impexp`+`dlr`/`wgt` 단일 값 행
- 시도·시군구 계열(13~17): `priodTitle`, `expUsdAmt`/`impUsdAmt`/`cmtrBlncAmt`, `expLnCnt` 등.
  16·17만 Swagger 에 `body.totalCount` 가 있다

Swagger 어디에도 `pageNo`·`numOfRows`·`_type` 이 없다 — **페이지네이션과 JSON 지원은 실측 전까지 모른다.**

기타:

- 기간 `strtYymm`·`endYymm` 은 YYYYMM, **조회기간 1년 이내**(명세 문구)
- 금액 USD — 수출 FOB(신고금액), 수입 CIF(과세가격). 중량 순중량 kg.
- 출력: XML(포털 표기). JSON(`_type=json` 등) 지원 여부는 실측으로 확인.
- 갱신: 매월 15일경 전월까지 현행화(정정·취하 반영). 최근월은 잠정치라 바뀔 수 있다.
- 트래픽: 개발계정 API별 1일 10,000 건.
- 공통 코드(국가·품목·성질·세관·시도 등)는 `관세청조회코드_v1.3.xlsx` 에 있다.

## 산출물 — 문서 구조

```
docs/api/
├── README.md                # 포털 공통 규약 (기관 무관)
└── customs/
    ├── README.md            # 관세청 인덱스 + 관세청 공통(코드표, 기간 제한, 금액 기준)
    ├── <17개 API>.md
    ├── codes/<시트명>.csv    # 관세청조회코드_v1.3.xlsx 시트 11개를 CSV 로 (2단계 상수·검증용)
    └── _source/             # 포털 원본: <데이터ID>.swagger.json 또는 <데이터ID>.detail.html, 코드표 xlsx
scripts/portal-spec/harvest.py   # _source/ 와 codes/ 를 다시 만드는 수집 스크립트 (python3 표준 라이브러리만)
```

### docs/api/README.md — 포털 공통 규약

- 서비스키: 발급, **인코딩 키 vs 디코딩 키**(쿼리 이중 인코딩 함정), API 별 활용신청 필요
- 게이트웨이 URL 규칙(`apis.data.go.kr/{기관코드}/...`), 공통 인자(`pageNo`, `numOfRows`, `_type` 등 실측된 것만)
- 정상 응답 봉투(`response/header/body/items/item`, `totalCount`) — 실측 XML 발췌
- **게이트웨이 레벨 에러**(`OpenAPI_ServiceResponse` / `returnReasonCode` 형식)와 서비스 레벨
  에러(`resultCode` ≠ 00)의 구분 — 둘은 구조가 다르다. 실제로 재현한 것(잘못된 키, 미신청 API,
  필수 인자 누락)은 실 응답을 싣는다
- 결과 0건일 때의 모양(`items` 빈 태그 vs 생략, `item` 단건일 때 배열 아님 등 XML 함정)
- `api.odcloud.kr`(파일데이터 API) 존재와 봉투 차이를 한 단락으로 기록만 — 이번 범위 밖

### 각 API 문서 포맷 (ecos-go `docs/api` 포맷 준용)

1. 제목 + 요청 URL 한 줄 요약 + 포털 데이터 ID 링크
2. **기본 정보**: 메서드 / URL / 출력 포맷 / 갱신 주기
3. **요청 인자**: 인자명 · 명칭 · 타입·길이 · 필수 · 값 설명(코드표 링크)
4. **응답 필드**: 필드명 · 명칭 · 단위 · 설명 — 실 응답과 1:1 대조해 전수
5. **샘플**: 실 호출 URL(서비스키는 `{SERVICE_KEY}`) + 실 응답 발췌
6. **함정**: 실측에서 발견한 것(기간 제한, 합계 행 존재 여부, 숫자 형식 등)

## 검증 방법

API 마다:

1. 실 API 호출로 응답 필드 전수 확인 — 문서 표와 실제 XML 태그 1:1 대조
2. 필수 인자 누락·기간 1년 초과 호출로 에러 모양 확인
3. 포털 Swagger·기술문서와 교차 확인

완료 후 체크리스트:

- [ ] `docs/api/customs/` 에 17개 문서 + README 가 있다
- [ ] 각 문서 요청 인자 표가 Swagger 와 일치한다
- [ ] 각 문서 응답 필드 표가 실 응답 태그와 전수 일치한다
- [ ] 공통 README 에 게이트웨이 에러·서비스 에러 실 응답이 모두 실려 있다
- [ ] 서비스키가 어떤 파일에도 들어가지 않았다 (`grep`)
- [ ] 인덱스 링크가 모두 유효하고 파일이 UTF-8 이다

## 선행 조건 (사용자 작업)

- data.go.kr 서비스키를 `~/.zshrc` 에 `OPENDATA_API_KEY` 로 둔다 (디코딩 키 권장 — 실측으로 확정)
- 17개 API 각각 **활용신청**(개발계정, 자동승인). 미신청 API 는 게이트웨이 에러로 막힌다

키가 없으면 문서 작업 중 Swagger 기반 표까지만 쓰고 실측 대조는 멈춘다.

## 완료 기준

- 위 체크리스트 전 항목 통과
- 사용자 검수 후 2단계(라이브러리 구현) 플랜으로 진행

## (참고) 2단계 — 라이브러리 설계 방향

문서 검수 후 별도 스펙·플랜으로 확정. 지금 합의된 방향만 적는다.

- 루트 `opendata`: `NewClient(serviceKey, opts...)` / `NewClientFromEnv()`(`OPENDATA_API_KEY`),
  옵션 `WithBaseURL`·`WithTimeout`(기본 30s)·`WithHTTPClient`. 게이트웨이 호출·봉투 파싱·에러 매핑·
  페이지 반복을 기관 패키지에 제공. **에러 메시지·URL 로그에서 serviceKey 마스킹**(쿼리스트링에 실리므로 필수)
- `customs`: `customs.New(client)` + API 17개 메서드. 기간이 1년을 넘으면 1년 창으로 나눠 이어 붙이는
  `...All` 헬퍼. 숫자 필드는 원본 문자열 보존 + 파싱 헬퍼
- 에러: 결과 없음 → `ErrNoData`, 게이트웨이 에러(키·한도) → `*GatewayError`, 서비스 에러 → `*APIError`
- 외부 의존성 0(표준 라이브러리, `encoding/xml`), Go 1.25, `httptest` + `testdata/` fixture,
  `-tags integration` 실 호출
- 릴리스: `scripts/release.sh` 재사용 → v0.1.0 → moneyflow `go.mod` 에 태그로 추가
- moneyflow 쪽(DB 테이블·월간 크론·수출 페이지)은 moneyflow 저장소의 별도 스펙
