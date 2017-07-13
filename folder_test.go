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

// func TestRate_foo(t *testing.T) {
// 	s := "foo"
// 	now := time.Now()
// 	t1 := now.Add(-time.Second * 40)
// 	t2 := now.Add(-time.Hour * 18)
// 	t3 := now.Add(-time.Hour * 24)
// 	t4 := now.Add(-time.Hour * 24 * 2)
// 	t5 := now.Add(-time.Hour * 24 * 7 * 2)
// 	f := &Folder{
// 		Path:  "/home/foo",
// 		Count: 1,
// 		Times: Times{t1, t2, t3, t4, t5},
// 	}
// 	n := rate(s, f.Path, f.Times)
// 	fmt.Println(n)
// }

// func TestBar(t *testing.T) {
// 	s := "foo"
// 	r := nowTimeRepo()
// 	var a RatedFolders
// 	for _, f := range r.All() {
// 		rf := NewRatedFolder(f, s)
// 		if rf.Points == NoMatch {
// 			continue
// 		}
// 		//fmt.Println("append:", rf.Folder, " points:", rf.Points)
// 		a = append(a, rf)
// 	}
// 	fmt.Println(a)
// 	fmt.Println("-----------------------------------------------")
// 	sort.Sort(a)
// 	fmt.Println(a)
// }

func dummy() RatedTimeFolders {
	fn := func(p string, t time.Time) RatedFolder {
		return NewRatedFolder(NewFolder(p, 1, Times{t}), "")
	}
	now := time.Now()
	f1 := fn("/home/bar", now.Add(-time.Hour*18))
	f2 := fn("/home/zot", now.Add(-time.Hour*4))
	f3 := fn("/home/foo", now)
	return RatedTimeFolders{f1, f2, f3}
}

func TestTimeRatedSort(t *testing.T) {
	a := dummy()
	a.sort()
	if a[0].folder.Path != "/home/foo" {
		t.Fail()
	}
}

func TestRemoveOldestFolders(t *testing.T) {
	a := dummy()
	m := a.removeOldestFolders(2)
	if len(m) != 2 {
		t.Fail()
	}
	if _, ok := m["/home/bar"]; ok {
		t.Fail()
	}
}
