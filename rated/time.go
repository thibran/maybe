package rated

import (
	"sort"
	"time"
)

// MaxTimeEntries of time.Time entries in a TimeSlice.
const MaxTimeEntries = 6

// TimeSlice alias with time focused sort implementation.
type TimeSlice []*Rated

func (a TimeSlice) sort() {
	sort.Slice(a, func(i, j int) bool {
		if a[i].TimePoints == a[j].TimePoints {
			return a[i].UpdateCount > a[j].UpdateCount
		}
		return a[i].TimePoints > a[j].TimePoints
	})
}

// SortAndCut time entries and keep only MaxTimesEntries.
func SortAndCut(a ...time.Time) []time.Time {
	sort.Slice(a, func(i, j int) bool { return a[i].After(a[j]) })
	if len(a) > MaxTimeEntries {
		return a[:MaxTimeEntries]
	}
	return a
}
