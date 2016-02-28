package main

import (
	"errors"
	"sort"
	"time"
)

// Repo abstracts the data storage.
type Repo interface {
	Add(path string, t time.Time)         // Add new folder to the repo.
	Search(s string) (RatedFolder, error) // Search for the key s in the repo
}

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
	var a RatedFolders
	for _, f := range r.m {
		rf := NewRatedFolder(f, s)
		if rf.Points == NoMatch {
			continue
		}
		a = append(a, rf)
	}
	if len(a) == 0 {
		return RatedFolder{}, errors.New("no result")
	}
	sort.Sort(a)
	return a[0], nil
}
