package search

import (
	"regexp/syntax"

	"github.com/google/codesearch/index"
)

// Result describes a search result.
type Result struct {
	FilePath string
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
	re, err := syntax.Parse(regexp, syntax.POSIX)
	if err != nil {
		return nil, err
	}

	fileids := s.idx.PostingQuery(index.RegexpQuery(re))

	var results []Result
	for _, id := range fileids {
		results = append(results, Result{FilePath: s.idx.Name(id)})
	}

	return results, nil
}
