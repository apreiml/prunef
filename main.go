package main

import (
	"flag"
	"fmt"
	"git.sr.ht/apreiml/prunef/ctime"
	"os"
	"time"
)

func main() {
	var fl = flag.Int("flagname", 1234, "help message for flagname")
	flag.Parse()
	fmt.Printf("%d\n", *fl)

	if len(flag.Args()) != 1 {
		printUsage()
		os.Exit(1)
	}

	var formatString = flag.Args()[0]

	var val, err = time.Parse("thor.server.aktrophy.at-2019-11-02_19-00-03", formatString)
	fmt.Println(err)
	fmt.Println(val)

	var t = time.Now()
	fmt.Printf("%#v\n", t)
	ctime.Strptime("asdf", "asdf", &t)
}

func printUsage() {
	fmt.Printf("Wrong args\n")
}
