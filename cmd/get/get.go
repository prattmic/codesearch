// Binary codesearch simply provides a basic test for package get for now.
package main

import (
	"flag"
	"fmt"
	"os"
	"sort"

	"github.com/prattmic/codesearch/pkg/get"
)

var usageMessage = `usage: codesearch package

Prints all of the (recursive) dependencies of the package.`

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, usageMessage)
		flag.PrintDefaults()
	}
	flag.Parse()

	if len(flag.Args()) != 1 {
		flag.Usage()
		os.Exit(1)
	}

	deps, err := get.PackageDependencies(flag.Arg(0))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	sort.Strings(deps)
	for _, d := range deps {
		fmt.Println(d)
	}
}
