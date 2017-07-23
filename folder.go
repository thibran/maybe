package main

import (
	"fmt"
	"log"
	"path"
	"sort"
	"strings"
	"time"
)

// Folder entry.
type Folder struct {
	Path  string
	Count uint  // counts how often the folder has been updated
	Times Times // last MaxTimesEntries updates
}

// NewFolder object.
func NewFolder(path string, count uint, times ...time.Time) Folder {
	if len(times) == 0 {
		log.Fatal("NewFolder - must have at last one []time entry")
	}
	return Folder{
		Path:  path,
		Count: count,
		Times: times,
	}
}

// Times is a shorthand for a time slice.
type Times []time.Time

// MaxTimeEntries of time.Time entries in a Times slice.
const MaxTimeEntries = 6

// sort time entries and cut all entries longer than MaxTimesEntries.
func (t Times) sort() Times {

	sort.Slice(t, func(i, j int) bool { return t[i].After(t[j]) })
	if len(t) > MaxTimeEntries {
		return t[:MaxTimeEntries]
	}
	return t
}

// RatedFolder object.
type RatedFolder struct {
	timePoints       uint
	similarityPoints uint
	folder           Folder
}

// points return the points of a rated folder.
// If no similarity is found, time points are ignored.
func (rf *RatedFolder) points() uint {
	n := rf.similarityPoints
	if n == noMatch {
		return n
	}
	return n + rf.timePoints
}

// RatedFolders is an alias for RatedFolder-slice.
type RatedFolders []RatedFolder

// NewRatedFolder creates a new object.
func NewRatedFolder(f Folder, query string) RatedFolder {
	return RatedFolder{
		timePoints:       ratePassedTime(f.Times),
		similarityPoints: rateSimilarity(path.Base(f.Path), query),
		folder:           f,
	}
}

func (a RatedFolders) String() string {
	arr := make([]string, len(a))
	for k, v := range a {
		arr[k] = fmt.Sprintf("%v  tp: %v  sp: %v",
			v.folder, v.timePoints, v.similarityPoints)
	}
	return strings.Join(arr, "\n")
}

func (a RatedFolders) sort() {
	sort.Slice(a, func(i, j int) bool {
		pi := a[i].points()
		pj := a[j].points()
		if pi == pj {
			return a[i].folder.Count > a[j].folder.Count
		}
		return pi > pj
	})
}

// RatedTimeFolders alias with time focused sort implementation.
type RatedTimeFolders []RatedFolder

// RemoveOldestFolders from map m, keep newest n entries.
func RemoveOldestFolders(m map[string]Folder, n int) map[string]Folder {
	return fromFolderMap(m).removeOldestFolders(n)
}

func fromFolderMap(m map[string]Folder) RatedTimeFolders {
	a := make(RatedTimeFolders, len(m))
	for _, f := range m {
		a = append(a, NewRatedFolder(f, ""))
	}
	return a
}

// removeOldestFolders and keep n entries.
func (a RatedTimeFolders) removeOldestFolders(n int) map[string]Folder {
	a.sort()
	if len(a) > n {
		a = a[:n]
	}
	m := make(map[string]Folder, len(a))
	for _, rf := range a {
		m[rf.folder.Path] = rf.folder
	}
	return m
}

func (a RatedTimeFolders) sort() {
	sort.Slice(a, func(i, j int) bool {
		if a[i].timePoints == a[j].timePoints {
			return a[i].folder.Count > a[j].folder.Count
		}
		return a[i].timePoints > a[j].timePoints
	})
}
