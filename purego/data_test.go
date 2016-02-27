package main

import "testing"

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
