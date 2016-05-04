package main

import (
	"encoding/gob"
	"errors"
	"os"
	"sort"
	"time"
)

// FileRepo content is saved to the disk.
type FileRepo struct {
	m        map[string]Folder
	dataPath string
}

// NewFileRepo foo
func NewFileRepo(path string) *FileRepo {
	return &FileRepo{
		m:        make(map[string]Folder),
		dataPath: path,
	}
}

// Add path to repo. If the path is known, the repo data is updated, else
// a new entry will be created.
func (r *FileRepo) Add(path string, t time.Time) {
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
	f.Times = f.Times.sort() // sort and keep only data.MaxTimesEntries
	r.m[path] = f

	// ++++++++++++++++++
	//if 

}

// search for s and sort results.
func search(m map[string]Folder, s string, f sortRatedFolders) RatedFolders {
	var a RatedFolders
	for _, f := range m {
		rf := NewRatedFolder(f, s)
		if rf.Points == NoMatch {
			continue
		}
		a = append(a, rf)
	}
	f(a)
	return a
}

// Search repo for the key s.
func (r *FileRepo) Search(s string) (RatedFolder, error) {
	a := search(r.m, s, func(a RatedFolders) { sort.Sort(a) })
	if len(a) == 0 {
		return RatedFolder{}, errors.New("no result")
	}
	return a[0], nil
}

// Show returns n RatedFolders.
func (r *FileRepo) Show(s string, n int) RatedFolders {
	a := search(r.m, s, func(a RatedFolders) { sort.Sort(a) })
	if len(a) < n {
		n = len(a)
	}
	return a[0:n]
}

// Save repo map to dataPath.
func (r *FileRepo) Save() error {
	f, err := os.Create(r.dataPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	enc := gob.NewEncoder(f)
	enc.Encode(r.m)
	return nil
}

// Load repo map from dataPath.
func (r *FileRepo) Load() {
	f, err := os.Open(r.dataPath)
	if err != nil {
		// no file found
		return
	}
	defer f.Close()
	var m map[string]Folder
	dec := gob.NewDecoder(f)
	if err := dec.Decode(&m); err != nil {
		panic(err)
	}
	r.m = m
}

// // search for s and sort results.
// func (r *FileRepo) search(s string) RatedFolders {
// 	var a RatedFolders
// 	for _, f := range r.m {
// 		rf := NewRatedFolder(f, s)
// 		if rf.Points == NoMatch {
// 			continue
// 		}
// 		a = append(a, rf)
// 	}
// 	sort.Sort(a)
// 	return a
// }
