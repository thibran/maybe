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
	//args := []string{"--limit", "250", "--basename", s.key}
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
func (s *search) filter(paths []string, n int) entries {
	set := newEntrySet()
	for _, p := range paths {
		for _, filter := range s.filters {
			s.maybeAdd(set, filter, p)
		}
	}
	result := set.items()
	sort.Sort(result)
	if len(result) > n {
		result = result[:n]
	}
	return result
}

func (s *search) maybeAdd(set *entrySet, filter filterFn, p string) {
	if filter(p) {
		set.append(entry{v: p, search: s.key})
	} else if r := s.inPathSegment(p); len(r) > 0 {
		if filter(r) {
			set.append(entry{v: r, search: s.key})
		}
	}
}

func (s *search) printResultlist() {
	entries := s.list()
	var arr []string
	for _, r := range s.filter(entries, 15) {
		if isDir(r.v) {
			arr = append(arr, r.v)
		}
	}
	fmt.Println(strings.Join(arr, "\n"))
}

// selectBest returns the first directory entry.
func (s *search) selectBest(arr entries) string {
	for _, r := range arr {
		if isDir(r.v) {
			return r.v
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

// entrySet set of strings.
type entrySet struct {
	m map[entry]struct{}
}

type entry struct {
	v      string
	search string
}

type entries []entry

func (r entries) Len() int {
	return len(r)
}
func (r entries) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func (r entries) Less(i, j int) bool {
	a := filepath.Base(r[i].v)
	b := filepath.Base(r[j].v)
	key := strings.ToLower(r[i].search)
	a = strings.ToLower(a)
	b = strings.ToLower(b)
	prefixA := strings.HasPrefix(a, key)
	prefixB := strings.HasPrefix(b, key)
	// aName := shorten(r[i].v)
	// bName := shorten(r[j].v)
	if prefixA && prefixB {
		return len(r[i].v) < len(r[j].v)
	}
	return prefixA
}

func newEntrySet() *entrySet {
	return &entrySet{
		m: make(map[entry]struct{}),
	}
}

func (s *entrySet) append(arr ...entry) {
	for _, v := range arr {
		s.m[v] = struct{}{}
	}
}

func (s *entrySet) items() entries {
	arr := make(entries, len(s.m))
	var i = 0
	for k := range s.m {
		arr[i] = k
		i++
	}
	return arr
}
