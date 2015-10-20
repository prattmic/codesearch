// Binary get simply provides a basic test for package get for now.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/prattmic/codesearch/pkg/get"
)

var usageMessage = `usage: get package gopath

Download a package and all of its dependencies into gopath.`

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, usageMessage)
		flag.PrintDefaults()
	}
	flag.Parse()

	if len(flag.Args()) != 2 {
		flag.Usage()
		os.Exit(1)
	}

	pkg := flag.Arg(0)
	gopath := flag.Arg(1)

	if err := get.Get(pkg, gopath); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
