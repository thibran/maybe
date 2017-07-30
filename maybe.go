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

// TODO
// - write fish completion, using --show with a sub-command
//   http://fishshell.com/docs/current/index.html#completion-own
//   https://stackoverflow.com/questions/16657803/creating-autocomplete-script-with-sub-commands
//   https://github.com/fish-shell/fish-shell/issues/1217#issuecomment-31441757
// - multi-word search queries
// - maybe replace time rating with: fewer seconds from now > better
//   if a time value is not present, add penalty

const (
	appVersion    = "0.3.1"
	maxEntries    = 10000
	minMaxEntries = 200 // minimal value for the maxEntries variable
)

var verbose = false

type pref struct {
	dataDir    string
	add        string
	search     string
	list       string
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
	flag.StringVar(&p.list, "list", "", "list results for keyword")
	flag.BoolVar(&p.version, "version", false, "print maybe version")
	flag.IntVar(&p.maxEntries, "max-entries", maxEntries, "Maximum number of unique path-entries.")
	verb := flag.Bool("v", false, "print verbose info about app execution")
	flag.Parse()
	if p.maxEntries < minMaxEntries {
		p.maxEntries = minMaxEntries
	}
	if strings.TrimSpace(p.dataDir) == "" {
		log.Fatalf("datadir empty or consists only of whitespace")
	}
	verbose = *verb
	return p
}

func main() {
	p := parse()
	r := NewRepo(p.dataDir+"/maybe.data", p.maxEntries)
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
	if len(p.list) != 0 {
		handleList(r, p.list)
		return
	}
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(1)
}

func handleVersion(r *Repo, dataDir string) {
	fmt.Printf("maybe %s   entries: %d   %s\n",
		appVersion, r.Size(), runtime.Version())

	if verbose {
		fmt.Printf("\nDataDir: %s\n", dataDir)
	}
}

func handleList(r *Repo, show string) {
	a := r.Show(show, 10)
	if len(a) == 0 {
		return
	}
	fmt.Println("Points\tFolder")
	for _, rf := range a {
		fmt.Printf("%d\t%s\n", rf.points(), rf.Path)
	}
}

func handleSearch(r *Repo, query string) {
	// return path-query directly
	if strings.HasPrefix(query, "/") {
		fmt.Println(query)
		return
	}
	rf, err := r.Search(folderChecker(), query)
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

func userHome() string {
	user, err := user.Current()
	if err != nil {
		log.Fatalf("current user unknown: %v\n", err)
	}
	return user.HomeDir
}
