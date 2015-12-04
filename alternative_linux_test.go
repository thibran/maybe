package main

import (
	"fmt"
	"sort"
	"testing"
)

func TestInPathSegment(t *testing.T) {
	s := new(search)
	s.key = "down"
	arr := entries{
		{v: "/files/src/go/src/juju-core/downloader", search: s.key},
		{v: "/files/Downloads/programs/android", search: s.key},
		{v: "/files/src/ubuntu/godd/parts", search: s.key},
	}
	for _, r := range arr {
		res := s.inPathSegment(r.v)
		if len(res) == 0 {
			continue
		}
		if res != "/files/Downloads" {
			t.Fail()
		}
		//fmt.Println(res)
	}
}

func TestAlternative(t *testing.T) {
	arr := []string{
		"Downl",
		"u2u",
		"table",
		"Soundtrack",
		"down",
		"download",
	}
	for _, v := range arr {
		s := newSearch(v)
		fmt.Printf("%s:\t%s\n", v, s.alternative())
	}
}

func TestPrintResult(t *testing.T) {
	s := newSearch("down")
	s.printResultlist()
}

func TestFilter_keyDownload(t *testing.T) {
	search := "download"
	s := newSearch(search)
	arr := s.list()
	// arr := []string{
	// 	"/files/src/go/src/juju-core/downloader",
	// 	"/files/Downloads/programs/android-studio/gradle/gradle",
	// 	"/files/Music/Soundtrack/Black Hawk down",
	// }
	for _, r := range s.filter(arr, 10) {
		fmt.Printf("'%s'\n", r.v)
		//fmt.Println(r.v)
	}
}

func TestFilter_keyDown(t *testing.T) {
	search := "down"
	s := newSearch(search)
	arr := []string{
		"/files/src/go/src/juju-core/downloader",
		"/files/Downloads/programs/android-studio/gradle/gradle",
		"/files/Music/Soundtrack/Black Hawk down",
	}
	for _, r := range s.filter(arr, 10) {
		fmt.Printf("'%s'\n", r.v)
		//fmt.Println(r.v)
	}
}

func TestSort_keyTable(t *testing.T) {
	search := "table"
	arr := entries{
		{v: "/files/src/go/src/github.com/mattn/go-gtk/example/table", search: search},
		{v: "/files/src/go/src/thibaut/table", search: search},
		{v: "/files/src/ubuntu/godd/stage/etc/iproute2/rt_tables", search: search},
	}
	sort.Sort(arr)
	for _, r := range arr {
		fmt.Println(r.v)
	}
}

func TestHiddenFolderInPath(t *testing.T) {
	in := "/files/Downloads/.tools/android"
	if !hiddenFolderInPath(in) {
		t.Fail()
	}
}

func TestHiddenFolderInPath_negative(t *testing.T) {
	in := "/files/Downloads/tools/android"
	if hiddenFolderInPath(in) {
		t.Fail()
	}
}

func TestLooksLikeAFile(t *testing.T) {
	in := "/files/android.jpg"
	if !looksLikeAFile(in) {
		t.Fail()
	}
}

func TestLooksLikeAFile_negative(t *testing.T) {
	in := "/files/android"
	if looksLikeAFile(in) {
		t.Fail()
	}
}

func TestIsDir(t *testing.T) {
	in := "/files/src/go/src/thibaut/table"
	if !isDir(in) {
		t.Fail()
	}
}

func TestIsDir_negative(t *testing.T) {
	in := "/files/Desktop/ai_brainstrom.txt"
	if isDir(in) {
		t.Fail()
	}
}
