// +build !android !darwin !windows

package pref

import (
	"flag"
	"fmt"
	"log"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
)

// Verbose output
var Verbose = false

const (
	maxEntries    = 10000
	minMaxEntries = 200 // minimal value for the maxEntries variable
)

// Pref object.
type Pref struct {
	DataDir, HomeDir, Add string
	List, Search          Query
	Version, Init         bool
	MaxEntries            int
}

// Parse flags.
func Parse() Pref {
	homeDir := userHome()
	dataDir := filepath.Join(homeDir, ".local/share/maybe")

	var p Pref
	p.HomeDir = homeDir
	flagDatadirVar(&p.DataDir, "datadir", dataDir, "")
	flag.StringVar(&p.Add, "add", "", "add path to index")
	q := flag.String("search", "", "search for keyword")
	l := flag.String("list", "", "list results for keyword")
	flag.BoolVar(&p.Init, "init", false, "scan $HOME and add folders (six folder-level deep)")
	flag.BoolVar(&p.Version, "version", false, "print maybe version")
	flagMaxentriesVar(&p.MaxEntries, "max-entries", maxEntries, "maximum unique path-entries")
	verb := flag.Bool("v", false, "verbose")
	flag.Parse()

	p.Search = queryFrom(*q)
	p.List = queryFrom(*l)

	Verbose = *verb
	return p
}

// Query object
type Query struct {
	Start, Last string
}

// IsNotEmpty returns true if Query.Last contains data.
func (q *Query) IsNotEmpty() bool { return len(q.Last) > 0 }

func (q Query) String() string {
	return fmt.Sprintf("{start: %s  last: %s}", q.Start, q.Last)
}

func queryFrom(s string) Query {
	s = strings.TrimSpace(s)
	if s == "" {
		return Query{}
	}
	if arg := flag.Args(); len(arg) > 0 {
		return Query{Start: s, Last: arg[0]}
	}
	return Query{Last: s}
}

func userHome() string {
	user, err := user.Current()
	if err != nil {
		log.Fatalf("current user unknown: %v\n", err)
	}
	if strings.TrimSpace(user.HomeDir) == "" {
		log.Fatal("user home directory path is empty")
	}
	return user.HomeDir
}

type maxentries int

func flagMaxentriesVar(p *int, name string, value int, usage string) {
	*p = value
	flag.CommandLine.Var((*maxentries)(p), name, usage)
}

func (m *maxentries) String() string { return strconv.Itoa(int(*m)) }

func (m *maxentries) Set(s string) error {
	n, err := strconv.ParseInt(s, 0, 64)
	if n < minMaxEntries {
		n = minMaxEntries
	}
	*m = maxentries(n)
	return err
}

type datadir string

func flagDatadirVar(p *string, name string, value string, usage string) {
	*p = value
	flag.CommandLine.Var((*datadir)(p), name, usage)
}

func (m *datadir) String() string { return string(*m) }

func (m *datadir) Set(s string) error {
	if strings.TrimSpace(s) == "" {
		return fmt.Errorf("datadir empty or consists only of whitespace")
	}
	*m = datadir(s)
	return nil
}
