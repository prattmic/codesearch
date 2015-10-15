package search

import (
	"regexp"
	"regexp/syntax"
	"runtime"
	"strings"

	"github.com/google/codesearch/index"
	pt "github.com/monochromegane/the_platinum_searcher"
)

// Options are passed to Searcher.Search.
type Options struct {
	// Regexp is the search query.
	Regexp string

	// Context is the number of snippet lines to include
	// before and after the result.
	Context int
}

// Match describes a single match within a file.
type Match struct {
	// Start is the line number of the first line in the snippet.
	Start int

	// Snippet is a snippet from the file containing the match.
	Snippet string

	// Indicies is a slice of byte index pairs, one pair for each
	// match within Snippet.
	Indicies [][]int
}

// Result describes a search result.
type Result struct {
	// Path is the local path of the result file.
	Path string

	// Matches are the individual matches in the file.
	Matches []Match
}

// NewResult builds a result from a slice of pt.Match.
func MakeResult(path string, pattern *regexp.Regexp, ptMatches []*pt.Match) Result {
	var matches []Match
	for _, m := range ptMatches {
		var lines []string
		start := m.Num

		for _, l := range m.Befores {
			if l.Num < start {
				start = l.Num
			}
			lines = append(lines, l.Str)
		}

		lines = append(lines, m.Str)

		for _, l := range m.Afters {
			lines = append(lines, l.Str)
		}

		snippet := strings.Join(lines, "\n")

		var indicies [][]int
		r := pattern.FindAllStringIndex(snippet, -1)
		for _, i := range r {
			// The first two items are the match for the
			// entire expression.
			indicies = append(indicies, i[:2])
		}

		matches = append(matches, Match{
			Start:    start,
			Snippet:  snippet,
			Indicies: indicies,
		})
	}

	return Result{
		Path:    path,
		Matches: matches,
	}
}

// Searcher can search with a given index.
type Searcher struct {
	idx    *index.Index
	prefix string
}

// NewSearcher creates a Searcher for the provided index.
func NewSearcher(file string, prefix string) *Searcher {
	return &Searcher{
		idx:    index.Open(file),
		prefix: prefix,
	}
}

// Search returns matches for the given regexp.
func (s *Searcher) Search(opts Options) ([]Result, error) {
	// Package index needs a regexp from the syntax package.
	re, err := syntax.Parse(opts.Regexp, syntax.POSIX)
	if err != nil {
		return nil, err
	}

	// While package pt wants us to use their function to create a pattern.
	pat, err := pt.NewPattern(opts.Regexp, "", false, false, true)
	if err != nil {
		return nil, err
	}

	// Find candidate files.
	fileids := s.idx.PostingQuery(index.RegexpQuery(re))

	// Start searching. Grep takes files to search on in and sends
	// results to out.
	popt := pt.Option{
		Before: opts.Context,
		After:  opts.Context,
		Proc:   runtime.NumCPU(),
	}
	in := make(chan *pt.GrepParams, popt.Proc)
	out := make(chan *pt.PrintParams, popt.Proc)
	go pt.Grep(in, out, &popt)

	// Send files to search.
	go func() {
		for _, id := range fileids {
			in <- &pt.GrepParams{
				Path:    s.idx.Name(id),
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
			path := strings.TrimPrefix(p.Path, s.prefix)
			results = append(results, MakeResult(path, p.Pattern.Regexp, p.Matches))
		}
	}

	return results, nil
}
