package main

import (
	"flag"
	"fmt"
	"os"
	"os/user"
	"time"
)

// TODO
// Repo:
// 	max folder entries
//  history()  // return last 10 Folder objects

func main() {
	// flag
	dataDir := flag.String("datadir", defaultDataDir(), "")
	in := flag.String("add", "", "add path")
	s := flag.String("search", "", "search for")
	flag.Parse()
	// create data dir, if not existent
	if err := os.MkdirAll(*dataDir, 0777); err != nil {
		panic(err)
	}
	// load data
	var r Repo
	r = NewFileRepo(*dataDir + "/miaow.data")
	r.Load()
	// add path to miaow?
	if len(*in) != 0 {
		r.Add(*in, time.Now())
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
	return user.HomeDir + "/.local/share/miaow"
}
