package main

import (
	"testing"
	"time"
)

func TestSort(t *testing.T) {
	a := Times{f1.Times[0], f2.Times[0], f3.Times[0]}
	a = a.sort()
	if a[0] != f3.Times[0] {
		t.Fail()
	}
	if a[1] != f2.Times[0] {
		t.Fail()
	}
}

func TestCompare_equals(t *testing.T) {
	s := "foo"
	f1 := &Folder{
		Path:  "/home/foo",
		Count: 1,
		Times: Times{
			time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
		}}
	f2 := &Folder{
		Path:  "/home/foo",
		Count: 1,
		Times: Times{
			time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
		}}
	if f1.Compare(s, f2) != f1 {
		t.Fail()
	}
}

func TestCompare_directMatch(t *testing.T) {
	s := "foo"
	f1 := &Folder{
		Path:  "/home/foo",
		Count: 1,
		Times: Times{
			time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
		}}
	f2 := &Folder{
		Path:  "/home/nfoo",
		Count: 1,
		Times: Times{
			time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
		}}
	if f1.Compare(s, f2) != f1 {
		t.Fail()
	}
}

func TestCheckBaseSimilarity_quals(t *testing.T) {
	base := "foo"
	s := "foo"
	n := checkBaseSimilarity(base, s)
	if n != StrEquals {
		t.Fail()
	}
	//fmt.Println(n)
}

func TestCheckBaseSimilarity_wrongCase(t *testing.T) {
	base := "Foo"
	s := "foo"
	n := checkBaseSimilarity(base, s)
	if n != StrEqualsWrongCase {
		t.Fail()
	}
	//fmt.Println(n)
}

func TestCheckBaseSimilarity_noMatch(t *testing.T) {
	base := "foo"
	s := "Bar"
	n := checkBaseSimilarity(base, s)
	if n != NoMatch {
		t.Fail()
	}
	//fmt.Println(n)
}

func TestCheckBaseSimilarity_startsWith(t *testing.T) {
	base := "foobar"
	s := "foo"
	n := checkBaseSimilarity(base, s)
	if n != StrStartsEndsWith {
		t.Fail()
	}
	//fmt.Println(n)
}

func TestCheckBaseSimilarity_endsWith(t *testing.T) {
	base := "superfoo"
	s := "foo"
	n := checkBaseSimilarity(base, s)
	if n != StrStartsEndsWith {
		t.Fail()
	}
	//fmt.Println(n)
}

func TestCheckBaseSimilarity_contains(t *testing.T) {
	base := "nfooD"
	s := "foo"
	n := checkBaseSimilarity(base, s)
	if n != StrContains {
		t.Fail()
	}
	//fmt.Println(n)
}

func TestCheckBaseSimilarity_similar(t *testing.T) {
	base := "Bar"
	s := "bao"
	n := checkBaseSimilarity(base, s)
	if n != StrSimilar {
		t.Fail()
	}
	//fmt.Println(n)
}
