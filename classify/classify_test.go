package classify

import (
	"testing"
	"time"
)

func TestClassifyText(t *testing.T) {
	// pref.Verbose = true
	tt := []struct {
		name, base, query string
		exp               uint
		insensitive       bool
	}{
		{name: "equals 1", base: "foo", query: "foo", exp: StrEquals},
		{name: "equals 2", base: ".foo", query: ".foo", exp: StrEquals},
		{name: "equals 3", base: ".foo", query: "foo", exp: StrEquals},

		{name: "wrong case 1", base: "Foo", query: "foo",
			exp: StrEqualsWrongCase},
		{name: "wrong case 2", base: "Sync", query: "Sync",
			exp: StrEqualsWrongCase, insensitive: true},

		{name: "no match 1", base: "foo", query: "Bar", exp: NoMatch},
		{name: "no match 2", base: "pubip", query: "book", exp: NoMatch},
		{name: "no match 3", base: "tmp", query: "timer", exp: NoMatch},
		{name: "no match 4", base: "tm", query: "tmp", exp: NoMatch},

		{name: "starts with", base: "foobar", query: "foo", exp: StrStartsWith},
		{name: "ends with", base: "superfoo", query: "foo", exp: StrEndsWith},

		{name: "contains", base: "nfooD", query: "foo", exp: StrContains},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			n := classifyText(tc.base, tc.query, !tc.insensitive)
			if n != tc.exp {
				t.Errorf("%s - exp %v, got %v", tc.name, tc.exp, n)
			}
		})
	}
}

func TestSimilarity(t *testing.T) {
	// pref.Verbose = true
	tt := []struct {
		name, base, query string
		exp               uint
	}{
		{name: "similar 1",
			base: "bar", query: "bao", exp: StrSimilar},
		{name: "similar 2",
			base: "bar", query: "bart", exp: StrSimilar},
		{name: "similar 3",
			base: "hubertvomschuh", query: "hub3rtv@mschu", exp: StrSimilar},
		{name: "similar 4",
			base: "foo", query: ".foo", exp: StrSimilar},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			n := similarity(tc.base, tc.query)
			if n != tc.exp {
				t.Errorf("%s - exp %v, got %v", tc.name, tc.exp, n)
			}
		})
	}
}

func TestTimeHelper(t *testing.T) {
	// pref.Verbose = true
	tt := []struct {
		name               string
		nowTime, otherTime time.Time
		exp                uint
	}{
		{name: "less than a minute", exp: TimeLessThanMinute,
			nowTime:   time.Date(2009, time.November, 10, 23, 30, 10, 0, time.UTC),
			otherTime: time.Date(2009, time.November, 10, 23, 30, 0, 0, time.UTC)},

		{name: "less than five minutes", exp: TimeLessThanFiveMinutes,
			nowTime:   time.Date(2009, time.November, 10, 23, 31, 10, 0, time.UTC),
			otherTime: time.Date(2009, time.November, 10, 23, 30, 0, 0, time.UTC)},

		{name: "less than an hour", exp: TimeLessThanHour,
			nowTime:   time.Date(2009, time.November, 10, 23, 0, 10, 0, time.UTC),
			otherTime: time.Date(2009, time.November, 10, 22, 30, 0, 0, time.UTC)},

		{name: "less than six hours", exp: TimeLessThanSixHours,
			nowTime:   time.Date(2009, time.November, 10, 23, 0, 10, 0, time.UTC),
			otherTime: time.Date(2009, time.November, 10, 20, 30, 0, 0, time.UTC)},

		{name: "less than twelve hours", exp: TimeLessThanTwelveHours,
			nowTime:   time.Date(2009, time.November, 10, 12, 0, 10, 0, time.UTC),
			otherTime: time.Date(2009, time.November, 10, 01, 30, 0, 0, time.UTC)},

		{name: "less than a day", exp: TimeLessThanDay,
			nowTime:   time.Date(2009, time.November, 10, 23, 0, 10, 0, time.UTC),
			otherTime: time.Date(2009, time.November, 10, 01, 30, 0, 0, time.UTC)},

		{name: "less than two days", exp: TimeLessThanTwoDays,
			nowTime:   time.Date(2009, time.November, 11, 23, 0, 10, 0, time.UTC),
			otherTime: time.Date(2009, time.November, 10, 01, 30, 0, 0, time.UTC)},

		{name: "less than a week", exp: TimeLessThanWeek,
			nowTime:   time.Date(2009, time.November, 12, 23, 0, 10, 0, time.UTC),
			otherTime: time.Date(2009, time.November, 10, 01, 30, 0, 0, time.UTC)},

		{name: "less than two weeks", exp: TimeLessThanTwoWeeks,
			nowTime:   time.Date(2009, time.November, 12, 0, 0, 0, 0, time.UTC),
			otherTime: time.Date(2009, time.November, 01, 0, 0, 0, 0, time.UTC)},

		{name: "less than a month", exp: TimeLessThanMonth,
			nowTime:   time.Date(2009, time.December, 1, 0, 0, 0, 0, time.UTC),
			otherTime: time.Date(2009, time.November, 12, 0, 0, 0, 0, time.UTC)},

		{name: "less than two months", exp: TimeLessThanTwoMonths,
			nowTime:   time.Date(2009, time.March, 1, 0, 0, 0, 0, time.UTC),
			otherTime: time.Date(2009, time.January, 20, 0, 0, 0, 0, time.UTC)},

		{name: "less than six months", exp: TimeLessThanSixMonths,
			nowTime:   time.Date(2009, time.May, 1, 0, 0, 0, 0, time.UTC),
			otherTime: time.Date(2009, time.January, 20, 0, 0, 0, 0, time.UTC)},

		{name: "less than a year", exp: TimeLessThanYear,
			nowTime:   time.Date(2010, time.January, 1, 0, 0, 0, 0, time.UTC),
			otherTime: time.Date(2009, time.February, 1, 0, 0, 0, 0, time.UTC)},

		{name: "older than a year", exp: TimeOlderThanAYear,
			nowTime:   time.Date(2012, time.January, 1, 0, 0, 0, 0, time.UTC),
			otherTime: time.Date(2009, time.February, 1, 0, 0, 0, 0, time.UTC)},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if n := timeHelper(tc.nowTime, tc.otherTime); n != tc.exp {
				t.Errorf("%s - exp %v, got %v", tc.name, tc.exp, n)
			}
		})
	}
}

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
