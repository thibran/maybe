package ratedfolder

import (
	"sort"
	"time"
)

// MaxTimeEntries of time.Time entries in a Times slice.
const MaxTimeEntries = 6

// RatedTimeFolders alias with time focused sort implementation.
type RatedTimeFolders []*RatedFolder

func (a RatedTimeFolders) sort() {
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
