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
	timeOlderThanAYear      = 3

	strEquals          = 100
	strEqualsWrongCase = 80
	strStartsEndsWith  = 60
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

// TODO write startWith endWith checks
// if len(s) is combined in word -> strContains
func rateSimilarity(base, s string) uint {
	// equals
	if base == s {
		return strEquals
	}
	base = strings.ToLower(base)
	s = strings.ToLower(s)
	// equals wrong case
	if base == s {
		return strEqualsWrongCase
	}
	// starts or ends with
	if strings.HasPrefix(base, s) || strings.HasSuffix(base, s) {
		return strStartsEndsWith
	}
	// does base even contain s?
	if strings.Contains(base, s) {
		// TODO check how much different s is compared to base
		return strContains
	}
	// search for similarities
	baseLen := utf8.RuneCountInString(base)
	// no similar comarisons on short words
	if baseLen < 3 {
		return noMatch
	}
	var mindiff int
	if baseLen <= 4 {
		mindiff = 1
	} else if baseLen <= 10 {
		mindiff = 2
	} else {
		mindiff = 3
	}
	// find differences
	var diff int
	for _, r := range base {
		// check if rune in searched string
		if ok := strings.ContainsRune(s, r); !ok {
			diff++
		}
	}
	if diff <= mindiff {
		return strSimilar
	}
	return noMatch
}