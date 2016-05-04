package main

import (
	"flag"
	"fmt"
	"os"
	"os/user"
	"runtime"
	"time"
)

// TODO
// Repo:
// 	max folder entries
//  history()  // return last 10 Folder objects
//  don't show result if under a certain limit

var maxEntries = 300 // Maximum number of history entries to keep.
const minMaxEntries = 50

func main() {
	// flag
	dataDir := flag.String("datadir", defaultDataDir(), "")
	add := flag.String("add", "", "add path")
	search := flag.String("search", "", "search for")
	show := flag.String("show", "", "show results for")
	version := flag.Bool("version", false, "print maybe version")
	entries := flag.Int("max-entries", maxEntries, "Maximum number of unique saved path-entries (minimum 50).")
	flag.Parse()
	if *version {
		fmt.Printf("maybe 0.1.2   %s\n", runtime.Version())
		os.Exit(0)
	}
	if *entries < minMaxEntries {
		*entries = minMaxEntries
	}
	maxEntries = *entries
	// create data dir, if not existent
	if err := os.MkdirAll(*dataDir, 0770); err != nil {
		panic(err)
	}
	// load data
	var r Repo
	r = NewFileRepo(*dataDir + "/maybe.data")
	r.Load()
	// add path to maybe?
	if len(*add) != 0 {
		r.Add(*add, time.Now())
		if err := r.Save(); err != nil {
			panic(err)
		}
		return
	}
	// show flag
	if len(*show) != 0 {
		handleShow(r, *show)
		return
	}
	// searching for someting?
	if len(*search) != 0 {
		handleSearch(r, *search)
		return
	}
	os.Exit(1)
}

func handleShow(r Repo, show string) {
	a := r.Show(show, 10)
	if len(a) == 0 {
		return
	}
	fmt.Println("Points\tFolder")
	for _, rf := range a {
		fmt.Printf("%d\t%s\n", rf.Points, rf.Folder.Path)
	}
}

func handleSearch(r Repo, search string) {
	rf, err := r.Search(search)
	if err != nil {
		// no result found
		os.Exit(1)
	}
	fmt.Println(rf.Folder.Path)
}

func defaultDataDir() string {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	return user.HomeDir + "/.local/share/maybe"
}
