package main

import (
	"fmt"
	"os"
	fp "path/filepath"
	"strings"
	"time"
)

type osWalker struct {
	root    string
	lvlDeep uint
	r       *Repo
	now     time.Time
	count   int
}

func newOSWalker(r *Repo, root string) *osWalker {
	var lvlDeep uint = 6
	if !verbose {
		fmt.Println("initialize folders...")
	} else {
		fmt.Printf("initialize folders, %d level deep, "+
			"max-entries: %d, root: %q\n",
			lvlDeep, r.maxEntries, root)
	}
	return &osWalker{
		root:    fp.Clean(root), // converts e.g. /foo/ to /foo
		lvlDeep: lvlDeep,
		r:       r,
		now:     time.Now(),
		count:   len(r.m),
	}
}

func (w *osWalker) add(path string) error {
	if w.count == w.r.maxEntries {
		return errWalkMaxDirEntries
	}
	w.count++
	w.r.updateOrAddPath(path, w.now, true)
	return nil
}

func (w *osWalker) walk(path string) error {
	return walk(w.root, path, w.lvlDeep)
}

func (w *osWalker) toWalkFunc() fp.WalkFunc {
	return fp.WalkFunc(func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("fromDirwalker - %v", err)
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

func walk(root, path string, lvlDeep uint) error {
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
