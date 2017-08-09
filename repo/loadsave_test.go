package repo

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"testing"
	"thibaut/maybe/ratedfolder"
	"thibaut/maybe/ratedfolder/folder"
	"time"
)

func TestSave(t *testing.T) {
	// verbose = true
	tmp, err := ioutil.TempFile("", "maybe.data_")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(tmp.Name())
	r := New(tmp.Name(), 10)
	if err := r.Save(); err != nil {
		t.Fatal(err)
	}
}

func TestSaveGzip(t *testing.T) {
	// verbose = true
	var buf bytes.Buffer
	m := ratedfolder.Map{"/foo": folder.New("/foo", time.Now())}
	if err := saveGzip(&buf, m); err != nil {
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
	r := New(tmp.Name(), 10)
	r.Save()
	if err := r.Load(); err != nil {
		t.Fatal(err)
	}
	r = New("/zot/foo/abababa/bar", 1)
	if err := r.Load(); err != errNoFile {
		t.Fatal()
	}
}

func TestLoadGzip(t *testing.T) {
	// verbose = true
	var buf bytes.Buffer
	m := ratedfolder.Map{"/foo": folder.New("/foo", time.Now())}
	saveGzip(&buf, m)
	m2, err := loadGzip(&buf)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := m2["/foo"]; !ok {
		t.Fatal()
	}
}
