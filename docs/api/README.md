# 공공데이터포털(data.go.kr) OpenAPI 공통 규약

공공데이터포털 게이트웨이(`apis.data.go.kr`)를 거치는 OpenAPI 의 공통 규약입니다. 기관과 무관하게 적용되는 내용만 다루고, 기관별 명세는 하위 디렉토리에 둡니다.

- [관세청 수출입실적](customs/README.md) (`customs/`)

포털 공지와 관세청 API 실 호출 결과를 대조해 작성했습니다(실측일 2026-10-02). 실측으로 확인하지 못한 내용은 "공지 기준(미실측)"이라고 적었습니다.

## 서비스키

- 포털 마이페이지에서 일반 인증키(서비스키)를 발급받는다. 키 하나로 여러 API 를 부를 수 있지만 **API 마다 활용신청을 따로 해야 한다.** 신청하지 않은 API 를 부르면 게이트웨이 에러 30 이 온다(실측, 아래 [에러](#에러) 참고).
- 개발계정 활용신청은 자동승인이고, 트래픽은 API 당 1일 10,000건이다.
- 쿼리 인자 `serviceKey` 로 전달한다.
- 포털은 같은 키를 Encoding·Decoding 두 형태로 보여준다. 키에 `+`·`/`·`=` 같은 문자가 있으면 이중 인코딩 문제가 생기는 것으로 알려져 있지만, **이번에 발급받은 키는 64자 16진수**라 URL 인코딩해도 문자열이 바뀌지 않는다. 한 번 더 인코딩해 보내도 정상 응답이었다(실측). 다른 형식의 키에서는 확인하지 않았다.

## 게이트웨이 URL

```
https://apis.data.go.kr/{기관코드}/{서비스}/{오퍼레이션}?serviceKey={SERVICE_KEY}&...
```

- 관세청 기관코드는 `1220000`. 예: `https://apis.data.go.kr/1220000/nationtrade/getNationtradeList`
- 명세(Swagger `schemes`)에는 `https`·`http` 둘 다 있다. `http://` 도 리다이렉트 없이 그대로 응답한다(실측).
- 메서드는 GET.

## 공통 인자

포털의 다른 API 에서 흔히 쓰는 `numOfRows`·`pageNo`·`_type` 은 **관세청 API 에서는 쓸 수 없다.**

| 인자 | 관세청 API 실측 결과 |
| --- | --- |
| `numOfRows`, `pageNo` | **무시된다. 페이지네이션이 없다.** 인자가 없을 때와 `numOfRows=5&pageNo=2`·`numOfRows=1000&pageNo=3` 의 응답이 바이트 단위로 같았다(시군구별 186행, 신성질별 23,650행·5,074,745바이트). Swagger 명세에도 이 인자들은 없다 |
| `_type=json` | **쓰면 안 된다.** 데이터 대신 게이트웨이 에러 04 가 JSON 으로 온다(아래) |
| `type=json`, `dataType=JSON`, `Accept: application/json` 헤더 | 무시되고 XML 이 온다 |

정리하면 **관세청 API 는 XML 전용이고 페이지네이션이 없다.** 한 번 호출하면 조건에 맞는 전체 행이 한 응답에 온다. 응답 크기를 줄이려면 기간·필터 인자를 좁혀야 한다(신성질별 6개월 조회는 142,386행·30.5MB).

`_type=json` 응답(HTTP 200, `Content-Type: application/json`):

```json
{
  "OpenAPI_ServiceResponse": {
    "cmmMsgHeader": {
      "errMsg": "HTTP_ERROR",
      "returnAuthMsg": "HTTP 에러",
      "returnReasonCode": "04"
    }
  }
}
```

## 정상 응답

`response/header/{resultCode,resultMsg}` + `response/body/items/item[]`. `resultCode` 는 `00`, `resultMsg` 는 `정상서비스.` 이다. HTTP 상태는 200, `Content-Type: application/xml`.

```xml
<?xml version="1.0" encoding="UTF-8" standalone="yes"?><response><header><resultCode>00</resultCode><resultMsg>정상서비스.</resultMsg></header><body><items><item><balPayments>5244595000</balPayments><expCnt>196748</expCnt><expDlr>12006733027</expDlr><impCnt>1433422</impCnt><impDlr>6762138027</impDlr><statCd>US</statCd><statCdCntnKor1>미국</statCdCntnKor1><year>2026.01</year></item></items></body></response>
```

- **단건**: 위처럼 `<item>` 이 하나 온다(배열 감싸기 없음).
- **0건**: 에러가 아니라 `resultCode` 00 에 `items` 가 자체닫힘 태그로 온다. 없는 코드로 조회해도 이 모양이다.

  ```xml
  <response><header><resultCode>00</resultCode><resultMsg>정상서비스.</resultMsg></header><body><items/></body></response>
  ```

- **`totalCount`**: 관세청 17개 중 시군구 계열 2개에만 있다(`...</items><totalCount>186</totalCount></body>`). 값은 item 수와 같고 페이지네이션과는 관계없다.
- 응답 헤더에 `X-RateLimit-Limit: 10000`·`X-RateLimit-Remaining`(남은 일일 호출 수)이 붙는다. 게이트웨이 에러(HTTP 401/403/400) 응답에는 없다. `_type=json` 에러 응답에는 붙어 있으므로 이것도 호출 수에서 차감되는 것으로 보인다.
- 이름 값의 `&` 는 XML 표준대로 `&amp;` 로 이스케이프되어 온다(`<statCdCntnKor1>남조지아 &amp; 남샌드위치 군도`). XML 파서로 읽으면 `&` 가 된다. 문자열 치환으로 파싱하면 `&amp;` 가 그대로 남는다.

## 에러

에러는 두 층에서 나오고 모양이 완전히 다르다. 클라이언트는 둘 다 처리해야 한다.

| 층 | 루트 요소 | HTTP 상태 | 판별 |
| --- | --- | --- | --- |
| (a) 게이트웨이 | `OpenAPI_ServiceResponse/cmmMsgHeader` | 401·403·400 (`_type=json` 은 200) | `returnReasonCode` |
| (b) 기관 서비스 | `response/header` | **200** | `resultCode` ≠ `00` |

HTTP 상태만 보면 (b) 를 놓친다. 응답 본문의 루트 요소로 먼저 가른다.

### (a) 게이트웨이 에러

잘못된 서비스키(HTTP 403):

```xml
<?xml version="1.0" encoding="UTF-8"?>
<OpenAPI_ServiceResponse>
<cmmMsgHeader>
  <errMsg>SERVICE_KEY_IS_NOT_REGISTERED_ERROR</errMsg>
  <returnAuthMsg>등록되지 않은 서비스키</returnAuthMsg>
  <returnReasonCode>30</returnReasonCode>
</cmmMsgHeader>
</OpenAPI_ServiceResponse>
```

| 코드 | errMsg | 의미 | HTTP | 비고 |
| --- | --- | --- | --- | --- |
| 1 | APPLICATION_ERROR | 어플리케이션 에러 | | 공지 기준(미실측) |
| 04 | HTTP_ERROR | HTTP 에러 | 200 | 실측: `_type=json` 을 붙였을 때. 본문이 JSON |
| 10 | INVALID_REQUEST_PARAMETER_ERROR | 잘못된 요청 파라미터 | | 공지 기준(미실측) |
| 12 | NO_OPENAPI_SERVICE_ERROR | 해당 오픈API 서비스가 없거나 폐기됨 | 400 | 실측: 없는 오퍼레이션 |
| 20 | SERVICE_ACCESS_DENIED_ERROR | 서비스 접근거부 | 401 | 실측: `serviceKey` 누락. 이때 errMsg 는 공지와 달리 `SERVICE_KEY_IS_NULL` |
| 22 | LIMITED_NUMBER_OF_SERVICE_REQUESTS_EXCEEDS_ERROR | 일일 트래픽 초과 | | 공지 기준(미실측) |
| 30 | SERVICE_KEY_IS_NOT_REGISTERED_ERROR | 등록되지 않은 서비스키 | 403 | 실측: 잘못된 키, **활용신청하지 않은 API 호출**도 같은 코드 |
| 31 | DEADLINE_HAS_EXPIRED_ERROR | 활용기간 만료 | | 공지 기준(미실측) |
| 32 | UNREGISTERED_IP_ERROR | 등록되지 않은 IP | | 공지 기준(미실측) |
| 99 | UNKNOWN_ERROR | 기타 에러 | | 공지 기준(미실측) |

### (b) 기관 서비스 에러

관세청 API 는 서비스 에러의 `resultCode` 가 **전부 `99`** 이고, 원인은 `resultMsg` 로만 구분된다. HTTP 상태는 200.

```xml
<?xml version="1.0" encoding="UTF-8" standalone="yes"?><response><header><resultCode>99</resultCode><resultMsg>시작과 종료의 조회기간은 1년이내 기간만 가능합니다.</resultMsg></header></response>
```

실측한 `resultMsg`:

| 원인 | resultMsg |
| --- | --- |
| 조회기간 12개월 초과 | `시작과 종료의 조회기간은 1년이내 기간만 가능합니다.` |
| 종료년월 < 시작년월 | `종료년월이 시작년월보다 크거나 같아야 합니다.` |
| 년월 형식 오류(`2026-01`, `2026.01`) | `시작년월 형식이 맞지 않습니다.` |
| 필수 인자 누락 | `필수 요청변수가 누락되었습니다.` |
| 없는 국가코드(소문자 `us` 포함) | `존재하지 않는 국가코드입니다.` |
| 품목코드 자릿수 오류(품목별국가별) | `품목코드는 2,4,6,10 자리로 입력해야 합니다.` |

`body` 의 모양이 API 마다 다르다. 세 가지가 나왔다.

| 모양 | 나온 곳(예) |
| --- | --- |
| `body` 없음 (`</header></response>`) | 국가별(필수 누락·기간 초과 등) |
| `<body/>` | 성질별·신성질별국가별(필수 누락), 품목별국가별(품목코드 자릿수 오류) |
| `<body><totalCount>0</totalCount></body>` | 시군구별(필수 누락) |

파서는 `body` 가 없거나 비어 있어도 실패하지 않아야 하고, 에러 판별은 `resultCode` 로만 한다.

## odcloud 파일데이터 API

포털의 "파일데이터"를 API 로 제공하는 `api.odcloud.kr` 은 `apis.data.go.kr` 과 다른 게이트웨이다. JSON 을 기본으로 쓰고 응답 봉투(필드 구성)도 이 문서의 `response/header/body` 와 다르다. 이 SDK 1단계 범위 밖이라 다루지 않는다.
