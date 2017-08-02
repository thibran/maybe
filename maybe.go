// +build !android !darwin !windows

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const (
	appVersion    = "0.3.1"
	maxEntries    = 10000
	minMaxEntries = 200 // minimal value for the maxEntries variable
)

var verbose = false

func main() {
	p := parse()
	r := NewRepo(filepath.Join(p.dataDir, "maybe.data"), p.maxEntries)
	loadData(r, p.dataDir)
	// version
	if p.version {
		handleVersion(r, p.dataDir)
		return
	}
	// init
	if p.init {
		handleInit(r, p.homeDir)
		return
	}
	// add path
	if p.add != "" {
		handleAdd(r, p.add)
		return
	}
	// search
	if p.search.isNotEmpty() {
		handleSearch(r, p.search)
		return
	}
	// list
	if p.list.isNotEmpty() {
		handleList(r, p.list)
		return
	}
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(1)
}

func loadData(r *Repo, dataDir string) {
	if err := r.Load(); err != nil {
		if err != errNoFile {
			log.Fatalln(err)
		}
		// create data dir, if not existent
		if err := os.MkdirAll(dataDir, 0770); err != nil {
			log.Fatalf("main - create data dir: %s\n", err)
		}
	}
}

func handleVersion(r *Repo, dataDir string) {
	fmt.Printf("maybe %s   entries: %d   %s\n",
		appVersion, r.Size(), runtime.Version())

	if verbose {
		fmt.Printf("\nDataDir: %s\n", dataDir)
	}
}

func handleInit(r *Repo, homeDir string) {
	r.Walk(homeDir)
	if err := r.Save(); err != nil {
		log.Fatalf("handleInit failed with: %v\n", err)
	}
}

func handleAdd(r *Repo, path string) {
	if strings.TrimSpace(path) == "" {
		return
	}
	r.Add(path, time.Now())
	if err := r.Save(); err != nil {
		log.Fatalf("handleAdd - path: %s\n", err)
	}
}

func handleList(r *Repo, q query) {
	a := r.List(q, 8)
	if len(a) == 0 {
		return
	}
	if !verbose {
		fmt.Println("Rating\tFolder")
	} else {
		fmt.Println("Time\tText\tFolder")
	}
	for _, rf := range a {
		if !verbose {
			fmt.Printf("%d\t%s\n", rf.timePoints+rf.similarityPoints, rf.Path)
		} else {
			fmt.Printf("%d\t%d\t%s\n", rf.timePoints, rf.similarityPoints, rf.Path)
		}
	}
}

func handleSearch(r *Repo, q query) {
	// return path-query directly
	if q.start == "" && strings.HasPrefix(q.last, "/") {
		fmt.Println(q.last)
		return
	}
	rf, err := r.Search(folderChecker(), q)
	if err != nil {
		// all okay
		if err == errNoResult {
			os.Exit(1)
		}
		// hell should freez
		logln(err)
		os.Exit(2)
	}
	fmt.Println(rf.Path)
}
