package repo

import (
	"fmt"
	"os"
	fp "path/filepath"
	"strings"
	"thibaut/maybe/pref"
	"time"
)

var errWalkMaxDirEntries = fmt.Errorf("max entries reached")

type walker struct {
	root    string
	lvlDeep uint
	r       *Repo
	now     time.Time
	count   int
}

func newWalker(r *Repo, root string) *walker {
	var lvlDeep uint = 6
	if !pref.Verbose {
		fmt.Println("initialize folders...")
	} else {
		fmt.Printf("initialize folders, %d level deep, "+
			"max-entries: %d, root: %q\n",
			lvlDeep, r.maxEntries, root)
	}
	return &walker{
		root:    fp.Clean(root), // converts e.g. /foo/ to /foo
		lvlDeep: lvlDeep,
		r:       r,
		now:     time.Date(2000, time.January, 0, 0, 0, 0, 0, time.UTC),
		count:   len(r.m),
	}
}

func (w *walker) add(path string) error {
	if w.count == w.r.maxEntries {
		return errWalkMaxDirEntries
	}
	w.count++
	w.r.updateOrAdd(path, w.now, true)
	return nil
}

func (w *walker) walk(path string) error {
	return walkHelper(w.root, path, w.lvlDeep)
}

func walkHelper(root, path string, lvlDeep uint) error {
	fileName := fp.Base(path)
	if fileName == "." || fileName == osSep {
		return nil
	}
	// no hidden folders
	if strings.HasPrefix(fileName, ".") {
		return fp.SkipDir
	}
	// not too deep
	if isMaxFolderLevel(lvlDeep, root, path) {
		return fp.SkipDir
	}
	// no ignore-list folders
	if isInIgnoreList(fileName) {
		return fp.SkipDir
	}
	return nil
}

func (w *walker) toWalkFunc() fp.WalkFunc {
	return fp.WalkFunc(func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("toWalkFunc - %v", err)
			return fp.SkipDir
		}
		// no files
		if !fi.IsDir() {
			return nil
		}
		err = w.walk(path)
		if err == fp.SkipDir {
			return err
		}
		return w.add(path)
	})
}

func isMaxFolderLevel(levelToWalk uint, root, path string) bool {
	lvl := strings.Replace(path, root, "", 1)
	a := strings.Split(lvl, osSep)
	return len(a) > int(levelToWalk+1)
}

func isInIgnoreList(name string) bool {
	for _, ign := range ignoreSlice {
		if name == ign {
			return true
		}
	}
	return false
}
