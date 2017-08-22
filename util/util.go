// +build !android !darwin !windows

package util

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/thibran/maybe/pref"
)

const osSep = string(os.PathSeparator)

// Logf prints to stdout if pref.Verbose is true.
func Logf(format string, a ...interface{}) {
	if pref.Verbose {
		fmt.Printf(format, a...)
	}
}

// Logln prints to stdout if pref.Verbose is true.
func Logln(a ...interface{}) {
	if pref.Verbose {
		fmt.Println(a...)
	}
}

// ShortenPath when necessary, tries to keep the last two segments intact.
func ShortenPath(p string, max int) string {
	rlen := utf8.RuneCountInString
	{
		pathLen := rlen(p)
		// shorter than max
		if pathLen == 0 || pathLen <= max {
			return p
		}
	}
	// handle last segment too long
	a := strings.Split(p, osSep)
	last := filepath.Join(a[len(a)-2:]...)
	lastLen := rlen(last)
	if lastLen > max {
		last = "..." + last[lastLen+3-max:]
		return last
	}
	// try to keep start, replace mid path segments
	res := osSep
	// keep space for seperators
	max = max - (len(a) - 2)
	dotWithSepLen := 4
	for _, v := range a[:len(a)-2] {
		if rlen(res)+rlen(v)+lastLen+dotWithSepLen > max {
			break
		}
		res = filepath.Join(res, v)
	}
	return filepath.Join(res, "...", last)
}

// TermWidth returns the terminal with.
func TermWidth() (int, error) {
	cmd := exec.Command("tput", "cols")
	buf, err := cmd.Output()
	if err != nil {
		return 0, err
	}
	buf = bytes.TrimSpace(buf)
	width, err := strconv.Atoi(string(buf))
	if err != nil {
		return 0, err
	}
	return width, nil
}

// NormalOrVerbose returns the normal string if Pref.Verbose is not true.
func NormalOrVerbose(normal, verb string) string {
	if !pref.Verbose {
		return normal
	}
	return verb
}
