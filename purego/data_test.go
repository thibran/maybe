package main

import (
	"fmt"
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

func TestRateTime_lessThanMinute(t *testing.T) {
	now := time.Date(2009, time.November, 10, 23, 30, 10, 0, time.UTC)
	t1 := time.Date(2009, time.November, 10, 23, 30, 0, 0, time.UTC)
	n := rateTime(now, t1)
	if n != TimeLessThanMinute {
		t.Fail()
	}
	//fmt.Println(n)
}

func TestRateTime_lessThanFiveMinutes(t *testing.T) {
	now := time.Date(2009, time.November, 10, 23, 31, 10, 0, time.UTC)
	t1 := time.Date(2009, time.November, 10, 23, 30, 0, 0, time.UTC)
	n := rateTime(now, t1)
	if n != TimeLessThanFiveMinutes {
		t.Fail()
	}
	//fmt.Println(n)
}

func TestRateTime_lessThanHour(t *testing.T) {
	now := time.Date(2009, time.November, 10, 23, 0, 10, 0, time.UTC)
	t1 := time.Date(2009, time.November, 10, 22, 30, 0, 0, time.UTC)
	n := rateTime(now, t1)
	if n != TimeLessThanHour {
		t.Fail()
	}
	//fmt.Println(n)
}

func TestRateTime_lessThanSixHours(t *testing.T) {
	now := time.Date(2009, time.November, 10, 23, 0, 10, 0, time.UTC)
	t1 := time.Date(2009, time.November, 10, 20, 30, 0, 0, time.UTC)
	n := rateTime(now, t1)
	if n != TimeLessThanSixHours {
		t.Fail()
	}
	//fmt.Println(n)
}

func TestRateTime_lessThanTwelveHours(t *testing.T) {
	now := time.Date(2009, time.November, 10, 12, 0, 10, 0, time.UTC)
	t1 := time.Date(2009, time.November, 10, 01, 30, 0, 0, time.UTC)
	n := rateTime(now, t1)
	if n != TimeLessThanTwelveHours {
		t.Fail()
	}
	//fmt.Println(n)
}

func TestRateTime_lessThanDay(t *testing.T) {
	now := time.Date(2009, time.November, 10, 23, 0, 10, 0, time.UTC)
	t1 := time.Date(2009, time.November, 10, 01, 30, 0, 0, time.UTC)
	n := rateTime(now, t1)
	if n != TimeLessThanDay {
		t.Fail()
	}
	//fmt.Println(n)
}

func TestRateTime_lessThanTwoDays(t *testing.T) {
	now := time.Date(2009, time.November, 11, 23, 0, 10, 0, time.UTC)
	t1 := time.Date(2009, time.November, 10, 01, 30, 0, 0, time.UTC)
	n := rateTime(now, t1)
	if n != TimeLessThanTwoDays {
		t.Fail()
	}
	//fmt.Println(n)
}

func TestRateTime_lessThanWeek(t *testing.T) {
	now := time.Date(2009, time.November, 12, 23, 0, 10, 0, time.UTC)
	t1 := time.Date(2009, time.November, 10, 01, 30, 0, 0, time.UTC)
	n := rateTime(now, t1)
	if n != TimeLessThanWeek {
		t.Fail()
	}
	//fmt.Println(n)
}

func TestRateTime_lessThanTwoWeeks(t *testing.T) {
	now := time.Date(2009, time.November, 12, 0, 0, 0, 0, time.UTC)
	t1 := time.Date(2009, time.November, 01, 0, 0, 0, 0, time.UTC)
	n := rateTime(now, t1)
	if n != TimeLessThanTwoWeeks {
		t.Fail()
	}
	//fmt.Println(n)
}

func TestRateTime_lessThanMonth(t *testing.T) {
	now := time.Date(2009, time.December, 1, 0, 0, 0, 0, time.UTC)
	t1 := time.Date(2009, time.November, 12, 0, 0, 0, 0, time.UTC)
	n := rateTime(now, t1)
	if n != TimeLessThanMonth {
		t.Fail()
	}
	//fmt.Println(n)
}

func TestRateTime_lessThanTwoMonths(t *testing.T) {
	now := time.Date(2009, time.March, 1, 0, 0, 0, 0, time.UTC)
	t1 := time.Date(2009, time.January, 20, 0, 0, 0, 0, time.UTC)
	n := rateTime(now, t1)
	if n != TimeLessThanTwoMonths {
		t.Fail()
	}
	//fmt.Println(n)
}

func TestRateTime_lessThanSixMonths(t *testing.T) {
	now := time.Date(2009, time.May, 1, 0, 0, 0, 0, time.UTC)
	t1 := time.Date(2009, time.January, 20, 0, 0, 0, 0, time.UTC)
	n := rateTime(now, t1)
	if n != TimeLessThanSixMonths {
		t.Fail()
	}
	//fmt.Println(n)
}

func TestRateTime_lessThanYear(t *testing.T) {
	now := time.Date(2010, time.January, 1, 0, 0, 0, 0, time.UTC)
	t1 := time.Date(2009, time.February, 1, 0, 0, 0, 0, time.UTC)
	n := rateTime(now, t1)
	if n != TimeLessThanYear {
		t.Fail()
	}
	//fmt.Println(n)
}

func TestRateTime_olderThanAYea(t *testing.T) {
	now := time.Date(2012, time.January, 1, 0, 0, 0, 0, time.UTC)
	t1 := time.Date(2009, time.February, 1, 0, 0, 0, 0, time.UTC)
	n := rateTime(now, t1)
	if n != TimeOlderThanAYear {
		t.Fail()
	}
	//fmt.Println(n)
}

func TestRate_maxRating(t *testing.T) {
	s := "foo"
	now := time.Now().Add(-time.Second * 40)
	f := &Folder{
		Path:  "/home/foo",
		Count: 1,
		Times: Times{now},
	}
	n := f.rate(s)
	if n != StrEquals+TimeLessThanMinute {
		t.Fail()
	}
	fmt.Println(n)
}

func TestRate_foo(t *testing.T) {
	s := "foo"
	now := time.Now()
	t1 := now.Add(-time.Second * 40)
	t2 := now.Add(-time.Hour * 18)
	t3 := now.Add(-time.Hour * 24)
	t4 := now.Add(-time.Hour * 24 * 2)
	t5 := now.Add(-time.Hour * 24 * 7 * 2)
	f := &Folder{
		Path:  "/home/foo",
		Count: 1,
		Times: Times{t1, t2, t3, t4, t5},
	}
	n := f.rate(s)
	// if n != StrEquals+TimeLessThanMinute {
	// 	t.Fail()
	// }
	fmt.Println(n)
}
