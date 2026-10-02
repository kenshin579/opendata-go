package opendata

import (
	"errors"
	"testing"
)

func TestNumInt64(t *testing.T) {
	cases := []struct {
		in   Num
		want int64
		err  error
	}{
		{"915255303", 915255303, nil},
		{"      38,179,150", 38179150, nil},
		{"          -2,218", -2218, nil},
		{"0", 0, nil},
		{"-", 0, ErrNotNumber},
		{"", 0, ErrNotNumber},
		{"  ", 0, ErrNotNumber},
		{"1.5", 0, ErrNotNumber},
	}
	for _, c := range cases {
		got, err := c.in.Int64()
		if got != c.want || !errors.Is(err, c.err) {
			t.Errorf("Num(%q).Int64() = %d, %v; want %d, %v", c.in, got, err, c.want, c.err)
		}
	}
	if Num(" 1,2 ").String() != " 1,2 " {
		t.Error("String must return the raw value")
	}
}
