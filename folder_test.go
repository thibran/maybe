package main

import (
	"testing"
	"time"
)

func TestSort(t *testing.T) {
	a := Times{f1.Times[0], f2.Times[0], f3.Times[0]}
	a = a.sort()
	if a[0] != f3.Times[0] {
		t.Fail()
	}
	if a[1] != f2.Times[0] {
		t.Fail()
	}
}

func TestNewFolder(t *testing.T) {
	now := []time.Time{time.Now()}
	tt := []struct {
		name  string
		panic bool
		path  string
		times []time.Time
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

func dummy() RatedTimeFolders {
	fn := func(p string, t time.Time) RatedFolder {
		return NewRatedFolder(NewFolder(p, t), "")
	}
	now := time.Now()
	f1 := fn("/home/bar", now.Add(-time.Hour*18))
	f2 := fn("/home/zot", now.Add(-time.Hour*4))
	f3 := fn("/home/foo", now)
	return RatedTimeFolders{f1, f2, f3}
}

func TestTimeRatedSort(t *testing.T) {
	a := dummy()
	a.sort()
	if a[0].folder.Path != "/home/foo" {
		t.Fatal()
	}
}

func TestRemoveOldestFolders(t *testing.T) {
	setup := func() map[string]Folder {
		now := time.Now()
		fn := func(p string, t time.Time) Folder {
			return NewFolder(p, t)
		}
		f1 := fn("/home/bar", now.Add(-time.Hour*18))
		f2 := fn("/home/zot", now.Add(-time.Hour*4))
		f3 := fn("/home/foo", now)
		return map[string]Folder{
			f1.Path: f1,
			f2.Path: f2,
			f3.Path: f3,
		}
	}

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
			res := RemoveOldestFolders(setup(), tc.keepValues)
			if len(res) != tc.resultLen {
				t.Fatalf("expected len(res) %d, got %v", tc.keepValues, len(res))
			}
			if _, ok := res[tc.notInMap]; ok {
				t.Fatalf("%s should not be in the map", tc.notInMap)
			}
		})
	}
}
