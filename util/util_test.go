package util

import (
	"testing"
	"unicode/utf8"
)

func TestShortenPath(t *testing.T) {
	tt := []struct {
		name, path, exp string
		maxlen          int
	}{
		{name: "okay", path: "/1/2/3/4/5/6/7/8/9",
			exp: "/1/2/3/4/5/6/7/8/9", maxlen: 18},
		{name: "cut most", maxlen: 20,
			path: "/home/tux/src/other/lo-design/dark_with_sidebar"},
		{name: "cut mid", maxlen: 40,
			path: "/home/tux/src/other/lo-design/dark_with_sidebar"},
		{name: "cut a bit", maxlen: 50,
			path: "/home/tux/src/other/lo-design/dark_with_sidebar"},
	}
	rlen := utf8.RuneCountInString
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			res := ShortenPath(tc.path, tc.maxlen)
			if rlen(res) > tc.maxlen {
				t.Fatalf("should be not longer than %d, but is %d",
					tc.maxlen, rlen(res))
			}
			if tc.exp != "" && tc.exp != res {
				t.Fatalf("exp %q, got %q", tc.exp, res)
			}
		})
	}
}
