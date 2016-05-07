package main

import (
	"errors"
	"testing"
	"time"
)

// RepoDummy is a in-memory repo.
type RepoDummy struct {
	m map[string]Folder
}

// NewRepoDummy creates an in-memory repo containing dummy values.
func NewRepoDummy() *RepoDummy {
	return &RepoDummy{m: map[string]Folder{
		f1.path: f1,
		f2.path: f2,
		f3.path: f3,
	}}
}

// Add path to repo. If the path is known, the repo data is updated, else
// a new entry will be created.
func (r *RepoDummy) Add(path string, t time.Time) {
	f, ok := r.m[path]
	// new folder object
	if !ok {
		r.m[path] = NewFolder(
			path,
			1,
			Times{t},
		)
		return
	}
	// update existing folder object
	f.count++
	f.times = append(f.times, t)
	f.times = f.times.sort()
	r.m[path] = f
}

// Search repo for the key s.
func (r *RepoDummy) Search(s string) (RatedFolder, error) {
	a := search(r.m, s, func(a RatedFolders) { a.sort() })
	if len(a) == 0 {
		return RatedFolder{}, errors.New("no result")
	}
	return a[0], nil
}

// Show returns n RatedFolders.
func (r *RepoDummy) Show(s string, n int) RatedFolders {
	a := search(r.m, s, func(a RatedFolders) { a.sort() })
	if len(a) < n {
		n = len(a)
	}
	return a[0:n]
}

// Save method is ignored.
func (r *RepoDummy) Save() error {
	// TODO
	return nil
}

// Load method is ignored.
func (r *RepoDummy) Load() {
	// TODO
}

func (r *RepoDummy) Size() int { return len(r.m) }

var f1 = NewFolder(
	"/home/foo",
	1,
	Times{
		time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
	})

var f2 = NewFolder(
	"/home/tux",
	1,
	Times{
		time.Date(2012, time.February, 3, 11, 30, 0, 0, time.UTC),
	})

var f3 = NewFolder(
	"/etc/apt",
	1,
	Times{
		time.Date(2016, time.March, 20, 18, 0, 0, 0, time.UTC),
	})

func nowTimeRepo() *RepoDummy {
	now := time.Now()
	v1 := NewFolder(
		"/home/nfoo",
		1,
		Times{now.Add(-time.Second * 40)},
	)
	v2 := NewFolder(
		"/home/foo",
		1,
		Times{now.Add(-time.Hour * 18)},
	)
	v3 := NewFolder(
		"/etc/apt",
		1,
		Times{now.Add(-time.Hour * 24 * 7 * 2)},
	)
	return &RepoDummy{m: map[string]Folder{
		v1.path: v1,
		v2.path: v2,
		v3.path: v3,
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
	r.Add(f1.path, timeNow)
	f := r.m[f1.path]
	if f.count != 2 {
		t.Fail()
	}
	if f.times[0] != timeNow {
		t.Error("Times[0] should be equals timeNow.")
	}
	if len(f.times) > MaxTimesEntries {
		t.Fail()
	}
}

func TestSearch(t *testing.T) {
	r := nowTimeRepo()
	a := search(r.m, "foo", func(a RatedFolders) { a.sort() })
	if len(a) == 0 {
		t.Fail()
	}
	if a[0].points() < a[1].points() {
		t.Fail()
	}
}

func TestSearch_Repo(t *testing.T) {
	r := nowTimeRepo()
	rf, err := r.Search("foo")
	if err != nil {
		t.Fail()
	}
	if rf.folder.path != "/home/foo" {
		t.Fail()
	}
}

func TestShow(t *testing.T) {
	a := nowTimeRepo().Show("foo", 2)
	if len(a) == 0 {
		t.Fail()
	}
	if a[0].folder.path != "/home/foo" {
		t.Fail()
	}
}
