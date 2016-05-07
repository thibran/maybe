package main

import (
	"encoding/gob"
	"errors"
	"log"
	"os"
	"sync"
	"time"
)

// Saver abstracts saving of implementing object.
type Saver interface {
	Save() error
}

// Loader abstract loading of implementing object.
type Loader interface {
	Load()
}

// Repo abstracts the data storage.
type Repo interface {
	Add(path string, t time.Time)         // Add new folder to the repo.
	Search(s string) (RatedFolder, error) // Search for the key s in the repo
	Show(s string, n int) RatedFolders    // Show returns n RatedFolders.
	Size() int                            // Size returns the number of entries.
	Saver
	Loader
}

// FileRepo content is saved to the disk.
type FileRepo struct {
	m          map[string]Folder
	dataPath   string
	maxEntries int
}

// NewFileRepo foo
func NewFileRepo(path string, maxEntries int) *FileRepo {
	return &FileRepo{
		m:          make(map[string]Folder),
		dataPath:   path,
		maxEntries: maxEntries,
	}
}

// Add path to repo. If the path is known, the repo data is updated, else
// a new entry will be created.
func (r *FileRepo) Add(path string, t time.Time) {
	f, ok := r.m[path]
	// new folder object
	if !ok {
		r.m[path] = NewFolder(path, 1, Times{t})
		return
	}
	// update existing folder object
	f.count++
	f.times = append(f.times, t)
	f.times = f.times.sort() // sort and keep only data.MaxTimesEntries
	r.m[path] = f
	// remove oldest File entries if necessary
	if len(r.m) <= r.maxEntries {
		return
	}
	r.m = RemoveOldestFolders(r.m, r.maxEntries-r.maxEntries/3)
}

// Search repo for the key s.
func (r *FileRepo) Search(query string) (RatedFolder, error) {
	a := search(r.m, query, func(a RatedFolders) { a.sort() })
	if len(a) == 0 {
		return RatedFolder{}, errors.New("no result")
	}
	return a[0], nil
}

// Show returns n RatedFolders.
func (r *FileRepo) Show(query string, limit int) RatedFolders {
	a := search(r.m, query, func(a RatedFolders) { a.sort() })
	if len(a) < limit {
		limit = len(a)
	}
	return a[0:limit]
}

type sorterFn func(a RatedFolders)

// search for s and sort results.
func search(m map[string]Folder, query string, sort sorterFn) RatedFolders {
	if len(m) == 0 {
		return RatedFolders{}
	}
	var wg sync.WaitGroup
	workers := 32
	if len(m) < workers {
		workers = len(m)
	}
	wg.Add(workers)

	tasks := createTasks(m)
	results := make(chan RatedFolder)

	go func() {
		wg.Wait()
		close(results)
	}()

	for i := 0; i < workers; i++ {
		go func() {
			for folder := range tasks {
				// for _, f := range m {
				rf := NewRatedFolder(folder, query)
				if rf.points() == noMatch {
					continue
				}
				results <- rf
			}
			wg.Done()
		}()
	}
	return collectResults(results, sort)
}

func createTasks(m map[string]Folder) <-chan Folder {
	tasks := make(chan Folder)
	go func() {
		for _, folder := range m {
			tasks <- folder
		}
		close(tasks)
	}()
	return tasks
}

func collectResults(c <-chan RatedFolder, sort sorterFn) RatedFolders {
	var a RatedFolders
	for r := range c {
		a = append(a, r)
	}
	sort(a)
	return a
}

// Save repo map to dataPath.
func (r *FileRepo) Save() error {
	f, err := os.Create(r.dataPath)
	if err != nil {
		log.Fatalf("could not save filerepo: %s %v\n", r.dataPath, err)
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

// Size of the repository.
func (r *FileRepo) Size() int { return len(r.m) }
