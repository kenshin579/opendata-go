package customs

import (
	"reflect"
	"testing"
)

func TestPeriodQuery(t *testing.T) {
	ok := [][2]string{{"202601", "202601"}, {"202507", "202606"}, {"202601", "202612"}}
	for _, p := range ok {
		q, err := periodQuery(p[0], p[1])
		if err != nil || q.Get("strtYymm") != p[0] || q.Get("endYymm") != p[1] {
			t.Errorf("%v: q=%v err=%v", p, q, err)
		}
	}
	bad := [][2]string{
		{"202506", "202606"}, // 13개월
		{"202606", "202601"}, // 역순
		{"2026-01", "202601"},
		{"202613", "202613"},
		{"202600", "202601"},
		{"", "202601"},
	}
	for _, p := range bad {
		if _, err := periodQuery(p[0], p[1]); err == nil {
			t.Errorf("%v: want error", p)
		}
	}
}

func TestYears(t *testing.T) {
	got, err := Years("202407", "202606")
	want := []Period{{"202407", "202412"}, {"202501", "202512"}, {"202601", "202606"}}
	if err != nil || !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v %v", got, err)
	}
	got, _ = Years("202603", "202605")
	if !reflect.DeepEqual(got, []Period{{"202603", "202605"}}) {
		t.Fatalf("single year: %v", got)
	}
	for _, p := range got {
		if _, err := periodQuery(p.Start, p.End); err != nil {
			t.Fatalf("window %v rejected: %v", p, err)
		}
	}
	if _, err := Years("202606", "202601"); err == nil {
		t.Fatal("want error for reversed")
	}
}

type fakeRow struct{ total bool }

func (f fakeRow) IsTotal() bool { return f.total }

func TestWithoutTotal(t *testing.T) {
	got := WithoutTotal([]fakeRow{{true}, {false}, {false}})
	if len(got) != 2 {
		t.Fatalf("got %d", len(got))
	}
}
