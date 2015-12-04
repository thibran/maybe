package main

import (
	"flag"
	"fmt"
	"os"
)

// --version
// --sub     // subdir
// --history // last 10 entries
// --info    // show most used & total entry count
// --add     //

func main() {
	onlyList := flag.Bool("list", false, "list top 10 results for the keyword.")
	flag.Parse()
	if len(flag.Args()) == 0 {
		fmt.Println("E: no args specified!")
		os.Exit(1)
	}
	in := flag.Args()[0]
	if len(in) < 3 {
		fmt.Println("E: input too short, must be at last 3 characters long.")
		os.Exit(1)
	}
	s := newSearch(in)
	if *onlyList {
		s.printResultlist()
		os.Exit(0)
	}
	res := s.alternative()
	if len(res) == 0 {
		fmt.Println("*miaow*")
		os.Exit(1)
	}
	fmt.Println(res)
}
