package rated

import (
	"testing"
	"time"

	"github.com/thibran/maybe/classify"
	"github.com/thibran/maybe/rated/folder"
)

func TestSort(t *testing.T) {
	now := time.Now()
	ratedFn := func(path string, count uint32) *Rated {
		return &Rated{
			Rating: &classify.Rating{TimePoints: 0},
			Folder: &folder.Folder{Path: path, UpdateCount: count,
				Times: []time.Time{now}},
		}
	}
	tt := []struct {
		name, exp string
		Slice
	}{
		{name: "by path len", exp: "/home/tux/Documents",
			Slice: Slice{
				ratedFn("/home/tux/go/src/github.com/nsf/gocode/docs", 1),
				ratedFn("/home/tux/src/nim/Nim/tools/dochack", 1),
				ratedFn("/home/tux/src/nim/Nim/lib/packages/docutils", 1),
				ratedFn("/home/tux/Downloads", 1),
				ratedFn("/home/tux/Documents", 1),
			}},
		{name: "by count", exp: "/home/tux/src/nim/Nim/tools/dochack",
			Slice: Slice{
				ratedFn("/home/tux/go/src/github.com/nsf/gocode/docs", 1),
				ratedFn("/home/tux/src/nim/Nim/tools/dochack", 3),
				ratedFn("/home/tux/src/nim/Nim/lib/packages/docutils", 1),
				ratedFn("/home/tux/Downloads", 2),
				ratedFn("/home/tux/Documents", 1),
			}},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tc.Sort()
			if tc.Slice[0].Path != tc.exp {
				t.Errorf("exp %q, got %q", tc.exp, tc.Slice[0].Path)
			}
		})
	}
}

func TestFilterInPathOf(t *testing.T) {
	now := time.Now()
	newFolder := func(p string) *folder.Folder {
		return folder.New(p, now)
	}
	a := Slice{
		{Folder: newFolder("/bar/src/foo")},
		{Folder: newFolder("/bar/no/foo")},
		{Folder: newFolder("/baz/zot/foo")},
		{Folder: newFolder("/joe/.cargo/bin")},
		{Folder: newFolder("/joe/hugo/go/bin")},
	}
	tt := []struct {
		name, start, exp string
		len              int
	}{
		{name: "ok 1", start: "src", exp: "/bar/src/foo", len: 1},
		{name: "ok 2", start: "baz", exp: "/baz/zot/foo", len: 1},
		{name: "empty", start: " ", len: 5},
		{name: "not in path", start: "cat", len: 0},
		{name: "search for last segment", start: "foo", len: 0},
		{name: "ignore suffix", start: "go", exp: "/joe/hugo/go/bin", len: 1},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			arr := a
			arr.FilterInPathOf(tc.start)
			if len(arr) != tc.len {
				t.Fail()
			}
			if tc.exp != "" && arr[0].Path != tc.exp {
				t.Errorf("%s - exp %v, got %v", tc.name, tc.exp, arr[0].Path)
			}
		})
	}
}
