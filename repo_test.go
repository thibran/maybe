package main

import (
	"errors"
	"fmt"
	"sort"
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
		f1.Path: f1,
		f2.Path: f2,
		f3.Path: f3,
	}}
}

// // All returns all folder entries of the repo.
// func (r *RepoDummy) All() []Folder {
// 	a := make([]Folder, len(r.m))
// 	var i = 0
// 	for _, v := range r.m {
// 		a[i] = v
// 		i++
// 	}
// 	return a
// }

// Add path to repo. If the path is known, the repo data is updated, else
// a new entry will be created.
func (r *RepoDummy) Add(path string, t time.Time) {
	f, ok := r.m[path]
	// new folder object
	if !ok {
		r.m[path] = Folder{
			Path:  path,
			Count: 1,
			Times: Times{t},
		}
		return
	}
	// update existing folder object
	f.Count++
	f.Times = append(f.Times, t)
	f.Times = f.Times.sort()
	r.m[path] = f
}

// Search repo for the key s.
func (r *RepoDummy) Search(s string) (RatedFolder, error) {
	a := search(r.m, s, func(a RatedFolders) { sort.Sort(a) })
	if len(a) == 0 {
		return RatedFolder{}, errors.New("no result")
	}
	return a[0], nil
}

// Show returns n RatedFolders.
func (r *RepoDummy) Show(s string, n int) RatedFolders {
	a := search(r.m, s, func(a RatedFolders) { sort.Sort(a) })
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

var f1 = Folder{
	Path:  "/home/foo",
	Count: 1,
	Times: Times{
		time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
	}}

var f2 = Folder{
	Path:  "/home/tux",
	Count: 1,
	Times: Times{
		time.Date(2012, time.February, 3, 11, 30, 0, 0, time.UTC),
	}}

var f3 = Folder{
	Path:  "/etc/apt",
	Count: 1,
	Times: Times{
		time.Date(2016, time.March, 20, 18, 0, 0, 0, time.UTC),
	}}

func nowTimeRepo() *RepoDummy {
	now := time.Now()
	v1 := Folder{
		Path:  "/home/nfoo",
		Count: 1,
		Times: Times{now.Add(-time.Second * 40)},
	}
	v2 := Folder{
		Path:  "/home/foo",
		Count: 1,
		Times: Times{now.Add(-time.Hour * 18)},
	}
	v3 := Folder{
		Path:  "/etc/apt",
		Count: 1,
		Times: Times{now.Add(-time.Hour * 24 * 7 * 2)},
	}
	return &RepoDummy{m: map[string]Folder{
		v1.Path: v1,
		v2.Path: v2,
		v3.Path: v3,
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
	a := search(r.m, "foo", func(a RatedFolders) { sort.Sort(a) })
	if a[0].Points < a[1].Points {
		t.Fail()
	}
}

func TestSearch_Repo(t *testing.T) {
	r := nowTimeRepo()
	rf, err := r.Search("foo")
	if err != nil {
		t.Fail()
	}
	if rf.Folder.Path != "/home/foo" {
		t.Fail()
	}
}

func TestShow(t *testing.T) {
	a := nowTimeRepo().Show("foo", 2)
	if len(a) == 0 {
		t.Fail()
	}
	fmt.Println("Points\tFolder")
	for _, rf := range a {
		fmt.Printf("%d\t%s\n", rf.Points, rf.Folder.Path)
	}
}
