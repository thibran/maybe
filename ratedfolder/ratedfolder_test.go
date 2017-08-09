package ratedfolder

import (
	"testing"
	"thibaut/maybe/classify"
	"thibaut/maybe/ratedfolder/folder"
	"time"
)

func TestRatedFoldersSort(t *testing.T) {
	now := time.Now()
	ratedFn := func(path string, count uint32) *RatedFolder {
		return &RatedFolder{
			Rating: &classify.Rating{TimePoints: 0},
			Folder: &folder.Folder{Path: path, UpdateCount: count,
				Times: []time.Time{now}},
		}
	}
	tt := []struct {
		name, exp string
		RatedFolders
	}{
		{name: "by path len", exp: "/home/tux/Documents",
			RatedFolders: RatedFolders{
				ratedFn("/home/tux/go/src/github.com/nsf/gocode/docs", 1),
				ratedFn("/home/tux/src/nim/Nim/tools/dochack", 1),
				ratedFn("/home/tux/src/nim/Nim/lib/packages/docutils", 1),
				ratedFn("/home/tux/Downloads", 1),
				ratedFn("/home/tux/Documents", 1),
			}},
		{name: "by count", exp: "/home/tux/src/nim/Nim/tools/dochack",
			RatedFolders: RatedFolders{
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
			if tc.RatedFolders[0].Path != tc.exp {
				t.Errorf("exp %q, got %q", tc.exp, tc.RatedFolders[0].Path)
			}
		})
	}
}

func TestRemoveOldestFolders(t *testing.T) {
	fn := func(p string, t time.Time) *folder.Folder {
		return folder.New(p, t)
	}
	now := time.Now()
	f1 := fn("/home/bar", now.Add(-time.Hour*18))
	f2 := fn("/home/zot", now.Add(-time.Hour*4))
	f3 := fn("/home/foo", now)
	tt := []struct {
		name       string
		keepValues int
		resultLen  int
		notInMap   string
	}{
		{name: "remove oldest", keepValues: 2,
			resultLen: 2, notInMap: "/home/bar"},

		{name: "no change", keepValues: 10,
			resultLen: 3, notInMap: "/aaa"},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			m := Map{f1.Path: f1, f2.Path: f2, f3.Path: f3}
			m.RemoveOldest(tc.keepValues)
			if len(m) != tc.resultLen {
				t.Fatalf("expected len(res) %d, got %v", tc.keepValues, len(m))
			}
			if _, ok := m[tc.notInMap]; ok {
				t.Fatalf("%s should not be in the map", tc.notInMap)
			}
		})
	}
}

func TestFilterInPathOf(t *testing.T) {
	now := time.Now()
	newFolder := func(p string) *folder.Folder {
		return folder.New(p, now)
	}
	a := RatedFolders{
		{Folder: newFolder("/bar/src/foo")},
		{Folder: newFolder("/bar/no/foo")},
		{Folder: newFolder("/baz/zot/foo")},
	}
	tt := []struct {
		name, start, exp string
		len              int
	}{
		{name: "ok 1", start: "src", exp: "/bar/src/foo", len: 1},
		{name: "ok 2", start: "baz", exp: "/baz/zot/foo", len: 1},
		{name: "empty", start: " ", len: 3},
		{name: "not in path", start: "cat", len: 0},
		{name: "search for last segment", start: "foo", len: 0},
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
