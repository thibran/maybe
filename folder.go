package main

import (
	"path"
	"sort"
	"time"
)

// FolderMap type alias
type FolderMap map[string]Folder

// Folder entry.
type Folder struct {
	Path        string
	UpdateCount uint32 // counts how often the folder has been updated
	Times       Times  // last MaxTimesEntries updates
}

// Times is a shorthand for a time slice.
type Times []time.Time

// MaxTimeEntries of time.Time entries in a Times slice.
const MaxTimeEntries = 6

// NewFolder object.
func NewFolder(path string, times ...time.Time) Folder {
	if len(path) == 0 {
		panic("NewFolder - empty path is prohibited")
	}
	if len(times) == 0 {
		panic("NewFolder - must have at last one []time entry")
	}
	return Folder{
		Path:        path,
		UpdateCount: 1,
		Times:       times,
	}
}

// sortAndCut time entries and keep only MaxTimesEntries.
func (t Times) sortAndCut() Times {
	sort.Slice(t, func(i, j int) bool { return t[i].After(t[j]) })
	if len(t) > MaxTimeEntries {
		return t[:MaxTimeEntries]
	}
	return t
}

// RatedFolder object.
type RatedFolder struct {
	Folder
	Rating
}

// RatedFolders is an alias for RatedFolder-slice.
type RatedFolders []RatedFolder

// NewRatedFolder creates a new object.
func NewRatedFolder(f Folder, query string) RatedFolder {
	return RatedFolder{
		Folder: f,
		Rating: newRating(query, path.Base(f.Path), f.Times),
	}
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
type RatedTimeFolders []RatedFolder

// RemoveOldestFolders from map m, keep newest n entries.
func RemoveOldestFolders(m FolderMap, n int) FolderMap {
	// to time-folders
	a := make(RatedTimeFolders, len(m))
	var i int
	for _, f := range m {
		a[i] = NewRatedFolder(f, "")
		i++
	}
	// delete too old entries
	a.sort()
	if len(a) > n {
		a = a[:n]
	}
	m = make(FolderMap, len(a))
	for _, rf := range a {
		m[rf.Path] = rf.Folder
	}
	return m
}

func (a RatedTimeFolders) sort() {
	sort.Slice(a, func(i, j int) bool {
		if a[i].timePoints == a[j].timePoints {
			return a[i].UpdateCount > a[j].UpdateCount
		}
		return a[i].timePoints > a[j].timePoints
	})
}
