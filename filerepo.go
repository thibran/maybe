package main

import (
	"bytes"
	"compress/gzip"
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const osSep = string(os.PathSeparator)

var (
	ignoreSlice          = []string{".git", ".hg", ".svn", ".bzr"}
	errWalkMaxDirEntries = fmt.Errorf("max entries reached")
	errNoFile            = fmt.Errorf("file not found")
)

// Repo content is saved to the disk.
type Repo struct {
	m          FolderMap
	dataPath   string
	maxEntries int
}

// NewRepo foo
func NewRepo(path string, maxEntries int) *Repo {
	return &Repo{
		m:          make(FolderMap),
		dataPath:   path,
		maxEntries: maxEntries,
	}
}

// Walk adds directories from root, for count
// of Repo.maxEntries, osWalker.lvlDeep.
func (r *Repo) Walk(root string) {
	w := newOSWalker(r, root)
	err := filepath.Walk(root, w.toWalkFunc())
	if err != nil && err != errWalkMaxDirEntries {
		log.Fatalln(err)
	}
	if tmp := os.TempDir(); tmp != "" {
		r.updateOrAddPath(tmp, w.now, true)
	}
}

// Add path to repo. If the path is known, the repo data is updated, else
// a new entry will be created.
func (r *Repo) Add(path string, t time.Time) {
	segments := strings.Split(path, osSep)
	len := len(segments)
Loop:
	for i := 0; i < len-1; i++ {
		path = strings.Join(segments[:len-i], osSep)
		if isInIgnoreList(segments[len-1-i]) {
			logf("ignore: %s\n", path)
			continue Loop
		}
		r.updateOrAddPath(path, t, i > 0)
	}
}

// updateOrAddPath to repository. Sub-folders are added, when unknown, but
// their timestamps are not updated.
func (r *Repo) updateOrAddPath(path string, t time.Time, subfolder bool) {
	f, ok := r.m[path]

	// new folder object
	if !ok {
		var sf string
		if subfolder {
			sf = "sub-"
		}
		logf("new %sfolder: %s\n", sf, path)
		r.m[path] = NewFolder(path, t)

		// guarantee folder limit holds
		if len(r.m) > r.maxEntries {
			r.m = RemoveOldestFolders(r.m, r.maxEntries-r.maxEntries/3)
		}
		return
	}

	// update existing folder object
	if subfolder {
		return
	}
	logf("update timestamps: %q\n", path)
	if f.UpdateCount < math.MaxUint32 {
		f.UpdateCount++
	}
	f.Times = append(f.Times, t)
	f.Times = f.Times.sortAndCut() // keep only data.MaxTimesEntries
	r.m[path] = f
}

var errNoResult = errors.New("no result")

// ResourceChecker returns true when a resource exists.
type ResourceChecker interface {
	doesExist(string) bool
}

// ResourceCheckerFn converter to ResourceChecker interface.
type ResourceCheckerFn func(path string) bool

func (fn ResourceCheckerFn) doesExist(path string) bool {
	return fn(path)
}

// Search repo for query.
func (r *Repo) Search(ch ResourceChecker, q query) (RatedFolder, error) {
	a := search(r.m, q.last, func(a RatedFolders) { a.sort() })
	for _, v := range filterInPathOf(a, q.start) {
		// keep not found folders, they might re-exist in future
		if ch.doesExist(v.Path) {
			// if checkFolder(v.folder.Path) {
			return v, nil
		}
	}
	return RatedFolder{}, errNoResult
}

// filterInPathOf returns a slice of entries where the path
// contains the start-string in the non-last segment.
// When start is empty, the input is returned as-is.
func filterInPathOf(a RatedFolders, start string) RatedFolders {
	start = strings.TrimSpace(strings.ToLower(start))
	if start == "" {
		return a
	}
	var res RatedFolders
	for _, f := range a {
		// ignore the last path-segment
		// path /bar/src/foo becomes /bar/src/
		pathStart, _ := filepath.Split(f.Path)
		pathStart = strings.ToLower(pathStart)
		if strings.Contains(pathStart, start) {
			res = append(res, f)
		}
	}
	return res
}

type sorterFn func(a RatedFolders)

// search for s and sort results.
func search(m FolderMap, query string, sort sorterFn) RatedFolders {
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
			var rf RatedFolder
			for folder := range tasks {
				rf = NewRatedFolder(folder, query)
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

func folderChecker() ResourceCheckerFn {
	// func returns true if the path points to a exsisting folder.
	return ResourceCheckerFn(func(path string) bool {
		if strings.TrimSpace(path) == "" {
			return false
		}
		fi, err := os.Stat(path)
		return !(os.IsNotExist(err) || !fi.IsDir())
	})
}

// List returns n RatedFolders.
func (r *Repo) List(q query, limit int, cutLong bool) RatedFolders {
	a := search(r.m, q.last, func(a RatedFolders) { a.sort() })
	a = filterInPathOf(a, q.start)
	a = cutLongPaths(a, cutLong)
	if len(a) < limit {
		limit = len(a)
	}
	return a[:limit]
}

func cutLongPaths(a RatedFolders, cutLong bool) RatedFolders {
	if !cutLong {
		return a
	}
	var res RatedFolders
	// use terminal width, when possible
	maxLineLen := 64
	if w, err := termWidth(); err == nil {
		if verbose && w-16 > 0 {
			maxLineLen = w - 16
		} else if !verbose && w-10 > 0 {
			maxLineLen = w - 10
		}
	}
	for _, rf := range a {
		rf.Path = shortenPath(rf.Path, maxLineLen)
		res = append(res, rf)
	}
	return res
}

func createTasks(m FolderMap) <-chan Folder {
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
func (r *Repo) Save() error {
	f, err := os.Create(r.dataPath)
	if err != nil {
		log.Fatalf("could not save filerepo: %s %v\n", r.dataPath, err)
	}
	defer f.Close()
	return saveGzip(f, r.m)
}

func saveGzip(w io.Writer, data FolderMap) error {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	enc.Encode(data)
	wg := gzip.NewWriter(w)
	defer wg.Close()
	wg.Write(b.Bytes())
	return nil
}

// Load repo map from dataPath.
func (r *Repo) Load() error {
	f, err := os.Open(r.dataPath)
	if err != nil {
		if os.IsNotExist(err) {
			return errNoFile
		}
		return err
	}
	defer f.Close()

	m, err := loadGzip(f)
	if err != nil {
		return err
	}
	r.m = m
	return nil
}

func loadGzip(r io.Reader) (FolderMap, error) {
	gr, err := gzip.NewReader(r)
	defer gr.Close()
	if err != nil {
		return nil, err
	}
	var m FolderMap
	dec := gob.NewDecoder(gr)
	if err := dec.Decode(&m); err != nil {
		return nil, fmt.Errorf("could not decode: %v", err)
	}
	return m, nil
}

// Size of the repository.
func (r *Repo) Size() int { return len(r.m) }
