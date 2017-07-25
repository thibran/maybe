package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"
)

var (
	f1 = NewFolder("/home/foo",
		time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
	)

	f2 = NewFolder("/home/tux",
		time.Date(2012, time.February, 3, 11, 30, 0, 0, time.UTC),
	)

	f3 = NewFolder("/etc/apt",
		time.Date(2016, time.March, 20, 18, 0, 0, 0, time.UTC),
	)
)

func TestAdd_updateExisting(t *testing.T) {
	// verbose = true
	r := NewFileRepo("/baz/bar/zot", 10)
	now := time.Now()
	folder := NewFolder("/home/foo", time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC))
	r.m = map[string]Folder{folder.Path: folder}
	r.Add(f1.Path, now)

	f := r.m[f1.Path]
	if f.Count != 2 {
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
	r := NewFileRepo("/baz/bar/zot", 10)
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
	now := time.Now()
	r := NewFileRepo("/baz/bar/zot", 10)
	r.updateOrAddPath("/home/nfoo", now.Add(-time.Second*40), false)
	r.updateOrAddPath("/home/foo", now.Add(-time.Hour*18), false)
	r.updateOrAddPath("/etc/apt", now.Add(-time.Hour*24*7*2), false)
	a := search(r.m, "foo", func(a RatedFolders) { a.sort() })
	if len(a) == 0 {
		t.Fatal()
	}
	if a[0].points() < a[1].points() {
		t.Fatal()
	}
	emptyMap := map[string]Folder{}
	a = search(emptyMap, "foo", func(a RatedFolders) { a.sort() })
	if len(a) != 0 {
		t.Fatalf("empty map shoudl return no result, got %v", len(a))
	}
}

func TestCheckFolder(t *testing.T) {
	if checkFolder("/zot/baz_faz/moo") {
		t.Fatal()
	}
	if checkFolder("") {
		t.Fatal()
	}
}

func TestSave(t *testing.T) {
	// verbose = true
	tmp, err := ioutil.TempFile("", "maybe.data_")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(tmp.Name())
	r := NewFileRepo(tmp.Name(), 10)
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
	r := NewFileRepo(tmp.Name(), 10)
	r.Save()
	if err := r.Load(); err != nil {
		t.Fatal(err)
	}
	r = NewFileRepo("/zot/foo/abababa/bar", 1)
	if err := r.Load(); err != errNoFile {
		t.Fatal()
	}
}

func TestSaveGzip(t *testing.T) {
	// verbose = true
	var buf bytes.Buffer
	m := map[string]Folder{"/foo": NewFolder("/foo", time.Now())}
	if err := saveGzip(&buf, m); err != nil {
		t.Fatal(err)
	}
}

func TestLoadGzip(t *testing.T) {
	// verbose = true
	var buf bytes.Buffer
	m := map[string]Folder{"/foo": NewFolder("/foo", time.Now())}
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
	r := NewFileRepo("/baz/bar/zot", 10)
	r.updateOrAddPath("/home/nfoo", now.Add(-time.Second*40), false)
	r.updateOrAddPath("/home/foo", now.Add(-time.Hour*18), false)
	r.updateOrAddPath("/etc/apt", now.Add(-time.Hour*24*7*2), false)
	a := r.Show("foo", 2)
	if len(a) != 2 {
		t.Fatal()
	}
	if a[0].folder.Path != "/home/foo" {
		t.Fatal()
	}
	if a[1].folder.Path != "/home/nfoo" {
		t.Fatal()
	}
}

func TestAdd(t *testing.T) {
	//verbose = true
	r := NewFileRepo("/baz/bar/zot", 10)
	r.Add("/tmp/zot/hot", time.Now())
	r.Add("/tmp/zot", time.Now())
	if len(r.m) != 3 {
		t.Fatalf("len(r.m) shoud be 3, got %v", len(r.m))
	}
}

func TestUpdateOrAddPath(t *testing.T) {
	// verbose = true
	keepEntries := 2
	r := NewFileRepo("/baz/bar/zot", keepEntries)
	r.updateOrAddPath("/zot", time.Now(), false)
	r.updateOrAddPath("/bar", time.Now(), false)
	r.updateOrAddPath("/foo", time.Now(), false)
	if len(r.m) != keepEntries {
		t.Fatalf("expected %d, got %v", keepEntries, len(r.m))
	}
}
