package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"

	"github.com/prattmic/codesearch/pkg/search"
)

var (
	index       = flag.String("i", "", "Index to search")
	cpuProfile  = flag.String("cpu_profile", "", "Save a CPU profile to this file")
	heapProfile = flag.String("heap_profile", "", "Save a heap profile to this file")
)

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

	if *cpuProfile != "" {
		f, err := os.Create(*cpuProfile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	if *heapProfile != "" {
		f, err := os.Create(*heapProfile)
		if err != nil {
			log.Fatal("could not create heap profile: ", err)
		}
		defer pprof.Lookup("heap").WriteTo(f, 0)
	}

	s := search.NewSearcher(*index, "")

	results, err := s.Search(search.Options{Regexp: flag.Arg(0)})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, r := range results {
		for _, m := range r.Matches {
			snippet := m.SnippetBefore + m.SnippetMatch + m.SnippetAfter
			fmt.Printf("%s:%d: %s\n", r.Path, m.Start, snippet)
		}
	}
}
