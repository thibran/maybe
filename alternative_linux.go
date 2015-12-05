package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

type search struct {
	key     string
	filters []filterFn
}

type filterFn func(p string) bool

// TODO whitelist folders
// TODO check current directory too

func newSearch(s string) *search {
	var filters []filterFn
	if home := os.Getenv("HOME"); len(home) != 0 {
		filters = append(filters, func(p string) bool {
			return prefixFilter(home, p)
		})
	}
	filesFilter := func(p string) bool {
		return prefixFilter("/files", p)
	}
	tmpFilter := func(p string) bool {
		return maxDepthFilter("/tmp", p, 2)
	}
	filters = append(filters, filesFilter, tmpFilter)
	return &search{
		key:     s,
		filters: filters,
	}
}

// alternative returns the best match it could find for the search keyword.
func (s *search) alternative() string {
	arr := s.list()
	if len(arr) == 0 {
		return ""
	}
	entries := s.filter(arr, 10)
	return s.selectBest(entries)
}

// list all files & folders with search key
func (s *search) list() []string {
	args := []string{"--limit", "800", "--basename", s.key}
	b, err := exec.Command("locate", args...).Output()
	if err != nil {
		return []string{}
	}
	res := string(b)
	res = strings.TrimSpace(res)
	return strings.Split(res, "\n")
}

// filter returns the n best result paths filterd by filterFn.
// Short path are prefered over long paths.
func (s *search) filter(paths []string, n int) []string {
	set := make(stringSet)
	for _, p := range paths {
		for _, filter := range s.filters {
			s.maybeAdd(set, filter, p)
		}
	}
	e := &entries{
		search: s.key,
		arr:    set.items(),
	}
	sort.Sort(e)
	if len(e.arr) > n {
		e.arr = e.arr[:n]
	}
	return e.arr
}

func (s *search) maybeAdd(set stringSet, filter filterFn, p string) {
	if filter(p) {
		set.append(p)
	} else if r := s.inPathSegment(p); len(r) > 0 {
		if filter(r) {
			set.append(r)
		}
	}
}

func (s *search) printResultlist() {
	in := s.list()
	var arr []string
	for _, row := range s.filter(in, 15) {
		if isDir(row) {
			arr = append(arr, row)
		}
	}
	fmt.Println(strings.Join(arr, "\n"))
}

// selectBest returns the first directory entry.
func (s *search) selectBest(arr []string) string {
	for _, row := range arr {
		if isDir(row) {
			return row
		}
	}
	return ""
}

// inPathSegment returns a path if path segment starts with the keyword.
// Otherwise an empty string is returned.
func (s *search) inPathSegment(p string) string {
	arr := strings.Split(p, "/")
	arrLen := len(arr)
	if arrLen <= 1 {
		return ""
	}
	// cut last segment
	arr = arr[:arrLen-1]
	key := strings.ToLower(s.key)
	for i, seq := range arr {
		seq := strings.ToLower(seq)
		if strings.HasPrefix(seq, key) {
			return strings.Join(arr[:i+1], "/")
		}
	}
	return ""
}

// prefixFilter takes a prefix and path string.
func prefixFilter(prefix, p string) bool {
	if strings.HasPrefix(p, prefix) {
		if hiddenFolderInPath(p) {
			return false
		}
		if looksLikeAFile(p) {
			return false
		}
		return true
	}
	return false
}

// maxDepthFilter takes a prefix and path string.
func maxDepthFilter(prefix, p string, n int) bool {
	if seq := strings.Split(p, "/"); len(seq) > n {
		seq = seq[:n]
		p = strings.Join(seq, "/")
	}
	return prefixFilter(prefix, p)
}

func shorten(p string) string {
	s := strings.Split(p, "/")
	return strings.Join(s[len(s)-2:], "/")
}

func hiddenFolderInPath(p string) bool {
	for _, v := range strings.Split(p, "/") {
		if strings.HasPrefix(v, ".") {
			return true
		}
	}
	return false
}

func looksLikeAFile(s string) bool {
	return len(filepath.Ext(s)) > 0
}

func isDir(p string) bool {
	f, err := os.Open(p)
	if err != nil {
		return false
	}
	defer f.Close()
	s, err := f.Stat()
	if err != nil {
		return false
	}
	return s.IsDir()
}

// stringSet set of entry items.
type stringSet map[string]struct{}

type entries struct {
	arr    []string
	search string
}

func (r entries) Len() int {
	return len(r.arr)
}

func (r entries) Swap(i, j int) {
	r.arr[i], r.arr[j] = r.arr[j], r.arr[i]
}

func (r entries) Less(i, j int) bool {
	key := strings.ToLower(r.search)
	a := filepath.Base(r.arr[i])
	b := filepath.Base(r.arr[j])
	a = strings.ToLower(a)
	b = strings.ToLower(b)
	prefixA := strings.HasPrefix(a, key)
	prefixB := strings.HasPrefix(b, key)
	if prefixA && prefixB {
		return len(r.arr[i]) < len(r.arr[j])
	}
	return prefixA
}

func (s stringSet) append(arr ...string) {
	for _, v := range arr {
		s[v] = struct{}{}
	}
}

func (s stringSet) items() []string {
	arr := make([]string, len(s))
	var i = 0
	for k := range s {
		arr[i] = k
		i++
	}
	return arr
}
