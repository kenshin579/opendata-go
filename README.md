# opendata-go

공공데이터포털([data.go.kr](https://www.data.go.kr)) OpenAPI 의 Go 클라이언트 라이브러리.
포털 공통(서비스키·게이트웨이·응답 봉투·에러)을 루트 패키지가 맡고, 기관별 API 는 서브패키지로 붙인다.

> 상태: 1단계 — API 명세 문서. 라이브러리 구현은 2단계에서 진행한다.

## 지원 기관

| 기관 | 서브패키지 | API | 문서 |
| --- | --- | --- | --- |
| 관세청 (1220000) | `customs` | 수출입실적 17개 | [docs/api/customs](docs/api/customs/README.md) |

## 문서

- [포털 공통 규약](docs/api/README.md) — 서비스키, 게이트웨이, 응답 봉투, 에러
- 포털 명세 원본·코드표 다시 받기: `python3 scripts/portal-spec/harvest.py`

## 인증

data.go.kr 에서 발급받은 서비스키를 `OPENDATA_API_KEY` 환경변수로 둔다.
API 마다 포털에서 활용신청이 필요하다.
