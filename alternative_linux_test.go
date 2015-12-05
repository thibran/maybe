package main

import (
	"fmt"
	"sort"
	"testing"
)

func TestAlternative(t *testing.T) {
	arr := []string{
		"Downl",
		"u2u",
		"table",
		"Soundtrack",
		"down",
		"download",
		"vide",
		"tmp",
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

func TestFilter_keyDl(t *testing.T) {
	search := "dl"
	s := newSearch(search)
	arr := s.list()
	// arr := []string{
	// 	"/files/src/go/src/juju-core/downloader",
	// 	"/files/Downloads/programs/android-studio/gradle/gradle",
	// 	"/files/Music/Soundtrack/Black Hawk down",
	// }
	for _, row := range s.filter(arr, 10) {
		fmt.Printf("'%s'\n", row)
	}
}

func TestFilter_keyVide(t *testing.T) {
	search := "vide"
	s := newSearch(search)
	arr := s.list()
	// arr := []string{
	// 	"/files/src/go/src/juju-core/downloader",
	// 	"/files/Downloads/programs/android-studio/gradle/gradle",
	// 	"/files/Music/Soundtrack/Black Hawk down",
	// }
	for _, row := range s.filter(arr, 10) {
		fmt.Printf("'%s'\n", row)
	}
}

func TestFilter_keyDownload(t *testing.T) {
	search := "download"
	s := newSearch(search)
	arr := s.list()
	for _, row := range s.filter(arr, 10) {
		fmt.Printf("'%s'\n", row)
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
	for _, row := range s.filter(arr, 10) {
		fmt.Printf("'%s'\n", row)
	}
}

func TestSort_keyTable(t *testing.T) {
	search := "table"
	en := entries{
		search: search,
		arr: []string{
			"/files/src/go/src/github.com/mattn/go-gtk/example/table",
			"/files/src/go/src/thibaut/table",
			"/files/src/ubuntu/godd/stage/etc/iproute2/rt_tables",
		},
	}
	sort.Sort(en)
	for _, row := range en.arr {
		fmt.Println(row)
	}
}

func TestPrefixFilter(t *testing.T) {
	prefix := "/home/tux"
	p := "/home/tux/.bin"
	if !prefixFilter(prefix, p, []string{"/home/tux/.bin"}) {
		t.Fail()
	}
}

func TestPrefixFilter_negative(t *testing.T) {
	prefix := "/home/tux"
	p := "/home/tux/.bin"
	if prefixFilter(prefix, p, []string{"/home/tux/.bin/old"}) {
		t.Fail()
	}
}

func TestInPathSegment(t *testing.T) {
	s := new(search)
	s.key = "down"
	en := entries{
		search: s.key,
		arr: []string{
			"/files/src/go/src/juju-core/downloader",
			"/files/Downloads/programs/android",
			"/files/src/ubuntu/godd/parts",
		},
	}
	for _, row := range en.arr {
		res := s.inPathSegment(row)
		if len(res) == 0 {
			continue
		}
		if res != "/files/Downloads" {
			t.Fail()
		}
		//fmt.Println(res)
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

func TestMaxDepthFilter(t *testing.T) {
	in := "/usr/share/android"
	if !maxDepthFilter("/usr", in, 2, nil) {
		t.Fail()
	}
}
