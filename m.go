package main

import (
	"flag"
	"fmt"
	"os"
	"os/user"
	"strings"
)

// --version
// --sub     // subdir
// --history // last 10 entries
// --info    // show most used & total entry count
// --add     //

func main() {
	onlyList := flag.Bool("list", false, "list top 10 results for the keyword")
	printVersion := flag.Bool("version", false, "print version number")
	wlist := flag.String("whitelist", "", "folders not to ignore, sep by ':'")
	flag.Parse()
	if *printVersion {
		fmt.Println("version: 0.1.6")
		os.Exit(0)
	}
	if len(flag.Args()) == 0 {
		if u, err := user.Current(); err == nil {
			fmt.Println(u.HomeDir)
			os.Exit(0)
		} else {
			fmt.Println("E: no args specified + can't find home directory!")
			os.Exit(1)
		}
	}
	in := flag.Args()[0]
	if len(in) < 2 {
		fmt.Println("E: input too short, must be at last 2 characters long.")
		os.Exit(1)
	}
	var whitelist []string
	if len(*wlist) > 0 {
		whitelist = strings.Split(*wlist, ":")
	}
	s := newSearch(in, whitelist)
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
