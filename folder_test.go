package main

import (
	"testing"
	"time"
)

func TestTimesSort(t *testing.T) {
	var i int
	now := time.Now()
	timeFn := func() time.Time {
		i++
		return now.Add(time.Hour + time.Duration(i))
	}
	var a []time.Time
	for len(a) <= MaxTimeEntries {
		a = append(a, timeFn())
	}
	a = sortAndCut(a...)

	if len(a) != MaxTimeEntries {
		t.Errorf("len(a) should be %d, got %d", MaxTimeEntries, len(a))
	}
	exp := now.Add(time.Hour + time.Duration(1))
	if a[0].Hour() != exp.Hour() {
		t.Fatalf("exp %q, got %q", exp, a[0])
	}
	exp = now.Add(time.Hour + time.Duration(6))
	if a[5].Hour() != exp.Hour() {
		t.Fatalf("exp %q, got %q", exp, a[5])
	}
}

func TestNewFolder(t *testing.T) {
	now := []time.Time{time.Now()}
	tt := []struct {
		name, path string
		panic      bool
		times      []time.Time
	}{
		{name: "okay", panic: false, path: "/foo", times: now},
		{name: "no time value", panic: true, path: "/foo", times: nil},
		{name: "empty path", panic: true, path: "", times: now},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				if err := recover(); err == nil && tc.panic {
					t.Errorf("%q should not panic", tc.name)
				}
			}()
			NewFolder(tc.path, tc.times...)
		})
	}
}

func TestTimeRatedSort(t *testing.T) {
	now := time.Now()
	rated := func(path string, timePoints uint, count uint32) *RatedFolder {
		return &RatedFolder{
			Rating: &Rating{timePoints: timePoints},
			Folder: &Folder{Path: path, UpdateCount: count,
				Times: []time.Time{now}},
		}
	}
	tt := []struct {
		name, exp string
		folders   RatedTimeFolders
	}{
		{name: "by time points", exp: "/home/foo", folders: RatedTimeFolders{
			rated("/home/bar", 4, 1),
			rated("/home/foo", 10, 1),
			rated("/home/zot", 8, 1),
		}},
		{name: "by time count", exp: "/b", folders: RatedTimeFolders{
			rated("/a", 20, 1),
			rated("/b", 20, 2),
		}},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.folders.sort(); tc.folders[0].Path != tc.exp {
				t.Fatalf("exp %q, got %q", tc.exp, tc.folders[0].Path)
			}
		})
	}
}

func TestRemoveOldestFolders(t *testing.T) {
	fn := func(p string, t time.Time) *Folder {
		return NewFolder(p, t)
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
			m := FolderMap{f1.Path: f1, f2.Path: f2, f3.Path: f3}
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

func TestRatedFoldersSort(t *testing.T) {
	now := time.Now()
	ratedFn := func(path string, count uint32) *RatedFolder {
		return &RatedFolder{
			Rating: &Rating{timePoints: 0},
			Folder: &Folder{Path: path, UpdateCount: count,
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
			tc.sort()
			if tc.RatedFolders[0].Path != tc.exp {
				t.Errorf("exp %q, got %q", tc.exp, tc.RatedFolders[0].Path)
			}
		})
	}
}
