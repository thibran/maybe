package main

import (
	"fmt"
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

func logWithPrefix(prefix string) func(a ...interface{}) {
	return func(a ...interface{}) {
		if verbose {
			a = append([]interface{}{fmt.Sprintf("%s -", prefix)}, a...)
			fmt.Println(a...)
		}
	}
}
