package main

import (
    "strings"
    "fmt"
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
		//Points: rate("", f.Path, f.Times), // todo rate only by time
		Folder: f,
	}
}

func (a RatedTimeFolders) String() string {
	arr := make([]string, len(a))
	for k, v := range a {
		arr[k] = fmt.Sprintf("%v\t%v", v.Points, v.Folder.Path)
	}
	return strings.Join(arr, "\n")
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
