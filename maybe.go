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
	appVersion    = "0.3.2"
	maxEntries    = 10000
	minMaxEntries = 200 // minimal value for the maxEntries variable
)

var verbose = false

func main() {
	p := parse()
	r := NewRepo(filepath.Join(p.dataDir, "maybe.data"), p.maxEntries)
	r.LoadData(p.dataDir)
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

// LoadData from dataDir or create directory.
func (r *Repo) LoadData(dataDir string) {
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
	fmt.Println("entries:", r.Size())
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

func handleList(r *Repo, q query) {
	a := r.List(q, true)
	if len(a) == 0 {
		return
	}
	var res []string
	res = append(res, normalOrVerbose("Rating\tFolder", "Time\tText\tFolder"))
	pathExistFn := folderChecker()
	appendFn := func(rf *RatedFolder) {
		res = append(res, normalOrVerbose(
			fmt.Sprintf("%d\t%s", rf.points(), rf.Path),
			fmt.Sprintf("%d\t%d\t%s", rf.timePoints,
				rf.similarityPoints, rf.Path)))
	}
	entries := 0
	entryLimit := 8
	for _, rf := range a {
		if entries == entryLimit {
			break
		}
		if !pathExistFn(rf.Path) {
			continue
		}
		appendFn(rf)
		entries++
	}
	if len(res) == 1 {
		return
	}
	fmt.Println(strings.Join(res, "\n"))
}
