package main

import (
	"fmt"
	"path"
	"sort"
	"strings"
	"time"
	"unicode/utf8"
)

// maximum number of time.Time entries in a Times slice.
const MaxTimesEntries = 10

// Times is a shorthand for a time slice.
type Times []time.Time

func (t Times) Len() int {
	return len(t)
}

func (t Times) Less(i, j int) bool {
	return t[i].After(t[j])
}

func (t Times) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

// sort time entries and cut all entries longer than MaxTimesEntries.
func (t Times) sort() Times {
	sort.Sort(t)
	if len(t) > MaxTimesEntries {
		return t[:MaxTimesEntries]
	}
	return t
}

// Folder entry.
type Folder struct {
	Path  string
	Count int   // counts how often the folder has been updated
	Times Times // last MaxTimesEntries updates
}

// RatedFolder object.
type RatedFolder struct {
	Points int
	Folder Folder
}

// RatedFolders is an alias for RatedFolder-slice.
type RatedFolders []RatedFolder

// NewRatedFolder creates a new object.
func NewRatedFolder(f Folder, s string) RatedFolder {
	return RatedFolder{
		Points: rate(s, f.Path, f.Times),
		Folder: f,
	}
}

func (a RatedFolders) String() string {
	arr := make([]string, len(a))
	for k, v := range a {
		arr[k] = fmt.Sprintf("%v  points: %v", v.Folder, v.Points)
	}
	return strings.Join(arr, "\n")
}

func (a RatedFolders) Len() int {
	return len(a)
}

func (a RatedFolders) Less(i, j int) bool {
	if a[i].Points == a[j].Points {
		return a[i].Folder.Count > a[j].Folder.Count
	}
	return a[i].Points > a[j].Points
}

func (a RatedFolders) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

const (
	TimeLessThanMinute      = 42
	TimeLessThanFiveMinutes = 39
	TimeLessThanHour        = 36
	TimeLessThanSixHours    = 33
	TimeLessThanTwelveHours = 30
	TimeLessThanDay         = 27
	TimeLessThanTwoDays     = 24
	TimeLessThanWeek        = 21
	TimeLessThanTwoWeeks    = 18
	TimeLessThanMonth       = 15
	TimeLessThanTwoMonths   = 12
	TimeLessThanSixMonths   = 9
	TimeLessThanYear        = 6
	TimeOlderThanAYear      = 3

	StrEquals          = 100
	StrEqualsWrongCase = 80
	StrStartsEndsWith  = 60
	StrContains        = 40
	StrSimilar         = 20
	NoMatch            = 0
)

// rate search-term s for path p with time-slice a.
func rate(s, p string, a Times) int {
	base := path.Base(p)
	var n int
	n += rateKeySimilarity(base, s)
	//fmt.Println("s:", n)
	// no string match -> return
	if n == NoMatch {
		return n
	}
	timeRate := ratePassedTime(a)
	n += timeRate
	//fmt.Println("t:", timeRate)
	return n
}

// TODO rate time passed since
func ratePassedTime(a Times) int {
	var n int
	now := time.Now()
	for _, t := range a {
		n += rateTime(now, t)
	}
	return n
}

func rateTime(now, t time.Time) int {
	// < minute
	if now.Before(t.Add(time.Minute)) {
		return TimeLessThanMinute
	}
	// < 5 min
	if now.Before(t.Add(time.Minute * 5)) {
		return TimeLessThanFiveMinutes
	}
	// < hour
	if now.Before(t.Add(time.Hour)) {
		return TimeLessThanHour
	}
	// < 6 hours
	if now.Before(t.Add(time.Hour * 6)) {
		return TimeLessThanSixHours
	}
	// < 12 hours
	if now.Before(t.Add(time.Hour * 12)) {
		return TimeLessThanTwelveHours
	}
	// < day
	if now.Before(t.Add(time.Hour * 24)) {
		return TimeLessThanDay
	}
	// < 2 days
	if now.Before(t.Add(time.Hour * 48)) {
		return TimeLessThanTwoDays
	}
	// < week
	if now.Before(t.Add(time.Hour * 24 * 7)) {
		return TimeLessThanWeek
	}
	// < 2 weeks
	if now.Before(t.Add(time.Hour * 24 * 7 * 2)) {
		return TimeLessThanTwoWeeks
	}
	// < month
	if now.Before(t.Add(time.Hour * 24 * 7 * 4)) {
		return TimeLessThanMonth
	}
	// < 2 months
	if now.Before(t.Add(time.Hour * 24 * 7 * 4 * 2)) {
		return TimeLessThanTwoMonths
	}
	// < 6 months
	if now.Before(t.Add(time.Hour * 24 * 7 * 4 * 6)) {
		return TimeLessThanSixMonths
	}
	// < year
	if now.Before(t.Add(time.Hour * 24 * 7 * 4 * 12)) {
		return TimeLessThanYear
	}
	return TimeOlderThanAYear
}

// TODO write startWith endWith checks
// if len(s) is combined in word -> StrContains
func rateKeySimilarity(base, s string) int {
	// equals
	if base == s {
		return StrEquals
	}
	base = strings.ToLower(base)
	s = strings.ToLower(s)
	// equals wrong case
	if base == s {
		return StrEqualsWrongCase
	}
	// starts or ends with
	if strings.HasPrefix(base, s) || strings.HasSuffix(base, s) {
		return StrStartsEndsWith
	}
	// does base even contain s?
	if strings.Contains(base, s) {
		// TODO check how much different s is compared to base
		return StrContains
	}
	// search for similarities
	baseLen := utf8.RuneCountInString(base)
	// no similar comarisons on short words
	if baseLen < 3 {
		return NoMatch
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
		return StrSimilar
	}
	return NoMatch
}
