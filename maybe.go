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

func main() {
	// flag
	dataDir := flag.String("datadir", defaultDataDir(), "")
	add := flag.String("add", "", "add path")
	s := flag.String("search", "", "search for")
	version := flag.Bool("version", false, "print maybe version")
	flag.Parse()
	if *version {
		fmt.Printf("maybe 0.1.1   %s\n", runtime.Version())
		os.Exit(0)
	}
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
	// searching for someting?
	if len(*s) != 0 {
		handleSearch(r, *s)
		return
	}
	os.Exit(1)
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
