package main

import (
	"testing"
	"time"
)

// func TestRate_noMatch(t *testing.T) {
// 	n := rate("aaa", "/home/foo", []time.Time{time.Now()})
// 	if n != noMatch {
// 		t.Fail()
// 	}
// }

func TestRateSimilarity(t *testing.T) {
	// verbose = true
	tt := []struct {
		name, base, query string
		exp               uint
	}{
		{name: "equals", base: "foo", query: "foo", exp: strEquals},
		{name: "wrong case", base: "Foo", query: "foo", exp: strEqualsWrongCase},

		{name: "no match 1", base: "foo", query: "Bar", exp: noMatch},
		{name: "no match 2", base: "pubip", query: "book", exp: noMatch},
		{name: "no match 3", base: "tmp", query: "timer", exp: noMatch},
		{name: "no match 4", base: "tm", query: "tmp", exp: noMatch},

		{name: "starts with", base: "foobar", query: "foo", exp: strStartsWith},
		{name: "ends with", base: "superfoo", query: "foo", exp: strEndsWith},

		{name: "contains", base: "nfooD", query: "foo", exp: strContains},

		{name: "similar 1", base: "Bar", query: "bao", exp: strSimilar},
		{name: "similar 2", base: "Bar", query: "bart", exp: strSimilar},
		{name: "similar 3", base: "HubertVomSchuh", query: "Hub3rtV@mSchu", exp: strSimilar},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			n := rateSimilarity(tc.base, tc.query)
			if n != tc.exp {
				t.Errorf("%s - exp %v, got %v", tc.name, tc.exp, n)
			}
		})
	}
}

func TestRateTime(t *testing.T) {
	tt := []struct {
		name               string
		nowTime, otherTime time.Time
		exp                uint
	}{
		{name: "less than a minute", exp: timeLessThanMinute,
			nowTime:   time.Date(2009, time.November, 10, 23, 30, 10, 0, time.UTC),
			otherTime: time.Date(2009, time.November, 10, 23, 30, 0, 0, time.UTC)},

		{name: "less than five minutes", exp: timeLessThanFiveMinutes,
			nowTime:   time.Date(2009, time.November, 10, 23, 31, 10, 0, time.UTC),
			otherTime: time.Date(2009, time.November, 10, 23, 30, 0, 0, time.UTC)},

		{name: "less than an hour", exp: timeLessThanHour,
			nowTime:   time.Date(2009, time.November, 10, 23, 0, 10, 0, time.UTC),
			otherTime: time.Date(2009, time.November, 10, 22, 30, 0, 0, time.UTC)},

		{name: "less than six hours", exp: timeLessThanSixHours,
			nowTime:   time.Date(2009, time.November, 10, 23, 0, 10, 0, time.UTC),
			otherTime: time.Date(2009, time.November, 10, 20, 30, 0, 0, time.UTC)},

		{name: "less than twelve hours", exp: timeLessThanTwelveHours,
			nowTime:   time.Date(2009, time.November, 10, 12, 0, 10, 0, time.UTC),
			otherTime: time.Date(2009, time.November, 10, 01, 30, 0, 0, time.UTC)},

		{name: "less than a day", exp: timeLessThanDay,
			nowTime:   time.Date(2009, time.November, 10, 23, 0, 10, 0, time.UTC),
			otherTime: time.Date(2009, time.November, 10, 01, 30, 0, 0, time.UTC)},

		{name: "less than two days", exp: timeLessThanTwoDays,
			nowTime:   time.Date(2009, time.November, 11, 23, 0, 10, 0, time.UTC),
			otherTime: time.Date(2009, time.November, 10, 01, 30, 0, 0, time.UTC)},

		{name: "less than a week", exp: timeLessThanWeek,
			nowTime:   time.Date(2009, time.November, 12, 23, 0, 10, 0, time.UTC),
			otherTime: time.Date(2009, time.November, 10, 01, 30, 0, 0, time.UTC)},

		{name: "less than two weeks", exp: timeLessThanTwoWeeks,
			nowTime:   time.Date(2009, time.November, 12, 0, 0, 0, 0, time.UTC),
			otherTime: time.Date(2009, time.November, 01, 0, 0, 0, 0, time.UTC)},

		{name: "less than a month", exp: timeLessThanMonth,
			nowTime:   time.Date(2009, time.December, 1, 0, 0, 0, 0, time.UTC),
			otherTime: time.Date(2009, time.November, 12, 0, 0, 0, 0, time.UTC)},

		{name: "less than two months", exp: timeLessThanTwoMonths,
			nowTime:   time.Date(2009, time.March, 1, 0, 0, 0, 0, time.UTC),
			otherTime: time.Date(2009, time.January, 20, 0, 0, 0, 0, time.UTC)},

		{name: "less than six months", exp: timeLessThanSixMonths,
			nowTime:   time.Date(2009, time.May, 1, 0, 0, 0, 0, time.UTC),
			otherTime: time.Date(2009, time.January, 20, 0, 0, 0, 0, time.UTC)},

		{name: "less than a year", exp: timeLessThanYear,
			nowTime:   time.Date(2010, time.January, 1, 0, 0, 0, 0, time.UTC),
			otherTime: time.Date(2009, time.February, 1, 0, 0, 0, 0, time.UTC)},

		{name: "older than a year", exp: timeOlderThanAYear,
			nowTime:   time.Date(2012, time.January, 1, 0, 0, 0, 0, time.UTC),
			otherTime: time.Date(2009, time.February, 1, 0, 0, 0, 0, time.UTC)},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if n := rateTime(tc.nowTime, tc.otherTime); n != tc.exp {
				t.Errorf("%s - exp %v, got %v", tc.name, tc.exp, n)
			}
		})
	}
}

// func TestRate_maxRating(t *testing.T) {
// 	s := "foo"
// 	f := NewFolder("/home/foo", time.Now().Add(-time.Second*40))
// 	n := rate(s, f.Path, f.Times)
// 	if n != strEquals+timeLessThanMinute {
// 		t.Fail()
// 	}
// }

// type Dict []string

// func loadDict() (Dict, error) {
// 	f, err := os.Open("/usr/share/dict/american")
// 	if err != nil {
// 		return nil, err
// 	}
// 	buf, err := ioutil.ReadAll(f)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return strings.Split(string(buf), "\n"), nil
// }

// func (d Dict) randomPath() string {
// 	a := []string{"/"}
// 	for i := 0; i < random(2, 5); i++ {
// 		a = append(a, d.randomWord())
// 	}
// 	return filepath.Join(a...)
// }

// func (d Dict) randomWord() string {
// 	len := len(d)
// 	word := d[random(0, len-1)]
// 	word = strings.Replace(word, "'", "", -1)
// 	word = strings.Replace(word, " ", "_", -1)
// 	return word
// }

// // random with min inclusive and max exclusive
// func random(min, max int) int {
// 	return rand.Intn(max-min) + min
// }
