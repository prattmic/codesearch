package search

import (
	"regexp/syntax"
	"runtime"

	"github.com/google/codesearch/index"
	pt "github.com/monochromegane/the_platinum_searcher"
)

// Result describes a search result.
type Result struct {
	Path    string
	Matches []*pt.Match
}

// Searcher can search with a given index.
type Searcher struct {
	idx *index.Index
}

// NewSearcher creates a Searcher for the provided index.
func NewSearcher(file string) *Searcher {
	return &Searcher{
		idx: index.Open(file),
	}
}

// Search returns matches for the given regexp.
func (s *Searcher) Search(regexp string) ([]Result, error) {
	// Package index needs a regexp from the syntax package.
	re, err := syntax.Parse(regexp, syntax.POSIX)
	if err != nil {
		return nil, err
	}

	// While package pt wants us to use their function to create a pattern.
	pat, err := pt.NewPattern(regexp, "", false, false, true)
	if err != nil {
		return nil, err
	}

	// Find candidate files.
	fileids := s.idx.PostingQuery(index.RegexpQuery(re))

	// Start searching. Grep takes files to search on in and sends
	// results to out.
	in := make(chan *pt.GrepParams, 16)
	out := make(chan *pt.PrintParams, 16)
	go pt.Grep(in, out, &pt.Option{
		Proc: runtime.NumCPU(),
	})

	// Send files to search.
	go func() {
		for _, id := range fileids {
			in <- &pt.GrepParams{
				Path: s.idx.Name(id),
				Pattern: pat,
			}
		}
		// Grep stops when in is empty and closed.
		close(in)
	}()

	// Gather results.
	var results []Result
	for p := range out {
		if len(p.Matches) > 0 {
			results = append(results, Result{
				Path:    p.Path,
				Matches: p.Matches,
			})
		}
	}

	return results, nil
}
