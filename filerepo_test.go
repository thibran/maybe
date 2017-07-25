package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"math"
	"os"
	"strings"
	"testing"
	"time"
)

func TestAdd_updateExisting(t *testing.T) {
	// verbose = true
	r := NewRepo("/baz/bar/zot", 10)
	now := time.Now()
	folder := NewFolder("/home/foo", time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC))
	r.m = FolderMap{folder.Path: folder}
	r.Add(folder.Path, now)

	f := r.m[folder.Path]
	if f.UpdateCount != 2 {
		t.Fatal()
	}
	if f.Times[0] != now {
		t.Fatal("Times[0] should be equals time now.")
	}
	if len(f.Times) > MaxTimeEntries {
		t.Fatal()
	}
}

func TestAdd_ignoreFolders(t *testing.T) {
	// verbose = true
	r := NewRepo("/baz/bar/zot", 10)
	r.Add("/tmp/.git", time.Now())
	if _, ok := r.m["/tmp"]; !ok {
		t.Fatal()
	}
	if len(r.m) != 1 {
		t.Fatalf("len(r.m) should be 1, got %v", len(r.m))
	}
}

func TestSearch(t *testing.T) {
	// verbose = true
	type path struct {
		p string
		t time.Time
	}
	now := time.Now()
	tt := []struct {
		name, query, exp string
		paths            []path
	}{
		{name: "okay", query: "foo", exp: "/home/foo",
			paths: []path{
				{p: "/home/nfoo", t: now.Add(-time.Second * 40)},
				{p: "/home/foo", t: now.Add(-time.Hour * 18)},
				{p: "/etc/apt", t: now.Add(-time.Hour * 24)},
			},
		},
		{name: "not found", query: "zzz",
			paths: []path{
				{p: "/home/nfoo", t: now.Add(-time.Second * 40)},
				{p: "/home/foo", t: now.Add(-time.Hour * 18)},
				{p: "/etc/apt", t: now.Add(-time.Hour * 24)},
			},
		},
		{name: "no map entries", query: "foo"},
	}
	doesExist := func(path string) bool {
		return strings.TrimSpace(path) != ""
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			r := NewRepo("/baz/bar/zot", 10)
			for _, p := range tc.paths {
				r.updateOrAddPath(p.p, p.t, false)
			}
			rf, err := r.Search(ResourceCheckerFn(doesExist), tc.query)
			if err != nil && tc.exp != "" {
				t.Fatalf("exp %q, got %v", tc.exp, err)
			}
			if rf.Path != tc.exp {
				t.Fatalf("exp %q, got %q", tc.exp, rf.Path)
			}
			if tc.exp == "" && err != errNoResult {
				t.Fatalf("should be errNoResult, go %v", err)
			}
		})
	}
}

func TestSave(t *testing.T) {
	// verbose = true
	tmp, err := ioutil.TempFile("", "maybe.data_")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(tmp.Name())
	r := NewRepo(tmp.Name(), 10)
	if err := r.Save(); err != nil {
		t.Fatal(err)
	}
}

func TestLoad(t *testing.T) {
	// verbose = true
	tmp, err := ioutil.TempFile("", "maybe.data_")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(tmp.Name())
	r := NewRepo(tmp.Name(), 10)
	r.Save()
	if err := r.Load(); err != nil {
		t.Fatal(err)
	}
	r = NewRepo("/zot/foo/abababa/bar", 1)
	if err := r.Load(); err != errNoFile {
		t.Fatal()
	}
}

func TestSaveGzip(t *testing.T) {
	// verbose = true
	var buf bytes.Buffer
	m := FolderMap{"/foo": NewFolder("/foo", time.Now())}
	if err := saveGzip(&buf, m); err != nil {
		t.Fatal(err)
	}
}

func TestLoadGzip(t *testing.T) {
	// verbose = true
	var buf bytes.Buffer
	m := FolderMap{"/foo": NewFolder("/foo", time.Now())}
	saveGzip(&buf, m)
	m2, err := loadGzip(&buf)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := m2["/foo"]; !ok {
		t.Fatal()
	}
}

func TestShow(t *testing.T) {
	// verbose = true
	now := time.Now()
	type path struct {
		p string
		t time.Time
	}
	paths := []path{
		{p: "/home/nfoo", t: now.Add(-time.Second)},
		{p: "/home/foo", t: now.Add(-time.Hour * 10)},
		{p: "/etc/apt", t: now.Add(-time.Hour * 24)},
		{p: "/bbbbb/foo", t: now.Add(-time.Hour * 14)},
	}
	tt := []struct {
		name, exp, query     string
		index, limit, resLen int
	}{
		{name: "okay 1", query: "foo", exp: "/home/foo",
			index: 0, limit: 2, resLen: 2},
		{name: "okay 2", query: "foo", exp: "/bbbbb/foo",
			index: 1, limit: 2, resLen: 2},
		{name: "no result", query: "foo", exp: "",
			index: 0, limit: 0, resLen: 0},
		{name: "one result", query: "apt", exp: "/etc/apt",
			index: 0, limit: 3, resLen: 1},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			r := NewRepo("/baz/bar/zot", 10)
			for _, p := range paths {
				r.updateOrAddPath(p.p, p.t, false)
			}
			a := r.Show(tc.query, tc.limit)
			if len(a) != tc.resLen {
				t.Fatalf("len(a) should be %v, got %v", tc.resLen, len(a))
			}
			if tc.exp != "" && a[tc.index].Path != tc.exp {
				t.Fatalf("should be %q, got %q", tc.exp, a[tc.index].Path)
			}
		})
	}
}

func TestAdd(t *testing.T) {
	//verbose = true
	r := NewRepo("/baz/bar/zot", 10)
	r.Add("/tmp/zot/hot", time.Now())
	r.Add("/tmp/zot", time.Now())
	if len(r.m) != 3 {
		t.Fatalf("exp 3, got %v", len(r.m))
	}
}

func TestUpdateOrAddPath(t *testing.T) {
	// verbose = true
	keepEntries := 2
	r := NewRepo("/baz/bar/zot", keepEntries)
	r.updateOrAddPath("/zot", time.Now(), false)
	r.updateOrAddPath("/bar", time.Now(), false)
	r.updateOrAddPath("/foo", time.Now(), false)
	if len(r.m) != keepEntries {
		t.Fatalf("expected %d, got %v", keepEntries, len(r.m))
	}
	// test f.UpdateCount overflow protection
	exp := uint32(math.MaxUint32)
	f := r.m["/zot"]
	f.UpdateCount = exp
	r.m["/zot"] = f
	r.updateOrAddPath("/zot", time.Now(), false)
	if f = r.m["/zot"]; f.UpdateCount != exp {
		t.Errorf("UpdateCount should %v, got %v", exp, f.UpdateCount)
	}
}
