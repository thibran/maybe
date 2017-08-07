// +build !android !darwin !windows

package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"unicode/utf8"
)

func logf(format string, a ...interface{}) {
	if verbose {
		fmt.Printf(format, a...)
	}
}

func logln(a ...interface{}) {
	if verbose {
		fmt.Println(a...)
	}
}

// shortenPath when necessary, tries to keep the last two segments intact.
func shortenPath(p string, max int) string {
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

func termWidth() (int, error) {
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

func normalOrVerbose(normal, verb string) string {
	if !verbose {
		return normal
	}
	return verb
}
