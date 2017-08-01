package main

import (
	"flag"
	"fmt"
	"log"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
)

type pref struct {
	dataDir    string
	homeDir    string
	add        string
	search     query
	list       query
	init       bool
	version    bool
	maxEntries int
}

type query struct {
	start, last string
}

func (q *query) isNotEmpty() bool { return len(q.last) > 0 }

func (q query) String() string {
	return fmt.Sprintf("{start: %s  last: %s}", q.start, q.last)
}

func parse() pref {
	homeDir := userHome()
	dataDir := filepath.Join(homeDir, ".local/share/maybe")

	var p pref
	p.homeDir = homeDir
	flagDatadirVar(&p.dataDir, "datadir", dataDir, "")
	flag.StringVar(&p.add, "add", "", "add path to index")
	q := flag.String("search", "", "search for keyword")
	l := flag.String("list", "", "list results for keyword")
	flag.BoolVar(&p.init, "init", false, "scan $HOME and add folders (six folder-level deep)")
	flag.BoolVar(&p.version, "version", false, "print maybe version")
	flagMaxentriesVar(&p.maxEntries, "max-entries", maxEntries, "maximum unique path-entries")
	verb := flag.Bool("v", false, "verbose")
	flag.Parse()

	p.search = queryFrom(*q)
	p.list = queryFrom(*l)

	verbose = *verb
	return p
}

func queryFrom(s string) query {
	s = strings.TrimSpace(s)
	if s == "" {
		return query{}
	}
	if arg := flag.Args(); len(arg) > 0 {
		return query{start: s, last: arg[0]}
	}
	return query{last: s}
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
