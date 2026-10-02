# 관세청 수출입실적 API 명세 문서 작성 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 공공데이터포털 공통 규약과 관세청 수출입실적 API 17개의 명세를 실 API 호출로 검증하며 `docs/api/` 문서(공통 README + 관세청 README + API 17개)로 작성한다.

**Architecture:** 포털 공식 명세(Swagger·상세 HTML)와 코드표 xlsx 를 수집 스크립트로 저장소에 원본 그대로 남기고(`docs/api/customs/_source/`, `codes/`), 실 호출 fixture(XML)를 모아 문서의 요청·응답 표를 원본·실측과 1:1 대조해 쓴다. 스펙: `docs/superpowers/specs/2026-10-02-customs-api-docs-design.md`

**Tech Stack:** Markdown, python3(표준 라이브러리), curl, git. Go 코드 없음(라이브러리 구현은 2단계 별도 플랜).

---

## 공통 준비 (모든 태스크에서 사용)

- 작업 디렉토리: `/Users/frankoh/src/workspace_moneyflow/opendata-go` (브랜치 `feature/customs-api-docs`)
- fixture 디렉토리(커밋 안 함): `/private/tmp/claude-501/-Users-frankoh-src-workspace-moneyflow/3fcaf621-532b-4e7d-a6c6-767706e37219/scratchpad/customs-fixtures`
- 서비스키 로드 (실 호출하는 모든 셸 명령 앞에서):

```bash
K=$(grep -m1 'export DATA_GO_KR_API_KEY' ~/.zshrc | cut -d= -f2- | tr -d '"' | tr -d "'")
FIX=/private/tmp/claude-501/-Users-frankoh-src-workspace-moneyflow/3fcaf621-532b-4e7d-a6c6-767706e37219/scratchpad/customs-fixtures
B=https://apis.data.go.kr/1220000
```

- **서비스키 마스킹 규칙**: 문서의 URL 예시는 키 자리를 `{SERVICE_KEY}` 로 쓴다. 커밋 전마다 `grep -rF "$K" docs/ scripts/` 가 아무것도 출력하지 않아야 한다.
- **실측 우선 규칙**: 이 플랜의 표는 포털 명세 기준이다. fixture 실측과 다르면 실측을 문서에 쓰고, 차이는 해당 문서 "함정" 절에 기록한다.
- **커밋**: 경로를 명령에 직접 넘긴다(`git add <경로>`). 메시지 끝에 아래 두 줄을 붙인다.

```
Co-Authored-By: Claude Opus 5.5 <noreply@anthropic.com>
Claude-Session: https://claude.ai/code/session_01Qw5e3X1n3xofs8Vm5BbW7W
```

## API 목록 (모든 문서 태스크가 참조)

| # | 데이터ID | 문서 | 오퍼레이션 (`$B/...`) | 원본 |
|---|---|---|---|---|
| 1 | 15100475 | 품목별국가별.md | `nitemtrade/getNitemtradeList` | detail.html |
| 2 | 15101609 | 품목별.md | `Itemtrade/getItemtradeList` | swagger |
| 3 | 15101612 | 국가별.md | `nationtrade/getNationtradeList` | swagger |
| 4 | 15101630 | 대륙별.md | `continenttradet/getContinenttradeList` | swagger |
| 5 | 15101632 | 경제권별.md | `economytrade/getEconomytradeList` | swagger |
| 6 | 15102109 | 성질별.md | `Idfytempertrade/getIdfytempertradeList` | swagger |
| 7 | 15100476 | 성질별국가별.md | `ntempertrade/getNtempertradeList` | detail.html |
| 8 | 15101616 | 신성질별.md | `newtempertrade/getNewtempertradeList` | swagger |
| 9 | 15101607 | 신성질별국가별.md | `nnewtempertrade/getNnewtempertradeList` | swagger |
| 10 | 15101634 | 종류별.md | `kindtrade/getKindtradeList` | swagger |
| 11 | 15101633 | 세관별.md | `customstrade/getCustomstradeList` | swagger |
| 12 | 15101636 | 항구공항별.md | `porttrade/getPorttradeList` | swagger |
| 13 | 15101643 | 시도별.md | `sidotrade/getSidotradeList` | swagger |
| 14 | 15101641 | 시도별품목별.md | `sidoitemtrade/getSidoitemtradeList` | swagger |
| 15 | 15101639 | 시도별성질별.md | `sidotempertrade/getSidotempertradeList` | swagger |
| 16 | 15134344 | 시군구별.md | `sigunguperimexacrs/getSigunguPerImexAcrs` | swagger |
| 17 | 15134343 | 시군구별품목별.md | `sigunguperprlstperacrs/getSigunguPerPrlstPerAcrs` | swagger |

## API 문서 템플릿 (Task 5~9 가 그대로 따른다)

````markdown
# <데이터명> (<오퍼레이션명>)

> `GET https://apis.data.go.kr/1220000/<서비스>/<오퍼레이션>`

<한두 문장 설명 — 포털 설명에서 요약>. 포털: [<데이터ID>](https://www.data.go.kr/data/<데이터ID>/openapi.do)

## 기본 정보

| 메서드 | 출력포맷 | 갱신 | 트래픽(개발) |
| --- | --- | --- | --- |
| GET | <실측: XML / XML·JSON> | 매월 15일경 전월까지 | 10,000/일 |

## 요청 인자

| 인자 | 명칭 | 크기 | 필수 | 값 설명 |
| --- | --- | --- | --- | --- |
| serviceKey | 서비스키 | 100 | Y | [공통 규약](../README.md) |
| strtYymm | 시작년월 | 6 | Y | YYYYMM. 조회기간 1년 이내 |
| ... | | | | 코드는 [codes/<시트>.csv](codes/<시트>.csv) |

## 응답 필드 (`response/body/items/item`)

| 필드 | 명칭 | 단위 | 설명 |
| --- | --- | --- | --- |
| ... | | | |

## 샘플

```
GET https://apis.data.go.kr/1220000/<서비스>/<오퍼레이션>?serviceKey={SERVICE_KEY}&strtYymm=...
```

```xml
<실 응답 발췌 — item 2~3개, 합계 행이 있으면 그것도>
```

## 함정

- <실측에서 발견한 것. 없으면 "실측에서 특이사항 없음">
````

응답 필드 표의 행 집합은 **원본 명세의 item 필드 집합 ∪ fixture 의 item 태그 집합**과 같아야 한다(Task 10 검증 스크립트가 확인).

---

### Task 1: 포털 명세 원본·코드표 수집 스크립트

**Files:**
- Create: `scripts/portal-spec/harvest.py`
- Create (스크립트 산출물): `docs/api/customs/_source/*.swagger.json` (15), `docs/api/customs/_source/*.detail.html` (2), `docs/api/customs/_source/관세청조회코드.xlsx`, `docs/api/customs/codes/*.csv` (11)

- [ ] **Step 1: 스크립트 작성** — 아래 내용 그대로 `scripts/portal-spec/harvest.py` 에 쓴다.

```python
#!/usr/bin/env python3
"""공공데이터포털 관세청 수출입실적 API 17개의 공식 명세 원본과 코드표를 수집한다.

로그인 불필요. python3 표준 라이브러리만 사용한다.
  - Swagger 가 페이지에 인라인된 API  → _source/<데이터ID>.swagger.json
  - Swagger 가 없는 API              → _source/<데이터ID>.detail.html (selectApiDetailFunction.do 응답)
  - 관세청조회코드 xlsx               → _source/관세청조회코드.xlsx + codes/<시트명>.csv

사용: python3 scripts/portal-spec/harvest.py [출력 디렉토리, 기본 docs/api/customs]
"""
import csv
import html
import io
import json
import re
import sys
import urllib.parse
import urllib.request
import zipfile
from pathlib import Path

PORTAL = "https://www.data.go.kr"
DATA_IDS = [
    "15100475", "15101609", "15101612", "15101630", "15101632", "15102109",
    "15100476", "15101616", "15101607", "15101634", "15101633", "15101636",
    "15101643", "15101641", "15101639", "15134344", "15134343",
]
UA = {"User-Agent": "opendata-go-harvest/1.0"}


def fetch(url, data=None):
    body = urllib.parse.urlencode(data).encode() if data else None
    req = urllib.request.Request(url, data=body, headers=UA)
    with urllib.request.urlopen(req, timeout=60) as r:
        return r.read()


def harvest_spec(data_id, src):
    page = fetch(f"{PORTAL}/data/{data_id}/openapi.do").decode("utf-8")
    title = html.unescape(re.search(r"<title>(.*?)</title>", page, re.S).group(1).split("|")[0].strip())
    sj = re.search(r"const swaggerJson = `(.*?)`;", page, re.S)
    if sj and sj.group(1).strip():
        spec = json.loads(sj.group(1))
        (src / f"{data_id}.swagger.json").write_text(json.dumps(spec, ensure_ascii=False, indent=2) + "\n")
        return title, "swagger"
    pk = re.search(r'id="publicDataDetailPk" value="([^"]+)"', page).group(1)
    sel = re.search(r'id="open_api_detail_select".*?</select>', page, re.S).group(0)
    seqs = re.findall(r'<option value="(\d+)"', sel)
    if len(seqs) != 1:
        raise SystemExit(f"{data_id}: 오퍼레이션이 {len(seqs)}개 — 스크립트 확장 필요")
    detail = fetch(f"{PORTAL}/tcs/dss/selectApiDetailFunction.do",
                   {"oprtinSeqNo": seqs[0], "publicDataDetailPk": pk, "publicDataPk": data_id})
    (src / f"{data_id}.detail.html").write_bytes(detail)
    return title, "detail"


def find_code_file(data_id):
    page = fetch(f"{PORTAL}/data/{data_id}/openapi.do").decode("utf-8")
    m = re.search(r"fn_fileDownload\('(FILE_\d+)','(\d+)'\)", page)
    return m.groups() if m else None


def xlsx_to_csv(xlsx_bytes, out_dir):
    z = zipfile.ZipFile(io.BytesIO(xlsx_bytes))
    shared = []
    if "xl/sharedStrings.xml" in z.namelist():
        for si in re.findall(r"<si>(.*?)</si>", z.read("xl/sharedStrings.xml").decode(), re.S):
            shared.append(html.unescape("".join(re.findall(r"<t[^>]*>(.*?)</t>", si, re.S))))
    wb = z.read("xl/workbook.xml").decode()
    rels = z.read("xl/_rels/workbook.xml.rels").decode()
    target = dict(re.findall(r'<Relationship [^>]*Id="([^"]+)"[^>]*Target="([^"]+)"', rels))
    target.update({k: v for v, k in re.findall(r'<Relationship [^>]*Target="([^"]+)"[^>]*Id="([^"]+)"', rels)})
    names = []
    for attrs in re.findall(r"<sheet ([^>]*)/>", wb):
        name = html.unescape(re.search(r'name="([^"]+)"', attrs).group(1))
        rid = re.search(r'r:id="([^"]+)"', attrs).group(1)
        xml = z.read("xl/" + target[rid].lstrip("/").removeprefix("xl/")).decode()
        rows = []
        for row in re.findall(r"<row[^>]*>(.*?)</row>", xml, re.S):
            cells = {}
            for attrs_c, inner in re.findall(r"<c ([^>]*?)(?:/>|>(.*?)</c>)", row, re.S):
                ref = re.search(r'r="([A-Z]+)\d+"', attrs_c).group(1)
                col = 0
                for ch in ref:
                    col = col * 26 + (ord(ch) - 64)
                v = re.search(r"<v>(.*?)</v>", inner or "", re.S)
                t = re.search(r't="(\w+)"', attrs_c)
                if t and t.group(1) == "s" and v:
                    val = shared[int(v.group(1))]
                elif t and t.group(1) == "inlineStr":
                    val = html.unescape("".join(re.findall(r"<t[^>]*>(.*?)</t>", inner, re.S)))
                else:
                    val = html.unescape(v.group(1)) if v else ""
                cells[col] = val
            if cells:
                rows.append([cells.get(c, "") for c in range(1, max(cells) + 1)])
        with open(out_dir / f"{name}.csv", "w", newline="", encoding="utf-8") as f:
            csv.writer(f).writerows(rows)
        names.append((name, len(rows)))
    return names


def main():
    out = Path(sys.argv[1] if len(sys.argv) > 1 else "docs/api/customs")
    src, codes = out / "_source", out / "codes"
    src.mkdir(parents=True, exist_ok=True)
    codes.mkdir(parents=True, exist_ok=True)
    for data_id in DATA_IDS:
        title, kind = harvest_spec(data_id, src)
        print(f"{data_id}\t{kind}\t{title}")
    file_id, sn = find_code_file("15101609")
    xlsx = fetch(f"{PORTAL}/cmm/cmm/fileDownload.do?atchFileId={file_id}&fileDetailSn={sn}")
    (src / "관세청조회코드.xlsx").write_bytes(xlsx)
    for name, n in xlsx_to_csv(xlsx, codes):
        print(f"codes\t{name}\t{n} rows")


if __name__ == "__main__":
    main()
```

- [ ] **Step 2: 실행**

Run: `cd /Users/frankoh/src/workspace_moneyflow/opendata-go && python3 scripts/portal-spec/harvest.py`

Expected: 17줄(`15100475	detail	...`, `15101609	swagger	...` …) 뒤에 `codes	<시트명>	<n> rows` 11줄. 시트: 수출수입코드·품목코드·국가코드·성질분류코드·성질통합분류코드·대륙코드·경제권코드·세관구분코드·수출입종류코드·항구공항코드·시도코드.

- [ ] **Step 3: 산출물 확인**

Run: `ls docs/api/customs/_source | wc -l; ls docs/api/customs/codes | wc -l; grep -c '' docs/api/customs/codes/시도코드.csv; du -sh docs/api/customs`
Expected: `18`, `11`, `21`, 용량 수 MB 이내(품목코드·성질통합분류코드 CSV 가 큼 — 그대로 커밋한다).

- [ ] **Step 4: 커밋**

```bash
git add scripts/portal-spec/harvest.py docs/api/customs/_source docs/api/customs/codes
git commit -m "docs: 관세청 수출입실적 포털 명세 원본·코드표 수집 스크립트 추가" -m "Co-Authored-By: Claude Opus 5.5 <noreply@anthropic.com>
Claude-Session: https://claude.ai/code/session_01Qw5e3X1n3xofs8Vm5BbW7W"
```

---

### Task 2: 실 호출 fixture 수집 (서비스키 필요)

**Files:** `$FIX/*.xml`, `$FIX/*.headers` (커밋 안 함)

- [ ] **Step 1: 선행 조건 확인**

Run: `grep -c 'export DATA_GO_KR_API_KEY' ~/.zshrc`
Expected: `1`. **0 이면 여기서 멈추고 사용자에게 보고한다** — data.go.kr 서비스키 발급과 17개 API 활용신청이 필요하다(스펙 "선행 조건"). Task 3 이후 중 원본만으로 쓸 수 있는 부분(요청 인자 표·응답 필드 초안)은 진행할 수 있지만 "샘플"·"함정"·출력포맷은 비워 두지 말고 이 태스크가 끝난 뒤 채운다.

- [ ] **Step 2: 정상 케이스 17개** — 기간은 `202601`~`202606`(1년 이내). 샘플 코드: 국가 `US`/`CN`, HS `8542`(전자집적회로), HS6 `854232`(메모리), 시도 `41`(경기도), 수출 `imexTpcd=1`.

```bash
mkdir -p "$FIX"
P="strtYymm=202601&endYymm=202606"
get() { curl -s -D "$FIX/$1.headers" "$B/$2?serviceKey=$K&$P&$3" -o "$FIX/$1.xml"; }
get 01_nitemtrade       nitemtrade/getNitemtradeList                  "hsSgn=8542&cntyCd=CN"
get 02_itemtrade        Itemtrade/getItemtradeList                    "hsSgn=8542"
get 03_nationtrade      nationtrade/getNationtradeList                "cntyCd=US"
get 04_continent        continenttradet/getContinenttradeList         ""
get 05_economy          economytrade/getEconomytradeList              ""
get 06_idfytemper       Idfytempertrade/getIdfytempertradeList        "imexTpcd=1"
get 07_ntemper          ntempertrade/getNtempertradeList              "imexTpcd=1&cntyCd=US"
get 08_newtemper        newtempertrade/getNewtempertradeList          "imexTpcd=1"
# 성질통합분류코드.csv: 머리말 3행 뒤 데이터, 15번째 열 = 신성질별 세세분류코드(예: 11020101)
UNFC=$(python3 -c "import csv;print(list(csv.reader(open('docs/api/customs/codes/성질통합분류코드.csv')))[3][14])")
get 09_nnewtemper       nnewtempertrade/getNnewtempertradeList        "imexTpcd=1&imexTmprUnfcClsfCd=$UNFC&cntyCd=US"
get 10_kind             kindtrade/getKindtradeList                    "imexTpcd=1"
get 11_customs          customstrade/getCustomstradeList              ""
get 12_port             porttrade/getPorttradeList                    ""
get 13_sido             sidotrade/getSidotradeList                    ""
get 14_sidoitem         sidoitemtrade/getSidoitemtradeList            "sidoCd=41"
get 15_sidotemper       sidotempertrade/getSidotempertradeList        "sidoCd=41&imexTpcd=1"
get 16_sigungu          sigunguperimexacrs/getSigunguPerImexAcrs      "sidoCd=41"
get 17_sigunguitem      sigunguperprlstperacrs/getSigunguPerPrlstPerAcrs "sidoCd=41&HsSgn=854232"
for f in "$FIX"/[01]*.xml; do printf '%s\t%s\t%s\n' "$(basename "$f")" "$(grep -o '<resultCode>[^<]*' "$f" | head -1)" "$(grep -o '<item>' "$f" | wc -l)"; done
```

Expected: 17줄 모두 `<resultCode>00` 이고 item 수 > 0. item 이 0 이면 인자(코드·기간)를 코드표에서 다른 값으로 바꿔 다시 받는다. resultCode 가 없고 `OpenAPI_ServiceResponse` 가 있으면 게이트웨이 에러다 — `returnReasonCode` 를 보고 활용신청 누락(보통 `SERVICE_KEY_IS_NOT_REGISTERED`)이면 사용자에게 해당 API 활용신청을 요청한다.

- [ ] **Step 3: 공통 규약 실측**

```bash
# a) JSON 지원 여부
curl -s "$B/nationtrade/getNationtradeList?serviceKey=$K&$P&cntyCd=US&_type=json" -o "$FIX/c_type_json.out"
curl -s "$B/nationtrade/getNationtradeList?serviceKey=$K&$P&cntyCd=US&type=json"  -o "$FIX/c_type_json2.out"
# b) 페이지네이션 인자 반응 여부 (전체 행이 가장 많은 시군구로)
curl -s "$B/sigunguperimexacrs/getSigunguPerImexAcrs?serviceKey=$K&$P&sidoCd=41" -o "$FIX/c_page_all.xml"
curl -s "$B/sigunguperimexacrs/getSigunguPerImexAcrs?serviceKey=$K&$P&sidoCd=41&numOfRows=5&pageNo=2" -o "$FIX/c_page_5_2.xml"
# c) 기간 1년 초과
curl -s "$B/nationtrade/getNationtradeList?serviceKey=$K&strtYymm=202401&endYymm=202606&cntyCd=US" -o "$FIX/c_over_1y.xml"
# d) 필수 인자 누락
curl -s "$B/nationtrade/getNationtradeList?serviceKey=$K&strtYymm=202601" -o "$FIX/c_missing.xml"
# e) 잘못된 서비스키
curl -s "$B/nationtrade/getNationtradeList?serviceKey=WRONGKEY&$P&cntyCd=US" -o "$FIX/c_bad_key.xml"
# f) 서비스키 인코딩: 키를 한 번 더 URL 인코딩해서 보내기
KE=$(python3 -c 'import sys,urllib.parse;print(urllib.parse.quote(sys.argv[1],safe=""))' "$K")
curl -s "$B/nationtrade/getNationtradeList?serviceKey=$KE&$P&cntyCd=US" -o "$FIX/c_key_encoded.xml"
# g) 결과 0건 (미래 기간)
curl -s "$B/nationtrade/getNationtradeList?serviceKey=$K&strtYymm=203001&endYymm=203001&cntyCd=US" -o "$FIX/c_empty.xml"
# h) 단건 결과 모양
curl -s "$B/nationtrade/getNationtradeList?serviceKey=$K&strtYymm=202601&endYymm=202601&cntyCd=US" -o "$FIX/c_single.xml"
for f in "$FIX"/c_*; do echo "== $(basename "$f")"; head -c 600 "$f"; echo; done
```

Expected: 파일 10개. 각 결과를 아래 질문에 대한 답으로 메모해 둔다(Task 3 에서 사용): JSON 이 되는가(어느 인자로)? `numOfRows`/`pageNo` 가 행 수를 바꾸는가, `totalCount` 가 있는가? 1년 초과·필수 누락은 어떤 코드·메시지인가? 잘못된 키의 응답 구조(게이트웨이 `OpenAPI_ServiceResponse` 인가)? 키를 이중 인코딩하면 실패하는가? 0건일 때 `items` 는 빈 태그인가 생략인가? 단건일 때도 `<item>` 하나인가?

- [ ] **Step 4: 합계 행·숫자 형식 확인**

Run: `for f in "$FIX"/[01]*.xml; do echo "== $(basename "$f")"; grep -o '<item>.*' "$f" | head -c 0; python3 -c "import re,sys;t=open(sys.argv[1]).read();it=re.findall(r'<item>(.*?)</item>',t,re.S);print(len(it));print(it[0][:400] if it else '');print(it[-1][:400] if it else '')" "$f"; done`
Expected: 각 API 의 첫·마지막 item. 마지막 행이 `총계`/`합계` 같은 합계 행인지, 기간 표기(`2026.01` vs `202601`), 숫자에 콤마·소수점이 있는지를 메모한다(각 API 문서 "함정").

---

### Task 3: 포털 공통 규약 문서

**Files:**
- Create: `docs/api/README.md`

- [ ] **Step 1: 문서 작성** — 아래 절을 이 순서로 쓴다. 각 절의 내용은 Task 2 Step 3 메모와 fixture 원문에서 가져온다.

1. 제목 `# 공공데이터포털(data.go.kr) OpenAPI 공통 규약` + 한 단락: 이 문서는 기관 무관 공통, 기관별 문서는 하위 디렉토리(`customs/` → [관세청](customs/README.md))
2. `## 서비스키` — 마이페이지에서 발급, 일반 인증키가 Encoding·Decoding 두 개로 보이는 점, **어느 것을 쿼리에 넣어야 하는지 실측 결과**(c_key_encoded), API 마다 활용신청 필요·개발계정 자동승인·트래픽 1일 10,000
3. `## 게이트웨이 URL` — `https://apis.data.go.kr/{기관코드}/{서비스}/{오퍼레이션}`, 관세청 기관코드 `1220000`, HTTP/HTTPS 둘 다 명세에 있음
4. `## 공통 인자` — 실측으로 반응이 확인된 것만(`numOfRows`, `pageNo`, `_type` 등). 반응 없으면 "관세청 API 는 무시함(실측)"이라고 쓴다
5. `## 정상 응답` — `response/header/{resultCode,resultMsg}` + `response/body/items/item[]` (+ `totalCount` 유무) 실 XML 발췌. 0건(c_empty)·단건(c_single) 모양
6. `## 에러` — 두 층을 표로: (a) 게이트웨이 에러 `OpenAPI_ServiceResponse/cmmMsgHeader/{errMsg,returnAuthMsg,returnReasonCode}` — 실측한 것(c_bad_key) 원문 + 포털이 공지한 코드표(1 APPLICATION_ERROR, 10 INVALID_REQUEST_PARAMETER_ERROR, 12 NO_OPENAPI_SERVICE_ERROR, 20 SERVICE_ACCESS_DENIED_ERROR, 22 LIMITED_NUMBER_OF_SERVICE_REQUESTS_EXCEEDS_ERROR, 30 SERVICE_KEY_IS_NOT_REGISTERED_ERROR, 31 DEADLINE_HAS_EXPIRED_ERROR, 32 UNREGISTERED_IP_ERROR, 99 UNKNOWN_ERROR — 실측 못 한 코드는 "공지 기준"이라 표시), (b) 서비스 에러 `response/header/resultCode` ≠ `00` — 실측한 것(c_over_1y, c_missing) 원문. HTTP 상태 코드도 함께(`$FIX/*.headers` 참고)
7. `## odcloud 파일데이터 API` — `api.odcloud.kr` 은 별도 게이트웨이·JSON·다른 봉투라는 것만 한 단락. "이 SDK 1단계 범위 밖"

- [ ] **Step 2: 마스킹 확인 후 커밋**

```bash
grep -rF "$K" docs/ scripts/ && echo "KEY LEAK" || echo ok
git add docs/api/README.md
git commit -m "docs: 공공데이터포털 OpenAPI 공통 규약 문서 추가" -m "Co-Authored-By: Claude Opus 5.5 <noreply@anthropic.com>
Claude-Session: https://claude.ai/code/session_01Qw5e3X1n3xofs8Vm5BbW7W"
```
Expected: `ok`

---

### Task 4: 관세청 인덱스 문서

**Files:**
- Create: `docs/api/customs/README.md`

- [ ] **Step 1: 문서 작성** — 절 순서:

1. `# 관세청 수출입실적 OpenAPI` + 기관코드 `1220000`, 공통 규약 링크(`../README.md`)
2. `## API 목록` — 이 플랜의 "API 목록" 표에서 `원본` 열을 빼고 문서 열을 링크로(`[품목별국가별](품목별국가별.md)`)
3. `## 공통 사항` — 기간 `strtYymm`/`endYymm` YYYYMM·1년 이내(실측 c_over_1y 결과), 금액 USD(수출 FOB 신고금액·수입 CIF 과세가격), 중량 순중량 kg, 갱신 매월 15일경 전월까지·최근월 변동 가능
4. `## 응답 필드 계열` — 구 계열(1~12: `year`, `expDlr`/`impDlr`/`balPayments` 등) vs 시도·시군구 계열(13~17: `priodTitle`, `expUsdAmt`/`impUsdAmt`/`cmtrBlncAmt` 등)을 대응표로. 같은 의미의 필드를 한 행에 놓는다(예: `expDlr` ↔ `expUsdAmt`, `balPayments` ↔ `cmtrBlncAmt`, `year` ↔ `priodTitle`, `expCnt` ↔ `expLnCnt`?) — 의미가 같은지는 fixture 값으로 확인하고, 확인 못 한 행은 표에 "추정"이라 쓴다
5. `## 코드표` — `codes/` CSV 11개 표(파일 · 쓰는 인자 · 행 수). 출처 `_source/관세청조회코드.xlsx`(포털 첨부 v1.3), 다시 받는 법 `python3 scripts/portal-spec/harvest.py`. CSV 맨 위 2~3행이 항목명 머리말이라는 것도 적는다
6. `## 코드표 함정` — `시도코드.csv` 의 `12 전남광주통합특별시`(2026.7.1. 이후 신고건 조회용, 비고 원문 인용)처럼 비고 열에 적힌 것은 전부 옮긴다: `grep -h ',[^,]\+$' docs/api/customs/codes/*.csv | head -50` 로 비고가 있는 행을 훑는다

- [ ] **Step 2: 커밋**

```bash
git add docs/api/customs/README.md
git commit -m "docs: 관세청 수출입실적 API 인덱스·공통 사항·코드표 문서 추가" -m "Co-Authored-By: Claude Opus 5.5 <noreply@anthropic.com>
Claude-Session: https://claude.ai/code/session_01Qw5e3X1n3xofs8Vm5BbW7W"
```

---

### Task 5: 품목·국가 계열 문서 (1·2·3)

**Files:**
- Create: `docs/api/customs/품목별국가별.md`, `docs/api/customs/품목별.md`, `docs/api/customs/국가별.md`

- [ ] **Step 1: 원본 필드 확인**

Run:
```bash
python3 - <<'EOF'
import json,re,html
for i in ["15101609","15101612"]:
    s=json.load(open(f"docs/api/customs/_source/{i}.swagger.json"))
    for path,ops in s["paths"].items():
        print(i,path,[(p["name"],p.get("required",False),p.get("description")) for p in ops.get("parameters",[])])
        def walk(x,pre=""):
            for k,v in (x.get("properties") or {}).items():
                if v.get("properties"): walk(v,pre+k+".")
                else: print("   ",pre+k,"|",v.get("description"),"|",v.get("example"))
        walk(ops["get"]["responses"]["200"]["schema"])
t=open("docs/api/customs/_source/15100475.detail.html").read()
t=re.sub(r"<[^>]+>"," ",t); print(re.sub(r"\s+"," ",html.unescape(t))[:2500])
EOF
```
Expected: 3개 API 의 요청 인자(이름·필수·설명)와 응답 필드(설명·예시). 15100475 는 표 텍스트(국문·영문·크기·구분·샘플·설명 순).

- [ ] **Step 2: 문서 3개 작성** — 템플릿대로. 요청 인자: 1 = `strtYymm* endYymm* hsSgn cntyCd*`, 2 = `strtYymm* endYymm* hsSgn`, 3 = `strtYymm* endYymm* cntyCd`. 응답 필드: 1 = `year statCdCntnKor1 statCd statKor hsCd expWgt expDlr impWgt impDlr balPayments`, 2 = `year statKor hsCode expWgt expDlr impWgt impDlr balPayments`, 3 = `year statCd statCdCntnKor1 expCnt expDlr impCnt impDlr balPayments`. 샘플은 fixture `01_nitemtrade.xml`·`02_itemtrade.xml`·`03_nationtrade.xml`. 함정 후보: 1 은 `hsCd`, 2 는 `hsCode` 로 같은 뜻의 필드명이 다르다(실측으로 확인); `hsSgn` 자릿수(2·4·6·10)에 따라 행이 어떻게 나오는지; 합계 행 여부.

- [ ] **Step 3: 필드 대조** — Task 10 Step 1 스크립트를 이 3개 문서에만 돌려(`python3 $SCRATCH/check_fields.py 1 2 3`, 스크립트는 Task 10 에 있음 — 이 태스크에서 먼저 만들어 둔다) `MISSING`/`EXTRA` 가 없어야 한다.

- [ ] **Step 4: 커밋**

```bash
grep -rF "$K" docs/ && echo "KEY LEAK" || echo ok
git add docs/api/customs/품목별국가별.md docs/api/customs/품목별.md docs/api/customs/국가별.md
git commit -m "docs: 관세청 품목별·품목별국가별·국가별 수출입실적 API 문서 추가" -m "Co-Authored-By: Claude Opus 5.5 <noreply@anthropic.com>
Claude-Session: https://claude.ai/code/session_01Qw5e3X1n3xofs8Vm5BbW7W"
```

---

### Task 6: 대륙·경제권 계열 문서 (4·5)

**Files:**
- Create: `docs/api/customs/대륙별.md`, `docs/api/customs/경제권별.md`

- [ ] **Step 1: 원본 확인** — Task 5 Step 1 스크립트에서 리스트를 `["15101630","15101632"]` 로, detail.html 출력 줄은 지운다.
- [ ] **Step 2: 문서 2개 작성** — 요청 인자 둘 다 `strtYymm* endYymm* cntnEbkUnfcClsfCd`(코드는 4 → `codes/대륙코드.csv`, 5 → `codes/경제권코드.csv`). 응답 필드 둘 다 `year statCd statCdCntnKor1 expCnt expDlr impCnt impDlr balPayments`. 샘플 `04_continent.xml`·`05_economy.xml`. 함정 후보: 경로 `continenttradet` 의 끝 `t`(명세 원문 그대로 호출해 성공했는지 기록), 대륙·경제권이 같은 인자명을 쓰는데 코드표가 다르다는 점, 코드 생략 시 전체 대륙/경제권이 나오는지.
- [ ] **Step 3: 필드 대조** — `python3 $SCRATCH/check_fields.py 4 5` → 출력에 `MISSING`/`EXTRA` 없음.
- [ ] **Step 4: 커밋**

```bash
grep -rF "$K" docs/ && echo "KEY LEAK" || echo ok
git add docs/api/customs/대륙별.md docs/api/customs/경제권별.md
git commit -m "docs: 관세청 대륙별·경제권별 수출입실적 API 문서 추가" -m "Co-Authored-By: Claude Opus 5.5 <noreply@anthropic.com>
Claude-Session: https://claude.ai/code/session_01Qw5e3X1n3xofs8Vm5BbW7W"
```

---

### Task 7: 성질 계열 문서 (6·7·8·9)

**Files:**
- Create: `docs/api/customs/성질별.md`, `docs/api/customs/성질별국가별.md`, `docs/api/customs/신성질별.md`, `docs/api/customs/신성질별국가별.md`

- [ ] **Step 1: 원본 확인** — Task 5 Step 1 스크립트에서 리스트를 `["15102109","15101616","15101607"]`, detail.html 경로를 `15100476.detail.html` 로.
- [ ] **Step 2: 문서 4개 작성** — 요청 인자: 6 = `strtYymm* endYymm* imexTpcd* imexTmprClsfCd`, 7 = `strtYymm* endYymm* imexTpcd* imexTmprClsfCd cntyCd*`, 8 = `strtYymm* endYymm* imexTpcd* imexTmprUnfcClsfCd`, 9 = `strtYymm* endYymm* imexTpcd* imexTmprUnfcClsfCd* cntyCd*`. 코드: `imexTpcd` → `codes/수출수입코드.csv`(1 수출·2 수입), `imexTmprClsfCd` → `codes/성질분류코드.csv`, `imexTmprUnfcClsfCd` → `codes/성질통합분류코드.csv`. 응답 필드 4개 모두 `year impexp statCd statCdCntnKor1 godsCd godsKor wgt dlr`(6·8 의 `statCd`/`statCdCntnKor1` 이 국가가 아니라면 무엇인지 fixture 로 확인해 명칭에 쓴다). 샘플 `06_idfytemper.xml`~`09_nnewtemper.xml`. 함정 후보: "성질별"(구 분류)과 "신성질별"(2012 신설 분류)의 차이를 각 문서 첫 단락에 쓴다; 성질통합분류코드 CSV 가 HS10 → 신성질 대·중·소·세·세세분류 매핑 표라서 인자로 넣을 코드가 어느 열인지(실측 성공한 `$UNFC` 가 몇 번째 열에서 왔는지) 명시; 9 는 성질코드가 필수.
- [ ] **Step 3: 필드 대조** — `python3 $SCRATCH/check_fields.py 6 7 8 9` → `MISSING`/`EXTRA` 없음.
- [ ] **Step 4: 커밋**

```bash
grep -rF "$K" docs/ && echo "KEY LEAK" || echo ok
git add docs/api/customs/성질별.md docs/api/customs/성질별국가별.md docs/api/customs/신성질별.md docs/api/customs/신성질별국가별.md
git commit -m "docs: 관세청 성질별·신성질별(국가별 포함) 수출입실적 API 문서 추가" -m "Co-Authored-By: Claude Opus 5.5 <noreply@anthropic.com>
Claude-Session: https://claude.ai/code/session_01Qw5e3X1n3xofs8Vm5BbW7W"
```

---

### Task 8: 종류·세관·항구 계열 문서 (10·11·12)

**Files:**
- Create: `docs/api/customs/종류별.md`, `docs/api/customs/세관별.md`, `docs/api/customs/항구공항별.md`

- [ ] **Step 1: 원본 확인** — Task 5 Step 1 스크립트에서 리스트를 `["15101634","15101633","15101636"]`, detail.html 줄 삭제.
- [ ] **Step 2: 문서 3개 작성** — 요청 인자: 10 = `strtYymm* endYymm* imexKcd imexTpcd*`(→ `codes/수출입종류코드.csv`), 11 = `strtYymm* endYymm* cstmSgnYn`(→ `codes/세관구분코드.csv`; Y/N 의미를 fixture 로 확인), 12 = `strtYymm* endYymm* portAirptRegnCd`(→ `codes/항구공항코드.csv`). 응답 필드: 10 = `year impexp statCd statCdCntnKor1 cnt wgt dlr won`(원화 금액 `won` 이 있는 유일한 API — 함정에 기록), 11 = `year center cstm statCdCntnKor expCnt expDlr impCnt impDlr balPayments`, 12 = `year statKor portCd cstmSgn expCnt expDlr impCnt impDlr balPayments`. 샘플 `10_kind.xml`~`12_port.xml`.
- [ ] **Step 3: 필드 대조** — `python3 $SCRATCH/check_fields.py 10 11 12` → `MISSING`/`EXTRA` 없음.
- [ ] **Step 4: 커밋**

```bash
grep -rF "$K" docs/ && echo "KEY LEAK" || echo ok
git add docs/api/customs/종류별.md docs/api/customs/세관별.md docs/api/customs/항구공항별.md
git commit -m "docs: 관세청 종류별·세관별·항구공항별 수출입실적 API 문서 추가" -m "Co-Authored-By: Claude Opus 5.5 <noreply@anthropic.com>
Claude-Session: https://claude.ai/code/session_01Qw5e3X1n3xofs8Vm5BbW7W"
```

---

### Task 9: 시도·시군구 계열 문서 (13~17)

**Files:**
- Create: `docs/api/customs/시도별.md`, `docs/api/customs/시도별품목별.md`, `docs/api/customs/시도별성질별.md`, `docs/api/customs/시군구별.md`, `docs/api/customs/시군구별품목별.md`

- [ ] **Step 1: 원본 확인** — Task 5 Step 1 스크립트에서 리스트를 `["15101643","15101641","15101639","15134344","15134343"]`, detail.html 줄 삭제.
- [ ] **Step 2: 문서 5개 작성** — 요청 인자: 13 = `strtYymm* endYymm* sidoCd`, 14 = `strtYymm* endYymm* sidoCd*`, 15 = `strtYymm* endYymm* dtlTmprYn sidoCd* imexTpcd* imexTmprClsfCd`, 16 = `strtYymm* endYymm* sidoCd*`, 17 = `strtYymm* endYymm* HsSgn* sidoCd*`(`HsSgn` 은 HS **6단위**, 대문자 H — 소문자 `hsSgn` 으로도 되는지 실측해 기록). `sidoCd` → `codes/시도코드.csv`. 응답 필드: 13 = `priodTitle sidoNm expCnt expUsdAmt impCnt impUsdAmt cmtrBlncAmt`, 14 = `priodTitle hsSgn korePrlstNm expLnCnt expUsdAmt impLnCnt impUsdAmt cmtrBlncAmt`, 15 = `priodTitle tmprTpcd cdValtValNm imexLnCnt imexUsdAmt`, 16 = `priodTitle sidoSggNm expCnt expUsdAmt impCnt impUsdAmt cmtrBlncAmt` (+ body `totalCount`), 17 = `priodTitle sggNm hsSgn korePrlstNm expCnt expUsdAmt impCnt impUsdAmt cmtrBlncAmt` (+ body `totalCount`). 샘플 `13_sido.xml`~`17_sigunguitem.xml`. 첫 단락에 이 계열은 필드 이름 체계가 다르다는 것과 [관세청 README 대응표](README.md#응답-필드-계열) 링크. 함정 후보: `priodTitle` 형식(월별 행인지 기간 합계 한 행인지 — 6개월 조회 결과 행 수로 판단), `시도코드 12` 전남광주통합특별시 비고, 17 은 시군구 단위 반도체(854232) 수출을 볼 수 있어 2단계 사용처(moneyflow)에 중요하다는 한 줄.
- [ ] **Step 3: 필드 대조** — `python3 $SCRATCH/check_fields.py 13 14 15 16 17` → `MISSING`/`EXTRA` 없음.
- [ ] **Step 4: 커밋**

```bash
grep -rF "$K" docs/ && echo "KEY LEAK" || echo ok
git add docs/api/customs/시도별.md docs/api/customs/시도별품목별.md docs/api/customs/시도별성질별.md docs/api/customs/시군구별.md docs/api/customs/시군구별품목별.md
git commit -m "docs: 관세청 시도별·시군구별 수출입실적 API 문서 추가" -m "Co-Authored-By: Claude Opus 5.5 <noreply@anthropic.com>
Claude-Session: https://claude.ai/code/session_01Qw5e3X1n3xofs8Vm5BbW7W"
```

---

### Task 10: 누락 점검·README·PR

**Files:**
- Create (커밋 안 함): `$SCRATCH/check_fields.py` (`SCRATCH=/private/tmp/claude-501/-Users-frankoh-src-workspace-moneyflow/3fcaf621-532b-4e7d-a6c6-767706e37219/scratchpad`) — **Task 5 Step 3 에서 먼저 만든다**
- Modify: `README.md`

- [ ] **Step 1: 필드 대조 스크립트** — 문서의 "응답 필드" 표 첫 열 집합이 (원본 명세 item 필드 ∪ fixture item 태그)와 같은지 본다. 요청 인자 표도 원본 인자 집합과 비교한다.

```python
#!/usr/bin/env python3
import html, json, os, re, sys
REPO = "/Users/frankoh/src/workspace_moneyflow/opendata-go/docs/api/customs"
FIX = os.environ.get("FIX", "/private/tmp/claude-501/-Users-frankoh-src-workspace-moneyflow/3fcaf621-532b-4e7d-a6c6-767706e37219/scratchpad/customs-fixtures")
APIS = {1: ("15100475", "품목별국가별", "01_nitemtrade"), 2: ("15101609", "품목별", "02_itemtrade"),
        3: ("15101612", "국가별", "03_nationtrade"), 4: ("15101630", "대륙별", "04_continent"),
        5: ("15101632", "경제권별", "05_economy"), 6: ("15102109", "성질별", "06_idfytemper"),
        7: ("15100476", "성질별국가별", "07_ntemper"), 8: ("15101616", "신성질별", "08_newtemper"),
        9: ("15101607", "신성질별국가별", "09_nnewtemper"), 10: ("15101634", "종류별", "10_kind"),
        11: ("15101633", "세관별", "11_customs"), 12: ("15101636", "항구공항별", "12_port"),
        13: ("15101643", "시도별", "13_sido"), 14: ("15101641", "시도별품목별", "14_sidoitem"),
        15: ("15101639", "시도별성질별", "15_sidotemper"), 16: ("15134344", "시군구별", "16_sigungu"),
        17: ("15134343", "시군구별품목별", "17_sigunguitem")}

def source_sets(did):
    p = f"{REPO}/_source/{did}.swagger.json"
    if os.path.exists(p):
        s = json.load(open(p))
        (path, ops), = s["paths"].items()
        params = {x["name"] for x in ops.get("parameters", []) + ops["get"].get("parameters", [])}
        item = ops["get"]["responses"]["200"]["schema"]["properties"]["body"]["properties"]["items"]["properties"]["item"]
        if item.get("type") == "array":
            item = item["items"]
        return params, set(item["properties"])
    t = re.sub(r"<[^>]+>", " ", open(f"{REPO}/_source/{did}.detail.html").read())
    t = re.sub(r"\s+", " ", html.unescape(t))
    req, resp = t.split("출력결과(Response Element)", 1)
    names = lambda s: set(re.findall(r"\b([a-z][A-Za-z0-9]+) \d+ (?:필수|옵션)", s))
    return names(req), names(resp) - {"resultCode", "resultMsg"}

def doc_sets(name):
    t = open(f"{REPO}/{name}.md").read()
    sec = lambda h: re.search(rf"## {h}.*?\n(.*?)(?:\n## |\Z)", t, re.S).group(1)
    first = lambda s: {m for m in re.findall(r"^\| ([A-Za-z][A-Za-z0-9_]*) \|", s, re.M)}
    return first(sec("요청 인자")), first(sec("응답 필드"))

def fixture_tags(fx):
    p = f"{FIX}/{fx}.xml"
    if not os.path.exists(p):
        return None
    tags = set()
    for it in re.findall(r"<item>(.*?)</item>", open(p).read(), re.S):
        tags |= set(re.findall(r"<([A-Za-z][A-Za-z0-9]*)>", it))
    return tags

for n in map(int, sys.argv[1:] or APIS):
    did, name, fx = APIS[n]
    sp, sf = source_sets(did)
    dp, df = doc_sets(name)
    ft = fixture_tags(fx)
    want = sf | (ft or set())
    print(f"== {n} {name} fixture={'yes' if ft is not None else 'NO'}")
    if sp - dp: print("   MISSING param", sorted(sp - dp))
    if dp - sp: print("   EXTRA param (실측 근거를 문서에 적었는지 확인)", sorted(dp - sp))
    if want - df: print("   MISSING field", sorted(want - df))
    if df - want: print("   EXTRA field", sorted(df - want))
    if ft is not None and sf - ft: print("   note: 명세에만 있고 실응답엔 없는 필드", sorted(sf - ft))
```

Run: `python3 $SCRATCH/check_fields.py`
Expected: 17개 `== n 이름 fixture=yes` 줄, `MISSING` 줄 없음. `EXTRA` 는 실측 근거(예: 실응답에만 있는 `totalCount`)가 문서에 있으면 허용. `note:` 줄은 해당 문서 함정 절에 기록됐는지 확인.

- [ ] **Step 2: 체크리스트 실행**

```bash
cd /Users/frankoh/src/workspace_moneyflow/opendata-go
ls docs/api/customs/*.md | wc -l                       # 18 (README + 17)
grep -rF "$K" docs/ scripts/ && echo LEAK || echo ok   # ok
file -I docs/api/*.md docs/api/customs/*.md | grep -v 'charset=utf-8' || echo utf8-ok
python3 - <<'EOF'
import re, pathlib
bad = []
for p in pathlib.Path("docs/api").rglob("*.md"):
    for link in re.findall(r"\]\(([^)#]+)(?:#[^)]*)?\)", p.read_text()):
        if link.startswith("http"): continue
        if not (p.parent / link).exists(): bad.append((str(p), link))
print(bad or "links-ok")
EOF
```
Expected: `18`, `ok`, `utf8-ok`, `links-ok`.

- [ ] **Step 3: 저장소 README** — `README.md` 를 다음으로 교체:

```markdown
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

data.go.kr 에서 발급받은 서비스키를 `DATA_GO_KR_API_KEY` 환경변수로 둔다.
API 마다 포털에서 활용신청이 필요하다.
```

- [ ] **Step 4: 커밋·푸시·PR**

```bash
git add README.md
git commit -m "docs: 저장소 README 에 지원 기관·문서 안내 추가" -m "Co-Authored-By: Claude Opus 5.5 <noreply@anthropic.com>
Claude-Session: https://claude.ai/code/session_01Qw5e3X1n3xofs8Vm5BbW7W"
git push -u origin feature/customs-api-docs
gh pr create --title "docs: 공공데이터포털 공통 규약 + 관세청 수출입실적 API 17개 명세 문서" --body "$(cat <<'EOF'
## 요약
- 공공데이터포털 공통 규약 문서(서비스키·게이트웨이·응답 봉투·에러)
- 관세청 수출입실적 API 17개 명세 문서 — 포털 명세 원본·실 호출 응답과 대조
- 포털 명세 원본(Swagger/상세 HTML)·코드표(xlsx→CSV) 수집 스크립트와 산출물

## 검증
- 필드 대조 스크립트: 17개 API 문서 표 = 명세 ∪ 실응답 (MISSING 없음)
- 서비스키 미포함, UTF-8, 내부 링크 확인

🤖 Generated with [Claude Code](https://claude.com/claude-code)

https://claude.ai/code/session_01Qw5e3X1n3xofs8Vm5BbW7W
EOF
)"
```
Expected: PR URL 출력. **머지는 하지 않는다** — 사용자 검수 후 머지.
