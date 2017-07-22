package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// TODO:
//   exact match should win over containig match: timer_c vs timer_cc

// TODO:
//   initial start, add $HOME folders, $PATH, $TMP...

// TODO:
//   better recognize words in the middle like 'aaa' in 'b_aaa_c'!
//   don't show result if under a certain limit

const (
	appVersion    = "0.2.6"
	maxEntries    = 200
	minMaxEntries = 30 // minimal value for the maxEntries variable
)

var verbose = false

type pref struct {
	dataDir    string
	add        string
	search     string
	show       string
	version    bool
	maxEntries int
}

func parse() pref {
	userHome := userHome()
	dataDir := filepath.Join(userHome, ".local/share/maybe")

	var p pref
	flag.StringVar(&p.dataDir, "datadir", dataDir, "")
	flag.StringVar(&p.add, "add", "", "add path to maybe index")
	flag.StringVar(&p.search, "search", "", "search for keyword")
	flag.StringVar(&p.show, "show", "", "show results for keyword")
	flag.BoolVar(&p.version, "version", false, "print maybe version")
	flag.IntVar(&p.maxEntries, "max-entries", maxEntries, "Maximum number of unique path-entries (minimum 30).")
	verb := flag.Bool("v", false, "print verbose info about app execution")
	flag.Parse()
	if p.maxEntries < minMaxEntries {
		p.maxEntries = minMaxEntries
	}
	verbose = *verb
	return p
}

func main() {
	p := parse()
	r := NewFileRepo(p.dataDir+"/maybe.data", p.maxEntries)
	// load data
	if err := r.Load(); err != nil {
		if err != errNoFile {
			log.Fatalln(err)
		}
		// create data dir, if not existent
		if err := os.MkdirAll(p.dataDir, 0770); err != nil {
			log.Fatalf("main - create data dir: %s\n", err)
		}
	}
	// version
	if p.version {
		handleVersion(r, p.dataDir)
		os.Exit(0)
	}
	// add path
	if strings.TrimSpace(p.add) != "" {
		r.Add(p.add, time.Now())
		if err := r.Save(); err != nil {
			log.Fatalf("main - add path: %s\n", err)
		}
		return
	}
	// search
	if len(p.search) != 0 {
		handleSearch(r, p.search)
		return
	}
	// show
	if len(p.show) != 0 {
		handleShow(r, p.show)
		return
	}
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(1)
}

func handleVersion(r Repo, dataDir string) {
	fmt.Printf("maybe %s   entries: %d   %s\n",
		appVersion, r.Size(), runtime.Version())

	if verbose {
		fmt.Printf("\nDataDir: %s\n", dataDir)
	}
}

func handleShow(r Repo, show string) {
	a := r.Show(show, 10)
	if len(a) == 0 {
		return
	}
	fmt.Println("Points\tFolder")
	for _, rf := range a {
		fmt.Printf("%d\t%s\n", rf.points(), rf.folder.Path)
	}
}

func handleSearch(r Repo, search string) {
	rf, err := r.Search(search)
	if err != nil {
		// all okay
		if err == errNoResult {
			return
		}
		// hell should freez
		logln(err)
		os.Exit(2)
	}
	fmt.Println(rf.folder.Path)
}

func userHome() string {
	user, err := user.Current()
	if err != nil {
		log.Fatalf("current user unknown: %v\n", err)
	}
	return user.HomeDir
}
