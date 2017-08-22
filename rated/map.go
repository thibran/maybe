package rated

import (
	"runtime"
	"sync"

	"github.com/thibran/maybe/classify"
	"github.com/thibran/maybe/rated/folder"
)

// Map type alias
type Map map[string]*folder.Folder

type sorterFn func(a Slice)

// Search for s and sort results.
func (m *Map) Search(query string, sort sorterFn) Slice {
	if len(*m) == 0 {
		return Slice{}
	}
	var wg sync.WaitGroup
	workers := runtime.NumCPU()
	if len(*m) < workers {
		workers = len(*m)
	}
	wg.Add(workers)

	tasks := createTasks(*m)
	results := make(chan *Rated, workers+1)
	// results := make(chan *Rated)
	go func() {
		wg.Wait()
		close(results)
	}()

	for i := 0; i < workers; i++ {
		go func() {
			for f := range tasks {
				rf, err := New(f, query)
				if err != nil || rf.Points() == classify.NoMatch {
					continue
				}
				results <- rf
			}
			wg.Done()
		}()
	}
	return collectResults(results, sort)
}

func createTasks(m Map) <-chan *folder.Folder {
	tasks := make(chan *folder.Folder)
	go func() {
		for _, folder := range m {
			tasks <- folder
		}
		close(tasks)
	}()
	return tasks
}

func collectResults(c <-chan *Rated, sort sorterFn) Slice {
	var a Slice
	for r := range c {
		a = append(a, r)
	}
	sort(a)
	return a
}

// RemoveOldest folders from map and keep newest n entries.
func (m *Map) RemoveOldest(n int) {
	// to time-folders
	var a TimeSlice
	for _, f := range *m {
		if rf, err := New(f, ""); err == nil {
			a = append(a, rf)
		}
	}
	// delete too old entries
	a.sort()
	if len(a) > n {
		a = a[:n]
	}
	m2 := make(Map, len(a))
	for _, rf := range a {
		m2[rf.Path] = rf.Folder
	}
	*m = m2
}
