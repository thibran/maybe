package rated

import (
	"testing"
	"time"

	"github.com/thibran/maybe/rated/folder"
)

func TestRemoveOldest(t *testing.T) {
	fn := func(p string, t time.Time) *folder.Folder {
		return folder.New(p, t)
	}
	now := time.Now()
	f1 := fn("/home/bar", now.Add(-time.Hour*18))
	f2 := fn("/home/zot", now.Add(-time.Hour*4))
	f3 := fn("/home/foo", now)
	tt := []struct {
		name       string
		keepValues int
		resultLen  int
		notInMap   string
	}{
		{name: "remove oldest", keepValues: 2,
			resultLen: 2, notInMap: "/home/bar"},

		{name: "no change", keepValues: 10,
			resultLen: 3, notInMap: "/aaa"},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			m := Map{f1.Path: f1, f2.Path: f2, f3.Path: f3}
			m.RemoveOldest(tc.keepValues)
			if len(m) != tc.resultLen {
				t.Fatalf("expected len(res) %d, got %v", tc.keepValues, len(m))
			}
			if _, ok := m[tc.notInMap]; ok {
				t.Fatalf("%s should not be in the map", tc.notInMap)
			}
		})
	}
}
