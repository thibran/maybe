package main

import (
	"fmt"
	"os"
	"os/user"
	"testing"
	"time"
)

func TestSave(t *testing.T) {
	tmpFile := os.TempDir() + "/" + "miaow_save"
	r := NewFileRepo(tmpFile)
	r.Add("/foo/bar", time.Now())
	if err := r.Save(); err != nil {
		t.Error(err)
	}
}

func TestLoad(t *testing.T) {
	tmpFile := os.TempDir() + "/" + "miaow_load"
	r := NewFileRepo(tmpFile)
	r.Add("/foo/zot", time.Now())
	r.Save()
	r = NewFileRepo(tmpFile)
	r.Load()
	if len(r.m) != 1 {
		t.Fail()
	}
	fmt.Println(r.m)
}

func TestHomeDir(t *testing.T) {
	user, err := user.Current()
	if err != nil {
		t.Error(err)
	}

	cfgDir := user.HomeDir + "/.local/share/miaow"

	if err := os.MkdirAll(cfgDir, 0777); err != nil {
		panic(err)
	}

	fmt.Println(cfgDir)
}
