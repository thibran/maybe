package rated

import (
	"testing"
	"time"

	"github.com/thibran/maybe/classify"
	"github.com/thibran/maybe/rated/folder"
)

func TestTimeSort(t *testing.T) {
	now := time.Now()
	rated := func(path string, timePoints uint, count uint32) *Rated {
		return &Rated{
			Rating: &classify.Rating{TimePoints: timePoints},
			Folder: &folder.Folder{Path: path, UpdateCount: count,
				Times: []time.Time{now}},
		}
	}
	tt := []struct {
		name, exp string
		folders   TimeSlice
	}{
		{name: "by time points", exp: "/home/foo", folders: TimeSlice{
			rated("/home/bar", 4, 1),
			rated("/home/foo", 10, 1),
			rated("/home/zot", 8, 1),
		}},
		{name: "by time count", exp: "/b", folders: TimeSlice{
			rated("/a", 20, 1),
			rated("/b", 20, 2),
		}},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.folders.sort(); tc.folders[0].Path != tc.exp {
				t.Fatalf("exp %q, got %q", tc.exp, tc.folders[0].Path)
			}
		})
	}
}

func TestSortAndCut(t *testing.T) {
	var i int
	now := time.Now()
	timeFn := func() time.Time {
		i++
		return now.Add(time.Hour + time.Duration(i))
	}
	var a []time.Time
	for len(a) <= MaxTimeEntries {
		a = append(a, timeFn())
	}
	a = SortAndCut(a...)

	if len(a) != MaxTimeEntries {
		t.Errorf("len(a) should be %d, got %d", MaxTimeEntries, len(a))
	}
	exp := now.Add(time.Hour + time.Duration(1))
	if a[0].Hour() != exp.Hour() {
		t.Fatalf("exp %q, got %q", exp, a[0])
	}
	exp = now.Add(time.Hour + time.Duration(6))
	if a[5].Hour() != exp.Hour() {
		t.Fatalf("exp %q, got %q", exp, a[5])
	}
}
