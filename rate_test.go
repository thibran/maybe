package main

import (
	"testing"
	"time"
)

func testSimilarity(t *testing.T, base, s string, exp uint) uint {
	n := rateSimilarity(base, s)
	if n != exp {
		t.Fail()
	}
	return n
}

func TestRateKeySimilarity_quals(t *testing.T) {
	testSimilarity(t, "foo", "foo", strEquals)
}

func TestRateKeySimilarity_wrongCase(t *testing.T) {
	testSimilarity(t, "Foo", "foo", strEqualsWrongCase)
}

func TestRateKeySimilarity_noMatch(t *testing.T) {
	testSimilarity(t, "foo", "Bar", noMatch)
}

func TestRateKeySimilarity_startsWith(t *testing.T) {
	testSimilarity(t, "foobar", "foo", strStartsEndsWith)
}

func TestRateKeySimilarity_endsWith(t *testing.T) {
	testSimilarity(t, "superfoo", "foo", strStartsEndsWith)
}

func TestRateKeySimilarity_contains(t *testing.T) {
	testSimilarity(t, "nfooD", "foo", strContains)
}

func TestRateKeySimilarity_similar(t *testing.T) {
	testSimilarity(t, "Bar", "bao", strSimilar)
}

func testTime(t *testing.T, now, t1 time.Time, exp uint) uint {
	n := rateTime(now, t1)
	if n != exp {
		t.Fail()
	}
	return n
}

func TestRateTime_lessThanMinute(t *testing.T) {
	testTime(t,
		time.Date(2009, time.November, 10, 23, 30, 10, 0, time.UTC),
		time.Date(2009, time.November, 10, 23, 30, 0, 0, time.UTC),
		timeLessThanMinute,
	)
}

func TestRateTime_lessThanFiveMinutes(t *testing.T) {
	testTime(t,
		time.Date(2009, time.November, 10, 23, 31, 10, 0, time.UTC),
		time.Date(2009, time.November, 10, 23, 30, 0, 0, time.UTC),
		timeLessThanFiveMinutes,
	)
}

func TestRateTime_lessThanHour(t *testing.T) {
	testTime(t,
		time.Date(2009, time.November, 10, 23, 0, 10, 0, time.UTC),
		time.Date(2009, time.November, 10, 22, 30, 0, 0, time.UTC),
		timeLessThanHour,
	)
}

func TestRateTime_lessThanSixHours(t *testing.T) {
	testTime(t,
		time.Date(2009, time.November, 10, 23, 0, 10, 0, time.UTC),
		time.Date(2009, time.November, 10, 20, 30, 0, 0, time.UTC),
		timeLessThanSixHours,
	)
}

func TestRateTime_lessThanTwelveHours(t *testing.T) {
	testTime(t,
		time.Date(2009, time.November, 10, 12, 0, 10, 0, time.UTC),
		time.Date(2009, time.November, 10, 01, 30, 0, 0, time.UTC),
		timeLessThanTwelveHours,
	)
}

func TestRateTime_lessThanDay(t *testing.T) {
	testTime(t,
		time.Date(2009, time.November, 10, 23, 0, 10, 0, time.UTC),
		time.Date(2009, time.November, 10, 01, 30, 0, 0, time.UTC),
		timeLessThanDay,
	)
}

func TestRateTime_lessThanTwoDays(t *testing.T) {
	testTime(t,
		time.Date(2009, time.November, 11, 23, 0, 10, 0, time.UTC),
		time.Date(2009, time.November, 10, 01, 30, 0, 0, time.UTC),
		timeLessThanTwoDays,
	)
}

func TestRateTime_lessThanWeek(t *testing.T) {
	testTime(t,
		time.Date(2009, time.November, 12, 23, 0, 10, 0, time.UTC),
		time.Date(2009, time.November, 10, 01, 30, 0, 0, time.UTC),
		timeLessThanWeek,
	)
}

func TestRateTime_lessThanTwoWeeks(t *testing.T) {
	testTime(t,
		time.Date(2009, time.November, 12, 0, 0, 0, 0, time.UTC),
		time.Date(2009, time.November, 01, 0, 0, 0, 0, time.UTC),
		timeLessThanTwoWeeks,
	)
}

func TestRateTime_lessThanMonth(t *testing.T) {
	testTime(t,
		time.Date(2009, time.December, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2009, time.November, 12, 0, 0, 0, 0, time.UTC),
		timeLessThanMonth,
	)
}

func TestRateTime_lessThanTwoMonths(t *testing.T) {
	testTime(t,
		time.Date(2009, time.March, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2009, time.January, 20, 0, 0, 0, 0, time.UTC),
		timeLessThanTwoMonths,
	)
}

func TestRateTime_lessThanSixMonths(t *testing.T) {
	testTime(t,
		time.Date(2009, time.May, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2009, time.January, 20, 0, 0, 0, 0, time.UTC),
		timeLessThanSixMonths,
	)
}

func TestRateTime_lessThanYear(t *testing.T) {
	testTime(t,
		time.Date(2010, time.January, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2009, time.February, 1, 0, 0, 0, 0, time.UTC),
		timeLessThanYear,
	)
}

func TestRateTime_olderThanAYear(t *testing.T) {
	testTime(t,
		time.Date(2012, time.January, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2009, time.February, 1, 0, 0, 0, 0, time.UTC),
		timeOlderThanAYear,
	)
}

func TestRate_maxRating(t *testing.T) {
	s := "foo"
	now := time.Now().Add(-time.Second * 40)
	f := NewFolder(
		"/home/foo",
		1,
		Times{now},
	)
	n := rate(s, f.path, f.times)
	if n != strEquals+timeLessThanMinute {
		t.Fail()
	}
}
