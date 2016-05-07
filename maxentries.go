package main

import (
	"sort"
)

// RatedTimeFolder object.
type RatedTimeFolder struct {
	Points int
	Folder Folder
}

// RatedTimeFolders is an alias for RatedTimeFolder-slice.
type RatedTimeFolders []RatedTimeFolder

// NewRatedTimeFolder creates a new object.
func NewRatedTimeFolder(f Folder) RatedTimeFolder {
	return RatedTimeFolder{
		Points: ratePassedTime(f.Times),
		Folder: f,
	}
}

// NewRatedTimeFolderSlice creator.
func NewRatedTimeFolderSlice(m map[string]Folder) RatedTimeFolders {
	a := make(RatedTimeFolders, len(m))
	for _, f := range m {
		a = append(a, NewRatedTimeFolder(f))
	}
	return a
}

// removeOldestFolders and keep n entries.
func (a RatedTimeFolders) removeOldestFolders(n int) map[string]Folder {
	sort.Sort(a)
	if len(a) > n {
		a = a[:n]
	}
	m := make(map[string]Folder, len(a))
	for _, rf := range a {
		m[rf.Folder.Path] = rf.Folder
	}
	return m
}

func (a RatedTimeFolders) Len() int {
	return len(a)
}

func (a RatedTimeFolders) Less(i, j int) bool {
	if a[i].Points == a[j].Points {
		return a[i].Folder.Count > a[j].Folder.Count
	}
	return a[i].Points > a[j].Points
}

func (a RatedTimeFolders) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
