package rated

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/thibran/maybe/classify"
	"github.com/thibran/maybe/pref"
	"github.com/thibran/maybe/rated/folder"
	"github.com/thibran/maybe/util"
)

// Rated object.
type Rated struct {
	*folder.Folder
	*classify.Rating
}

const osSep = string(os.PathSeparator)

// New creates a new rated folder object.
func New(f *folder.Folder, query string) (*Rated, error) {
	if f == nil {
		return nil, fmt.Errorf("rated.New - *Folder is nil")
	}
	r, err := classify.NewRating(query, path.Base(f.Path), f.Times...)
	if err != nil {
		return nil, fmt.Errorf("rated.New - %v", err)
	}
	return &Rated{
		Folder: f,
		Rating: r,
	}, nil
}

// Slice is an alias for Slice.
type Slice []*Rated

// Sort a RatedFolders.
func (rs Slice) Sort() {
	var pi, pj uint
	sort.Slice(rs, func(i, j int) bool {
		pi = rs[i].Points()
		pj = rs[j].Points()
		if pi == pj {
			if rs[i].UpdateCount == rs[j].UpdateCount {
				return rs[i].Path < rs[j].Path
			}
			return rs[i].UpdateCount > rs[j].UpdateCount
		}
		return pi > pj
	})
}

// FilterInPathOf returns only entries where the path
// contains the start-string in the non-last segment.
// When start is empty nothing is changed.
func (rs *Slice) FilterInPathOf(start string) {
	start = strings.TrimSpace(strings.ToLower(start))
	if start == "" {
		return
	}
	var a Slice
	for _, f := range *rs {
		// ignore the last path-segment
		// path /bar/src/foo becomes /bar/src/
		pathStart, _ := filepath.Split(f.Path)
		pathStart = strings.ToLower(pathStart)
		for _, seg := range strings.Split(pathStart, osSep) {
			if len(seg) == 0 {
				continue
			}
			trimmed := strings.TrimSuffix(seg, start)
			// exact match or non-suffix match
			if len(trimmed) == 0 || strings.Contains(trimmed, start) {
				a = append(a, f)
				break
			}
		}
	}
	*rs = a
}

// CutLongPaths if too long.
func (rs *Slice) CutLongPaths(cutLong bool) {
	if !cutLong {
		return
	}
	var a Slice
	// use terminal width, when possible
	maxLineLen := 64
	if w, err := util.TermWidth(); err == nil {
		if pref.Verbose && w-16 > 0 {
			maxLineLen = w - 16
		} else if !pref.Verbose && w-10 > 0 {
			maxLineLen = w - 10
		}
	}
	for _, rf := range *rs {
		rf.Path = util.ShortenPath(rf.Path, maxLineLen)
		a = append(a, rf)
	}
	*rs = a
}
