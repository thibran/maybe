package main

import (
	"path"
	"strings"
	"time"
	"unicode/utf8"
)

const (
	timeLessThanMinute      = 42
	timeLessThanFiveMinutes = 39
	timeLessThanHour        = 36
	timeLessThanSixHours    = 33
	timeLessThanTwelveHours = 30
	timeLessThanDay         = 27
	timeLessThanTwoDays     = 24
	timeLessThanWeek        = 21
	timeLessThanTwoWeeks    = 18
	timeLessThanMonth       = 15
	timeLessThanTwoMonths   = 12
	timeLessThanSixMonths   = 9
	timeLessThanYear        = 6
	timeOlderThanAYear      = 0

	strEquals          = 200 // 160
	strEqualsWrongCase = 80
	strStartsWith      = 70
	strEndsWith        = 60
	strContains        = 40
	strSimilar         = 20
	noMatch            = 0
)

// Rating of a search query
type Rating struct {
	timePoints       uint
	similarityPoints uint
}

// points return the point sum of a rateing.
// If no similarity is found, time points are ignored.
func (r *Rating) points() uint {
	return r.similarityPoints + r.timePoints
}

// newRating rates search-term s for path p within time-slice a.
func newRating(s, p string, a Times) Rating {
	base := path.Base(p)
	n := rateSimilarity(base, s)
	if n == noMatch {
		return Rating{}
	}
	timeRate := ratePassedTime(a)
	return Rating{similarityPoints: n, timePoints: timeRate}
}

func ratePassedTime(a Times) uint {
	var n uint
	now := time.Now()
	for _, t := range a {
		n += rateTime(now, t)
	}
	return n
}

func rateTime(now, t time.Time) uint {
	beforeNow := func(n time.Duration) bool {
		return now.Before(t.Add(n))
	}
	// < minute
	if beforeNow(time.Minute) {
		return timeLessThanMinute
	}
	// < 5 min
	if beforeNow(time.Minute * 5) {
		return timeLessThanFiveMinutes
	}
	// < hour
	if beforeNow(time.Hour) {
		return timeLessThanHour
	}
	// < 6 hours
	if beforeNow(time.Hour * 6) {
		return timeLessThanSixHours
	}
	// < 12 hours
	if beforeNow(time.Hour * 12) {
		return timeLessThanTwelveHours
	}
	// < day
	if beforeNow(time.Hour * 24) {
		return timeLessThanDay
	}
	// < 2 days
	if beforeNow(time.Hour * 48) {
		return timeLessThanTwoDays
	}
	// < week
	if beforeNow(time.Hour * 24 * 7) {
		return timeLessThanWeek
	}
	// < 2 weeks
	if beforeNow(time.Hour * 24 * 7 * 2) {
		return timeLessThanTwoWeeks
	}
	// < month
	if beforeNow(time.Hour * 24 * 7 * 4) {
		return timeLessThanMonth
	}
	// < 2 months
	if beforeNow(time.Hour * 24 * 7 * 4 * 2) {
		return timeLessThanTwoMonths
	}
	// < 6 months
	if beforeNow(time.Hour * 24 * 7 * 4 * 6) {
		return timeLessThanSixMonths
	}
	// < year
	if beforeNow(time.Hour * 24 * 7 * 4 * 12) {
		return timeLessThanYear
	}
	return timeOlderThanAYear
}

// if len(s) is combined in word -> strContains
func rateSimilarity(base, query string) uint {
	if base == query {
		return strEquals
	}
	base = strings.ToLower(base)
	query = strings.ToLower(query)
	// equals wrong case
	if base == query {
		return strEqualsWrongCase
	}
	// starts with
	if strings.HasPrefix(base, query) {
		return strStartsWith
	}
	// ends with
	if strings.HasSuffix(base, query) {
		return strEndsWith
	}
	// does base even contain s?
	if strings.Contains(base, query) {
		return strContains
	}
	return strSimilarity(base, query)
}

func strSimilarity(base, query string) uint {
	baseLen := utf8.RuneCountInString(base)
	// don't compare too short words
	if baseLen < 3 {
		return noMatch
	}

	var maxdiff int
	if baseLen <= 4 {
		maxdiff = 1
	} else if baseLen <= 10 {
		maxdiff = 2
	} else {
		maxdiff = 3
	}

	// find differences, e.g.: foo & foa are similare
	var diff int
	runes := []rune(query)
	searchLen := len(runes)
	for k, v := range base {
		if k < searchLen && v == runes[k] {
			continue
		}
		diff++
	}

	if diff <= maxdiff {
		return strSimilar
	}
	return noMatch
}
