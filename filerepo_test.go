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
	f1 = NewFolder("/home/foo", 1,
		time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
	)

	f2 = NewFolder("/home/tux", 1,
		time.Date(2012, time.February, 3, 11, 30, 0, 0, time.UTC),
	)

	f3 = NewFolder("/etc/apt", 1,
		time.Date(2016, time.March, 20, 18, 0, 0, 0, time.UTC),
	)
)

func TestAdd_updateExisting(t *testing.T) {
	// verbose = true
	r := NewFileRepo("/baz/bar/zot", 10)
	now := time.Now()
	folder := NewFolder("/home/foo", 1, time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC))
	r.m = map[string]Folder{folder.Path: folder}
	r.Add(f1.Path, now)

	f := r.m[f1.Path]
	if f.Count != 2 {
		t.Fail()
	}
	if f.Times[0] != now {
		t.Error("Times[0] should be equals time now.")
	}
	if len(f.Times) > MaxTimeEntries {
		t.Fail()
	}
}

func TestAdd_ignoreFolders(t *testing.T) {
	// verbose = true
	r := NewFileRepo("/baz/bar/zot", 10)
	r.Add("/tmp/.git", time.Now())
	if _, ok := r.m["/tmp"]; !ok {
		t.Fail()
	}
	if len(r.m) != 1 {
		t.Fail()
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
		t.Fail()
	}
	if a[0].points() < a[1].points() {
		t.Fail()
	}
}

func TestCheckFolder(t *testing.T) {
	if checkFolder("/zot/baz_faz/moo") {
		t.Fail()
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
		t.Error(err)
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
		t.Error(err)
	}
	r = NewFileRepo("/zot/foo/abababa/bar", 1)
	if err := r.Load(); err != errNoFile {
		t.Fail()
	}
}

func TestSaveGzip(t *testing.T) {
	// verbose = true
	var buf bytes.Buffer
	m := map[string]Folder{"/foo": NewFolder("/foo", 1, time.Now())}
	if err := saveGzip(&buf, m); err != nil {
		t.Error(err)
	}
}

func TestLoadGzip(t *testing.T) {
	// verbose = true
	var buf bytes.Buffer
	m := map[string]Folder{"/foo": NewFolder("/foo", 1, time.Now())}
	saveGzip(&buf, m)
	m2, err := loadGzip(&buf)
	if err != nil {
		t.Error(err)
	}
	if _, ok := m2["/foo"]; !ok {
		t.Fail()
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
		t.Fail()
	}
	if a[0].folder.Path != "/home/foo" {
		t.Fail()
	}
	if a[1].folder.Path != "/home/nfoo" {
		t.Fail()
	}
}

func TestAdd(t *testing.T) {
	//verbose = true
	r := NewFileRepo("/baz/bar/zot", 10)
	r.Add("/tmp/zot/hot", time.Now())
	r.Add("/tmp/zot", time.Now())
	if len(r.m) != 3 {
		t.Fail()
	}
}
