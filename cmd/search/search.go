package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/prattmic/codesearch/pkg/search"
)

var index = flag.String("i", "", "Index to search")

var usageMessage = `usage: search -i index regexp

Search for regexp within index.`

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, usageMessage)
		flag.PrintDefaults()
	}
	flag.Parse()

	if *index == "" {
		flag.Usage()
		os.Exit(1)
	}

	if len(flag.Args()) != 1 {
		flag.Usage()
		os.Exit(1)
	}

	s := search.NewSearcher(*index, "")

	results, err := s.Search(search.Options{Regexp: flag.Arg(0)})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, r := range results {
		for _, m := range r.Matches {
			fmt.Printf("%s:%d: %s\n", r.Path, m.Start, m.Snippet)
		}
	}
}
