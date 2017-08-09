package repo

import (
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"strings"
	"thibaut/maybe/pref"
	"thibaut/maybe/rated"
	"thibaut/maybe/rated/folder"
	"thibaut/maybe/util"
	"time"
)

const osSep = string(os.PathSeparator)

var (
	// ErrNoResult - search has no result
	ErrNoResult = errors.New("no result")
	ignoreSlice = []string{".git", ".hg", ".svn", ".bzr"}
	errNoFile   = fmt.Errorf("file not found")
)

// Repo content is saved to the disk.
type Repo struct {
	m          rated.Map
	dataPath   string
	maxEntries int
}

// New repo object.
func New(path string, maxEntries int) *Repo {
	return &Repo{
		m:          make(rated.Map),
		dataPath:   path,
		maxEntries: maxEntries,
	}
}

// Walk adds directories from root, for count
// of Repo.maxEntries, osWalker.lvlDeep.
func (r *Repo) Walk(root string) {
	w := newWalker(r, root)
	err := filepath.Walk(root, w.toWalkFunc())
	if err != nil && err != errWalkMaxDirEntries {
		log.Fatalln(err)
	}
	if tmp := os.TempDir(); tmp != "" {
		r.updateOrAdd(tmp, w.now, true)
	}
}

// Add path to repo. If the path is known, the repo data is updated, else
// a new entry will be created.
func (r *Repo) Add(path string, t time.Time) {
	segments := strings.Split(path, osSep)
	len := len(segments)
	for i := 0; i < len-1; i++ {
		path = strings.Join(segments[:len-i], osSep)
		if isInIgnoreList(segments[len-1-i]) {
			util.Logf("ignore: %s\n", path)
			continue
		}
		r.updateOrAdd(path, t, i > 0)
	}
}

// updateOrAdd to repository. Sub-folders are added, when unknown, but
// their timestamps are not updated.
func (r *Repo) updateOrAdd(path string, t time.Time, subfolder bool) {
	f, ok := r.m[path]
	// new folder object
	if !ok {
		var sf string
		if subfolder {
			sf = "sub-"
		}
		util.Logf("new %sfolder: %s\n", sf, path)
		r.m[path] = folder.New(path, t)

		// guarantee folder limit holds
		if len(r.m) > r.maxEntries {
			r.m.RemoveOldest(r.maxEntries - r.maxEntries/3)
		}
		return
	}
	// update existing folder object
	if subfolder {
		return
	}
	util.Logf("update timestamps: %q\n", path)
	if f.UpdateCount < math.MaxUint32 {
		f.UpdateCount++
	}
	f.Times = append(f.Times, t)
	f.Times = rated.SortAndCut(f.Times...) // keep only data.MaxTimesEntries
	r.m[path] = f
}

// ResourceChecker returns true when a resource exists.
type ResourceChecker interface {
	DoesExist(string) bool
}

// Search repo for query.
func (r *Repo) Search(ch ResourceChecker, q pref.Query) (*rated.Rated, error) {
	a := r.m.Search(q.Last, func(a rated.Slice) { a.Sort() })
	a.FilterInPathOf(q.Start)
	for _, v := range a {
		// keep not found folders, they might re-exist in future
		if ch.DoesExist(v.Path) {
			// if checkFolder(v.folder.Path) {
			return v, nil
		}
	}
	return nil, ErrNoResult
}

// List returns all RatedSlice for the query q.
func (r *Repo) List(q pref.Query, cutLong bool) rated.Slice {
	a := r.m.Search(q.Last, func(a rated.Slice) { a.Sort() })
	a.FilterInPathOf(q.Start)
	a.CutLongPaths(cutLong)
	return a
}

// Size of the repository.
func (r *Repo) Size() int { return len(r.m) }
