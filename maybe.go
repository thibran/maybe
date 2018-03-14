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

	"github.com/thibran/maybe/pref"
	"github.com/thibran/maybe/rated"
	"github.com/thibran/maybe/rated/folder"
	"github.com/thibran/maybe/repo"
	"github.com/thibran/maybe/util"
)

const appVersion = "0.5.0"

func main() {
	p := pref.Parse()
	r := repo.New(filepath.Join(p.DataDir, "maybe.data"), p.MaxEntries)
	r.Load(p.DataDir)
	// version
	if p.Version {
		handleVersion(r, p.DataDir)
		return
	}
	// init
	if p.Init {
		handleInit(r, p.HomeDir)
		return
	}
	// add path
	if p.Add != "" {
		handleAdd(r, p.Add)
		return
	}
	// search
	if p.Search.IsNotEmpty() {
		handleSearch(r, p.Search)
		return
	}
	// list
	if p.List.IsNotEmpty() {
		handleList(r, p.List)
		return
	}
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(1)
}

func handleVersion(r *repo.Repo, dataDir string) {
	fmt.Printf("maybe %s   entries: %d   %s\n",
		appVersion, r.Size(), runtime.Version())

	if pref.Verbose {
		fmt.Printf("\nDataDir: %s\n", dataDir)
	}
}

func handleInit(r *repo.Repo, homeDir string) {
	r.Walk(homeDir)
	if err := r.Save(); err != nil {
		log.Fatalf("handleInit failed with: %v\n", err)
	}
	fmt.Println("entries:", r.Size())
}

func handleAdd(r *repo.Repo, path string) {
	if strings.TrimSpace(path) == "" {
		return
	}
	r.Add(path, time.Now())
	if err := r.Save(); err != nil {
		log.Fatalf("handleAdd - path: %s\n", err)
	}
}

func handleSearch(r *repo.Repo, q pref.Query) {
	// return path-query directly
	if q.Start == "" && strings.HasPrefix(q.Last, "/") {
		fmt.Println(q.Last)
		return
	}
	rf, err := r.Search(folder.CheckerFn(), q)
	if err != nil {
		// all okay
		if err == repo.ErrNoResult {
			os.Exit(1)
		}
		// hell should freez
		util.Logln(err)
		os.Exit(2)
	}
	fmt.Printf(rf.Path)
}

func handleList(r *repo.Repo, q pref.Query) {
	a := r.List(q, true)
	if len(a) == 0 {
		return
	}
	var res []string
	res = append(res, util.NormalOrVerbose("Rating\tFolder", "Time\tText\tFolder"))
	pathExistFn := folder.CheckerFn()
	appendFn := func(rf *rated.Rated) {
		res = append(res, util.NormalOrVerbose(
			fmt.Sprintf("%d\t%s", rf.Points(), rf.Path),
			fmt.Sprintf("%d\t%d\t%s", rf.TimePoints,
				rf.SimilarityPoints, rf.Path)))
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
