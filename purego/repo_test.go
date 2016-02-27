package main

import (
	"testing"
	"time"
)

func TestNewRepoDummy(t *testing.T) {
	r := NewRepoDummy()
	if r == nil {
		t.Fail()
	}
}

func TestAll(t *testing.T) {
	r := NewRepoDummy()
	if len(r.All()) != len(r.m) {
		t.Fail()
	}
}

func TestAdd_newObj(t *testing.T) {
	r := NewRepoDummy()
	oldLen := len(r.m)
	r.Add("/foo/bar", time.Now())
	if len(r.m) != oldLen+1 {
		t.Fail()
	}
}

func TestAdd_updateExisting(t *testing.T) {
	r := NewRepoDummy()
	timeNow := time.Now()
	r.Add(f1.Path, timeNow)
	f := r.m[f1.Path]
	if f.Count != 2 {
		t.Fail()
	}
	if f.Times[0] != timeNow {
		t.Error("Times[0] should be equals timeNow.")
	}
	if len(f.Times) > MaxTimesEntries {
		t.Fail()
	}
}

// func TestSortTimes(t *testing.T) {
// 	a := Times{f1.Times[0], f2.Times[0], f3.Times[0]}
// 	for _, v := range a {
// 		fmt.Println(v)
// 	}
// 	fmt.Println("---------------")
// 	sort.Sort(a)
//
// 	if len(a) > MAX_TIMES_VALUES {
// 		a = a[:MAX_TIMES_VALUES]
// 	}
// 	for _, v := range a {
// 		fmt.Println(v)
// 	}
// }
