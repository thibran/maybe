package main

import (
	"time"
)

type sortRatedFolders func(a RatedFolders)

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
	Add(path string, t time.Time)          // Add new folder to the repo.
	Search(s string) (RatedFolder, error)  // Search for the key s in the repo
	Show(s string, n int) (a RatedFolders) // Show returns n RatedFolders.
	Saver
	Loader
}
