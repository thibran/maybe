package main

import (
	"fmt"
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
	timeOlderThanAYear      = 3

	strEquals          = 160
	strEqualsWrongCase = 80
	strStartsWith      = 70
	strEndsWith        = 60
	strContains        = 40
	strSimilar         = 20
	noMatch            = 0
)

// rate search-term s for path p with time-slice a.
func rate(s, p string, a Times) uint {
	base := path.Base(p)
	var n uint
	n += rateSimilarity(base, s)
	// no string match -> return
	if n == noMatch {
		return n
	}
	timeRate := ratePassedTime(a)
	n += timeRate
	return n
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
	l := logWithPrefix("rateSimilarity")
	l(fmt.Sprintf("base: %s, search for: %q", base, query))
	// l("base:", base, "search for:", s)
	if base == query {
		l("strEquals:", strEquals)
		return strEquals
	}
	base = strings.ToLower(base)
	query = strings.ToLower(query)
	// equals wrong case
	if base == query {
		l("strEqualsWrongCase:", strEqualsWrongCase)
		return strEqualsWrongCase
	}
	// starts with
	if strings.HasPrefix(base, query) {
		l("strStartsWith:", strStartsWith)
		return strStartsWith
	}
	// ends with
	if strings.HasSuffix(base, query) {
		l("strEndsWith:", strEndsWith)
		return strEndsWith
	}
	// does base even contain s?
	if strings.Contains(base, query) {
		l("strContains:", strContains)
		return strContains
	}
	return strSimilarity(base, query)
}

func strSimilarity(base, query string) uint {
	baseLen := utf8.RuneCountInString(base)
	// don't compare too short words
	l := logWithPrefix("strSimilarity")
	if baseLen < 3 {
		l("baseLen < 3:", noMatch)
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
		l("diff <= maxdiff:", strSimilar)
		return strSimilar
	}
	l("end:", "noMatch:", noMatch)
	return noMatch
}
