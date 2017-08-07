package main

import (
	"fmt"
	"path"
	"sort"
	"time"
)

// FolderMap type alias
type FolderMap map[string]*Folder

// Folder entry.
type Folder struct {
	Path        string
	UpdateCount uint32      // counts how often the folder has been updated
	Times       []time.Time // last MaxTimesEntries updates
}

// MaxTimeEntries of time.Time entries in a Times slice.
const MaxTimeEntries = 6

// NewFolder object.
func NewFolder(path string, times ...time.Time) *Folder {
	if path == "" {
		panic("NewFolder - empty path is prohibited")
	}
	if len(times) == 0 {
		panic("NewFolder - must have at last one []time entry")
	}
	return &Folder{
		Path:        path,
		UpdateCount: 1,
		Times:       times,
	}
}

// sortAndCut time entries and keep only MaxTimesEntries.
func sortAndCut(a ...time.Time) []time.Time {
	sort.Slice(a, func(i, j int) bool { return a[i].After(a[j]) })
	if len(a) > MaxTimeEntries {
		return a[:MaxTimeEntries]
	}
	return a
}

// RatedFolder object.
type RatedFolder struct {
	*Folder
	*Rating
}

// RatedFolders is an alias for RatedFolder-slice.
type RatedFolders []*RatedFolder

// NewRatedFolder creates a new object.
func NewRatedFolder(f *Folder, query string) (*RatedFolder, error) {
	if f == nil {
		return nil, fmt.Errorf("NewRatedFolder - *Folder is nil")
	}
	r, err := NewRating(query, path.Base(f.Path), f.Times...)
	if err != nil {
		return nil, fmt.Errorf("NewRatedFolder - %v", err)
	}
	return &RatedFolder{
		Folder: f,
		Rating: r,
	}, nil
}

func (a RatedFolders) sort() {
	var pi, pj uint
	sort.Slice(a, func(i, j int) bool {
		pi = a[i].points()
		pj = a[j].points()
		if pi == pj {
			if a[i].UpdateCount == a[j].UpdateCount {
				return a[i].Path < a[j].Path
			}
			return a[i].UpdateCount > a[j].UpdateCount
		}
		return pi > pj
	})
}

// RatedTimeFolders alias with time focused sort implementation.
type RatedTimeFolders []*RatedFolder

// RemoveOldest folders from map m and keep newest n entries.
func (fm *FolderMap) RemoveOldest(n int) {
	// to time-folders
	var a RatedTimeFolders
	for _, f := range *fm {
		if rf, err := NewRatedFolder(f, ""); err == nil {
			a = append(a, rf)
		}
	}
	// delete too old entries
	a.sort()
	if len(a) > n {
		a = a[:n]
	}
	m := make(FolderMap, len(a))
	for _, rf := range a {
		m[rf.Path] = rf.Folder
	}
	*fm = m
}

func (a RatedTimeFolders) sort() {
	sort.Slice(a, func(i, j int) bool {
		if a[i].timePoints == a[j].timePoints {
			return a[i].UpdateCount > a[j].UpdateCount
		}
		return a[i].timePoints > a[j].timePoints
	})
}
