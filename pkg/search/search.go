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

	// SnippetBefore is the portion of the snippet before the match.
	// It includes text on previous context lines and the same line
	// as the match.
	SnippetBefore string

	// SnippetMatch is the exact matching string.
	SnippetMatch string

	// SnippetAfter is the portion of the snippet after the match.
	// It includes text on the same line as the match as well as
	// following context lines.
	SnippetAfter string
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
		start := m.Num

		var snippetBefore string
		for i, l := range m.Befores {
			if l.Num < start {
				start = l.Num
			}
			if i != 0 {
				snippetBefore += "\n"
			}
			snippetBefore += l.Str
		}

		if len(m.Befores) > 0 {
			snippetBefore += "\n"
		}

		// Find the exact match on the matching line.
		i := pattern.FindStringIndex(m.Str)

		snippetBefore += m.Str[:i[0]]
		snippetMatch := m.Str[i[0]:i[1]]
		snippetAfter := m.Str[i[1]:]

		for _, l := range m.Afters {
			snippetAfter += "\n" + l.Str
		}

		matches = append(matches, Match{
			Start:         start,
			SnippetBefore: snippetBefore,
			SnippetMatch:  snippetMatch,
			SnippetAfter:  snippetAfter,
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
