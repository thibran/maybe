package main

import (
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

// Compare folder with another folder object.
func (f *Folder) Compare(s string, other *Folder) *Folder {
	if f.rate(s) >= other.rate(s) {
		return f
	}
	return other
}

const (
	StrEquals          = 100
	StrEqualsWrongCase = 80
	StrContains        = 40
	StrSimilar         = 20
	NoMatch            = 0
)

func (f *Folder) rate(s string) int {
	base := path.Base(f.Path)
	var n int
	n += checkBaseSimilarity(base, s)
	// if base == s {
	// 	n += StrEquals
	// }
	return n
}

// TODO rate time passed since
func ratePassedTime(t1 Times) int {

	return 0
}

// TODO write startWith endWith checks
// if len(s) is combined in word -> StrContains
func checkBaseSimilarity(base, s string) int {
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
