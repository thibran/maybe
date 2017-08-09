package folder

import (
	"os"
	"strings"
	"time"
)

// Folder entry.
type Folder struct {
	Path        string
	UpdateCount uint32      // counts how often the folder has been updated
	Times       []time.Time // last MaxTimesEntries updates
}

// New folder object.
func New(path string, times ...time.Time) *Folder {
	if path == "" {
		panic("NewFolder - empty path is prohibited")
	}
	if len(times) == 0 {
		panic("NewFolder - must have at last one []time entry")
	}
	return &Folder{
		Path:        path,
		UpdateCount: 1,
		Times:       times,
	}
}

// ResourceCheckerFn converter to ResourceChecker interface.
type ResourceCheckerFn func(path string) bool

// DoesExist implementation for ResourceCheckerFn.
func (fn ResourceCheckerFn) DoesExist(path string) bool {
	return fn(path)
}

// CheckerFn returns a ResourceChecker function.
func CheckerFn() ResourceCheckerFn {
	// func returns true if the path points to a exsisting folder.
	return ResourceCheckerFn(func(path string) bool {
		if strings.TrimSpace(path) == "" {
			return false
		}
		fi, err := os.Stat(path)
		return !(os.IsNotExist(err) || !fi.IsDir())
	})
}
