package classify

import (
	"fmt"
	"path"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/thibran/maybe/pref"
)

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
	TimeOlderThanAYear      = 0
	StrEquals               = 60
	StrEqualsWrongCase      = 50
	StrStartsWith           = 40
	StrEndsWith             = 30
	StrContains             = 20
	StrSimilar              = 10
	NoMatch                 = 0
)

// Rating of a search query
type Rating struct {
	TimePoints       uint
	SimilarityPoints uint
}

// Points return the point sum of a rateing.
// If no similarity is found, time points are ignored.
func (r *Rating) Points() uint {
	return r.SimilarityPoints + r.TimePoints
}

// NewRating rates search-term s for path p within time-slice a.
func NewRating(s, p string, a ...time.Time) (*Rating, error) {
	base := path.Base(p)
	n := classifyText(base, s, pref.CaseSensitive)
	if n == NoMatch {
		return nil, fmt.Errorf("NewRating - similarity: noMatch")
	}
	timeRate := classifyTime(a...)
	return &Rating{SimilarityPoints: n, TimePoints: timeRate}, nil
}

func classifyTime(a ...time.Time) uint {
	var n uint
	now := time.Now()
	for _, t := range a {
		n += timeHelper(now, t)
	}
	return n
}

func timeHelper(now, t time.Time) uint {
	beforeNow := func(n time.Duration) bool {
		return now.Before(t.Add(n))
	}
	// < minute
	if beforeNow(time.Minute) {
		return TimeLessThanMinute
	}
	// < 5 min
	if beforeNow(time.Minute * 5) {
		return TimeLessThanFiveMinutes
	}
	// < hour
	if beforeNow(time.Hour) {
		return TimeLessThanHour
	}
	// < 6 hours
	if beforeNow(time.Hour * 6) {
		return TimeLessThanSixHours
	}
	// < 12 hours
	if beforeNow(time.Hour * 12) {
		return TimeLessThanTwelveHours
	}
	// < day
	if beforeNow(time.Hour * 24) {
		return TimeLessThanDay
	}
	// < 2 days
	if beforeNow(time.Hour * 48) {
		return TimeLessThanTwoDays
	}
	// < week
	if beforeNow(time.Hour * 24 * 7) {
		return TimeLessThanWeek
	}
	// < 2 weeks
	if beforeNow(time.Hour * 24 * 7 * 2) {
		return TimeLessThanTwoWeeks
	}
	// < month
	if beforeNow(time.Hour * 24 * 7 * 4) {
		return TimeLessThanMonth
	}
	// < 2 months
	if beforeNow(time.Hour * 24 * 7 * 4 * 2) {
		return TimeLessThanTwoMonths
	}
	// < 6 months
	if beforeNow(time.Hour * 24 * 7 * 4 * 6) {
		return TimeLessThanSixMonths
	}
	// < year
	if beforeNow(time.Hour * 24 * 7 * 4 * 12) {
		return TimeLessThanYear
	}
	return TimeOlderThanAYear
}

// classifyText compares base to the query sting.
// If the case is determined is set using the sensitive bool.
func classifyText(base, query string, sensitive bool) uint {
	// remove leading dot from base if not found in query
	if r := []rune(query); len(r) > 0 && r[0] != '.' {
		if r := []rune(base); len(r) > 0 && r[0] == '.' {
			base = strings.Replace(base, ".", "", 1)
		}
	}
	// equals
	if base == query {
		return StrEquals
	}
	base = strings.ToLower(base)
	query = strings.ToLower(query)
	// equals wrong case
	if base == query {
		if sensitive {
			return StrEqualsWrongCase
		}
		return StrEquals
	}
	// starts with
	if strings.HasPrefix(base, query) {
		return StrStartsWith
	}
	// ends with
	if strings.HasSuffix(base, query) {
		return StrEndsWith
	}
	// does base even contain s?
	if strings.Contains(base, query) {
		return StrContains
	}
	return similarity(base, query)
}

func similarity(base, query string) uint {
	baseLen := utf8.RuneCountInString(base)
	// don't compare too short words
	if baseLen < 3 {
		return NoMatch
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
	// remove leading dot rune
	if len(runes) > 0 && runes[0] == '.' {
		runes = runes[1:]
	}
	searchLen := len(runes)
	for k, v := range base {
		if k < searchLen && v == runes[k] {
			continue
		}
		diff++
	}
	if diff <= maxdiff {
		return StrSimilar
	}
	return NoMatch
}
