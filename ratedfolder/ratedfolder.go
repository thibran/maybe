package ratedfolder

import (
	"fmt"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"thibaut/maybe/classify"
	"thibaut/maybe/pref"
	"thibaut/maybe/ratedfolder/folder"
	"thibaut/maybe/util"
)

// RatedFolder object.
type RatedFolder struct {
	*folder.Folder
	*classify.Rating
}

// RatedFolders is an alias for RatedFolder-slice.
type RatedFolders []*RatedFolder

// New creates a new rated folder object.
func New(f *folder.Folder, query string) (*RatedFolder, error) {
	if f == nil {
		return nil, fmt.Errorf("NewRatedFolder - *Folder is nil")
	}
	r, err := classify.NewRating(query, path.Base(f.Path), f.Times...)
	if err != nil {
		return nil, fmt.Errorf("NewRatedFolder - %v", err)
	}
	return &RatedFolder{
		Folder: f,
		Rating: r,
	}, nil
}

// Sort a RatedFolders.
func (a RatedFolders) Sort() {
	var pi, pj uint
	sort.Slice(a, func(i, j int) bool {
		pi = a[i].Points()
		pj = a[j].Points()
		if pi == pj {
			if a[i].UpdateCount == a[j].UpdateCount {
				return a[i].Path < a[j].Path
			}
			return a[i].UpdateCount > a[j].UpdateCount
		}
		return pi > pj
	})
}

// Map type alias
type Map map[string]*folder.Folder

// RemoveOldest folders from map m and keep newest n entries.
func (fm *Map) RemoveOldest(n int) {
	// to time-folders
	var a RatedTimeFolders
	for _, f := range *fm {
		if rf, err := New(f, ""); err == nil {
			a = append(a, rf)
		}
	}
	// delete too old entries
	a.sort()
	if len(a) > n {
		a = a[:n]
	}
	m := make(Map, len(a))
	for _, rf := range a {
		m[rf.Path] = rf.Folder
	}
	*fm = m
}

type sorterFn func(a RatedFolders)

// Search for s and sort results.
func Search(m Map, query string, sort sorterFn) RatedFolders {
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
	results := make(chan *RatedFolder)
	go func() {
		wg.Wait()
		close(results)
	}()
	for i := 0; i < workers; i++ {
		go func() {
			for f := range tasks {
				rf, err := New(f, query)
				if err != nil || rf.Points() == classify.NoMatch {
					continue
				}
				results <- rf
			}
			wg.Done()
		}()
	}
	return collectResults(results, sort)
}

func createTasks(m Map) <-chan *folder.Folder {
	tasks := make(chan *folder.Folder)
	go func() {
		for _, folder := range m {
			tasks <- folder
		}
		close(tasks)
	}()
	return tasks
}

func collectResults(c <-chan *RatedFolder, sort sorterFn) RatedFolders {
	var a RatedFolders
	for r := range c {
		a = append(a, r)
	}
	sort(a)
	return a
}

// ResourceChecker returns true when a resource exists.
type ResourceChecker interface {
	DoesExist(string) bool
}

// FilterInPathOf returns a slice of entries where the path
// contains the start-string in the non-last segment.
// When start is empty, the input is returned as-is.
func (rf *RatedFolders) FilterInPathOf(start string) {
	start = strings.TrimSpace(strings.ToLower(start))
	if start == "" {
		return
	}
	var res RatedFolders
	for _, f := range *rf {
		// ignore the last path-segment
		// path /bar/src/foo becomes /bar/src/
		pathStart, _ := filepath.Split(f.Path)
		pathStart = strings.ToLower(pathStart)
		if strings.Contains(pathStart, start) {
			res = append(res, f)
		}
	}
	*rf = res
}

// CutLongPaths if too long.
func (a *RatedFolders) CutLongPaths(cutLong bool) {
	if !cutLong {
		return
	}
	var res RatedFolders
	// use terminal width, when possible
	maxLineLen := 64
	if w, err := util.TermWidth(); err == nil {
		if pref.Verbose && w-16 > 0 {
			maxLineLen = w - 16
		} else if !pref.Verbose && w-10 > 0 {
			maxLineLen = w - 10
		}
	}
	for _, rf := range *a {
		rf.Path = util.ShortenPath(rf.Path, maxLineLen)
		res = append(res, rf)
	}
	*a = res
}
