package main

import (
	"testing"
	"time"
)

func nowTimeRepo() *RepoDummy {
	now := time.Now()
	f1 := Folder{
		Path:  "/home/nfoo",
		Count: 1,
		Times: Times{now.Add(-time.Second * 40)},
	}
	f2 := Folder{
		Path:  "/home/foo",
		Count: 1,
		Times: Times{now.Add(-time.Hour * 18)},
	}
	f3 := Folder{
		Path:  "/etc/apt",
		Count: 1,
		Times: Times{now.Add(-time.Hour * 24 * 7 * 2)},
	}
	return &RepoDummy{m: map[string]Folder{
		f1.Path: f1,
		f2.Path: f2,
		f3.Path: f3,
	}}
}

func TestNewRepoDummy(t *testing.T) {
	r := NewRepoDummy()
	if r == nil {
		t.Fail()
	}
}

func TestAdd_newObj(t *testing.T) {
	r := NewRepoDummy()
	oldLen := len(r.m)
	r.Add("/foo/bar", time.Now())
	if len(r.m) != oldLen+1 {
		t.Fail()
	}
}

func TestAdd_updateExisting(t *testing.T) {
	r := NewRepoDummy()
	timeNow := time.Now()
	r.Add(f1.Path, timeNow)
	f := r.m[f1.Path]
	if f.Count != 2 {
		t.Fail()
	}
	if f.Times[0] != timeNow {
		t.Error("Times[0] should be equals timeNow.")
	}
	if len(f.Times) > MaxTimesEntries {
		t.Fail()
	}
}

func TestSearch(t *testing.T) {
	r := nowTimeRepo()
	rf, err := r.Search("foo")
	if err != nil {
		t.Fail()
	}
	if rf.Folder.Path != "/home/foo" {
		t.Fail()
	}
}
