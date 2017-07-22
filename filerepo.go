package main

import (
	"bytes"
	"compress/gzip"
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

// Saver abstracts saving of implementing object.
type Saver interface {
	Save() error
}

// Loader abstract loading of implementing object.
type Loader interface {
	Load() error
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

var ignoreSlice = []string{".git", ".hg", ".svn", ".bzr"}

const osSep = string(os.PathSeparator)

// Add path to repo. If the path is known, the repo data is updated, else
// a new entry will be created.
func (r *FileRepo) Add(path string, t time.Time) {
	segments := strings.Split(path, osSep)
	len := len(segments)

Loop:
	for i := 0; i < len-1; i++ {
		path = strings.Join(segments[:len-i], osSep)
		// check if folder is in the ignore list
		for _, ign := range ignoreSlice {
			if segments[len-1-i] == ign {
				logf("Add - ignore: %s\n", path)
				continue Loop
			}
		}
		r.updateOrAddPath(path, t, i > 0)
	}
}

// updateOrAddPath to repository. Sub-folders are added, when unknown, but
// their timestamps are not updated.
func (r *FileRepo) updateOrAddPath(path string, t time.Time, subfolder bool) {
	f, ok := r.m[path]
	// new folder object
	if !ok {
		logln("updateOrAddPath - new folder:", path)
		r.m[path] = NewFolder(path, 1, Times{t})
		return
	}
	if subfolder {
		return
	}
	logln("updateOrAddPath - update timestamps:", path)

	// update existing folder object
	f.Count++
	f.Times = append(f.Times, t)
	f.Times = f.Times.sort() // sort and keep only data.MaxTimesEntries
	r.m[path] = f
	// remove oldest File entries if necessary
	if len(r.m) <= r.maxEntries {
		return
	}
	r.m = RemoveOldestFolders(r.m, r.maxEntries-r.maxEntries/3)
}

var errNoResult = errors.New("no result")

// Search repo for the key s.
func (r *FileRepo) Search(query string) (RatedFolder, error) {
	a := search(r.m, query, func(a RatedFolders) { a.sort() })
	for _, v := range a {
		// keep not found folders, they might re-exist in future
		if checkFolder(v.folder.Path) {
			return v, nil
		}
	}
	return RatedFolder{}, errNoResult
}

// checkFolder returns true if the path points to a folder which exists.
func checkFolder(path string) bool {
	if strings.TrimSpace(path) == "" {
		return false
	}
	fi, err := os.Stat(path)
	return !(os.IsNotExist(err) || !fi.IsDir())
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
	return saveGzip(r.dataPath, r.m)
}

func saveGzip(path string, data map[string]Folder) error {
	f, err := os.Create(path)
	if err != nil {
		log.Fatalf("could not save filerepo: %s %v\n", path, err)
	}
	defer f.Close()

	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	enc.Encode(data)

	w := gzip.NewWriter(f)
	defer w.Close()
	w.Write(b.Bytes())
	return nil
}

var errNoFile = fmt.Errorf("file not found")

// Load repo map from dataPath.
func (r *FileRepo) Load() error {
	m, err := loadGzip(r.dataPath)
	if err != nil {
		return err
	}
	r.m = m
	return nil
}

func loadGzip(path string) (map[string]Folder, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, errNoFile
	}
	defer f.Close()

	r, err := gzip.NewReader(f)
	defer r.Close()
	if err != nil {
		return nil, err
	}

	var m map[string]Folder
	dec := gob.NewDecoder(r)
	if err := dec.Decode(&m); err != nil {
		return nil, fmt.Errorf("can not decode: %s %v", path, err)
	}
	return m, nil
}

// Size of the repository.
func (r *FileRepo) Size() int { return len(r.m) }
