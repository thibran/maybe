package repo

import (
	"path/filepath"
	"testing"
)

func TestWalkHelper(t *testing.T) {
	// pref.Verbose = true
	root := "/foo"
	j := filepath.Join
	ignoreSlice = append(ignoreSlice, "_a")
	errSkip := filepath.SkipDir
	tt := []struct {
		name, path string
		exp        error
	}{
		{name: "ok", path: j(root, "bar"), exp: nil},
		{name: "too deep", path: j(root, "zot/bar"), exp: errSkip},
		{name: "ignore", path: j(root, "_a"), exp: errSkip},
		{name: "hidden folder", path: j(root, ".zot"), exp: errSkip},
		{name: "empty path", path: root, exp: nil},
	}
	var lvlDeep uint = 1
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := walkHelper(root, tc.path, lvlDeep)
			if err != tc.exp {
				t.Fatalf("exp %q, got %v - for: %q", tc.exp, err, tc.path)
			}
		})
	}
}

func TestIsMaxFolderLevel(t *testing.T) {
	root := "/foo"
	path := "/foo/zot/bar"
	if !isMaxFolderLevel(1, root, path) {
		t.Fail()
	}
}

// func TestWalk(t *testing.T) {
// 	verbose = true
// 	// root := "/home/tux"
// 	root := "/tmp/test"
// 	r := NewRepo("/baz/bar/zot", 10)
// 	r.Walk(root)
// 	if len(r.m) == 0 {
// 		t.Fail()
// 	}
// }
