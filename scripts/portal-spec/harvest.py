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


def must(m, data_id, what):
    """정규식 매치가 없으면 데이터 ID 와 찾지 못한 대상을 밝히고 종료한다(포털 페이지 구조 변경 대비)."""
    if m is None:
        raise SystemExit(f"{data_id}: {what} 을(를) 페이지에서 찾지 못함 — 포털 페이지 구조가 바뀌었는지 확인")
    return m


def harvest_spec(data_id, src):
    page = fetch(f"{PORTAL}/data/{data_id}/openapi.do").decode("utf-8")
    title = html.unescape(must(re.search(r"<title>(.*?)</title>", page, re.S), data_id, "<title>").group(1).split("|")[0].strip())
    sj = re.search(r"const swaggerJson = `(.*?)`;", page, re.S)
    if sj and sj.group(1).strip():
        spec = json.loads(sj.group(1))
        (src / f"{data_id}.swagger.json").write_text(json.dumps(spec, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")
        return title, "swagger"
    pk = must(re.search(r'id="publicDataDetailPk" value="([^"]+)"', page), data_id, "publicDataDetailPk").group(1)
    sel = must(re.search(r'id="open_api_detail_select".*?</select>', page, re.S),
               data_id, "open_api_detail_select").group(0)
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
    file_id, sn = must(find_code_file("15101609"), "15101609", "코드표 첨부파일(fn_fileDownload)")
    xlsx = fetch(f"{PORTAL}/cmm/cmm/fileDownload.do?atchFileId={file_id}&fileDetailSn={sn}")
    (src / "관세청조회코드.xlsx").write_bytes(xlsx)
    for name, n in xlsx_to_csv(xlsx, codes):
        print(f"codes\t{name}\t{n} rows")


if __name__ == "__main__":
    main()
