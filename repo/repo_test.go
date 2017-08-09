package repo

import (
	"math"
	"strings"
	"testing"
	"thibaut/maybe/pref"
	"thibaut/maybe/rated"
	"thibaut/maybe/rated/folder"
	"time"
)

var _ ResourceChecker = (*folder.ResourceCheckerFn)(nil)

func TestAdd(t *testing.T) {
	//verbose = true
	r := New("/baz/bar/zot", 10)
	r.Add("/tmp/zot/hot", time.Now())
	r.Add("/tmp/zot", time.Now())
	if len(r.m) != 3 {
		t.Fatalf("exp 3, got %v", len(r.m))
	}
}

func TestAdd_ignoreFolders(t *testing.T) {
	// verbose = true
	r := New("/baz/bar/zot", 10)
	r.Add("/tmp/.git", time.Now())
	if _, ok := r.m["/tmp"]; !ok {
		t.Fatal()
	}
	if len(r.m) != 1 {
		t.Fatalf("len(r.m) should be 1, got %v", len(r.m))
	}
}

func TestAdd_updateExisting(t *testing.T) {
	// verbose = true
	r := New("/baz/bar/zot", 10)
	now := time.Now()
	fo := folder.New("/home/foo", time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC))
	r.m = rated.Map{fo.Path: fo}
	r.Add(fo.Path, now)

	f := r.m[fo.Path]
	if f.UpdateCount != 2 {
		t.Fatal()
	}
	if f.Times[0] != now {
		t.Fatal("Times[0] should be equals time now.")
	}
	if len(f.Times) > rated.MaxTimeEntries {
		t.Fatal()
	}
}

func TestUpdateOrAdd(t *testing.T) {
	// verbose = true
	tt := []struct {
		name                   string
		maxEntries, expEntries int
		paths                  []string
		overflow               bool
		expCount               uint32
	}{
		{name: "keep newest", maxEntries: 2, expEntries: 2, expCount: 1,
			paths: []string{"/zot", "/bar", "/foo"}},

		{name: "counter overflow", maxEntries: 2, expEntries: 1,
			expCount: uint32(math.MaxUint32),
			overflow: true, paths: []string{"/zot"}},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			r := New("/baz/bar/zot", tc.maxEntries)
			for _, p := range tc.paths {
				r.updateOrAdd(p, time.Now(), false)
			}
			if len(r.m) != tc.expEntries {
				t.Fatalf("exp expEntries %d, got %v", tc.expEntries, len(r.m))
			}
			if tc.overflow {
				f := r.m["/zot"]
				f.UpdateCount = tc.expCount
			}
			for k, f := range r.m {
				if f.UpdateCount != tc.expCount {
					t.Errorf("UpdateCount for %s should be %v, got %v",
						k, tc.expCount, f.UpdateCount)
				}
			}
		})
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
		name, search, exp string
		paths             []path
	}{
		{name: "okay", search: "foo", exp: "/home/foo",
			paths: []path{
				{p: "/home/nfoo", t: now.Add(-time.Second * 40)},
				{p: "/home/foo", t: now.Add(-time.Hour * 18)},
				{p: "/etc/apt", t: now.Add(-time.Hour * 24)},
			},
		},
		{name: "not found", search: "zzz",
			paths: []path{
				{p: "/home/nfoo", t: now.Add(-time.Second * 40)},
				{p: "/home/foo", t: now.Add(-time.Hour * 18)},
				{p: "/etc/apt", t: now.Add(-time.Hour * 24)},
			},
		},
		{name: "no map entries", search: "foo"},
	}
	doesExist := func(path string) bool {
		return strings.TrimSpace(path) != ""
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			r := New("/baz/bar/zot", 10)
			for _, p := range tc.paths {
				r.updateOrAdd(p.p, p.t, false)
			}
			rf, err := r.Search(folder.ResourceCheckerFn(
				doesExist), pref.Query{Last: tc.search})
			if err != nil && tc.exp != "" {
				t.Fatalf("exp %q, got %v", tc.exp, err)
			}
			if rf != nil && rf.Path != tc.exp {
				t.Fatalf("exp %q, got %q", tc.exp, rf.Path)
			}
			if tc.exp == "" && err != ErrNoResult {
				t.Fatalf("should be errNoResult, go %v", err)
			}
		})
	}
}

func TestList(t *testing.T) {
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
		name, exp, search string
		index, resLen     int
	}{
		{name: "okay 1", search: "foo", exp: "/home/foo",
			index: 0, resLen: 3},
		{name: "okay 2", search: "foo", exp: "/bbbbb/foo",
			index: 1, resLen: 3},
		{name: "one result", search: "apt", exp: "/etc/apt",
			index: 0, resLen: 1},
		{name: "no result", search: "zot", exp: "",
			index: 0, resLen: 0},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			r := New("/baz/bar/zot", 10)
			for _, p := range paths {
				r.updateOrAdd(p.p, p.t, false)
			}
			a := r.List(pref.Query{Last: tc.search}, false)
			if len(a) != tc.resLen {
				t.Fatalf("len(a) should be %v, got %v", tc.resLen, len(a))
			}
			if tc.exp != "" && a[tc.index].Path != tc.exp {
				t.Fatalf("should be %q, got %q", tc.exp, a[tc.index].Path)
			}
		})
	}
}
