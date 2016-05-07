package main

import (
	"sort"
	"testing"
	"time"
)

func dummyRatedTimeFolders() RatedTimeFolders {
	now := time.Now()
	f1 := NewRatedTimeFolder(Folder{
		Path:  "/home/bar",
		Count: 1,
		Times: Times{now.Add(-time.Hour * 18)},
	})
	f2 := NewRatedTimeFolder(Folder{
		Path:  "/home/zot",
		Count: 1,
		Times: Times{now.Add(-time.Hour * 4)},
	})
	f3 := NewRatedTimeFolder(Folder{
		Path:  "/home/foo",
		Count: 1,
		Times: Times{now},
	})
	return RatedTimeFolders{f1, f2, f3}
}

func TestNewRatedTimeFolder(t *testing.T) {
	a := dummyRatedTimeFolders()
	sort.Sort(a)
	if a[0].Folder.Path != "/home/foo" {
		t.Fail()
	}
}

func TestRemoveOldestFolders(t *testing.T) {
	a := dummyRatedTimeFolders()
	m := a.removeOldestFolders(2)
	if len(m) != 2 {
		t.Fail()
	}
	if _, ok := m["/home/bar"]; ok {
		t.Fail()
	}
}
