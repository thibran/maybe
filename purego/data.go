package main

import (
	"sort"
	"time"
)

// Times is a shorthand for a time slice.
type Times []time.Time

func (t Times) Len() int {
	return len(t)
}

func (t Times) Less(i, j int) bool {
	return t[i].After(t[j])
}

func (t Times) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

// sort time entries and cut all entries longer than MaxTimesEntries.
func (t Times) sort() Times {
	sort.Sort(t)
	if len(t) > MaxTimesEntries {
		return t[:MaxTimesEntries]
	}
	return t
}

// Folder entry.
type Folder struct {
	Path  string
	Count int   // counts how often the folder has been updated
	Times Times // last MaxTimesEntries updates
}

// maximum number of time.Time entries in a Times slice.
const MaxTimesEntries = 10

// func NewFolder(path string, t time.Time) Folder {
// 	return Folder{
// 		Path:  path,
// 		Count: 1,
// 		Times: Times{t},
// 	}
// }

// type Event struct {
// 	Time time.Time
// }
