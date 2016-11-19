package search

import (
	"bytes"
	"os"
	"regexp"
	"regexp/syntax"
	"sync"

	"github.com/google/codesearch/index"
	"github.com/prattmic/codesearch/pkg/grep"
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

// NewResult builds a result from a slice of grep.Match.
func MakeResult(path string, re *regexp.Regexp, grepMatches []grep.Match) Result {
	var matches []Match
	for _, m := range grepMatches {
		start := m.LineNum - len(m.ContextBefore)

		snippetBefore := string(bytes.Join(m.ContextBefore, []byte{'\n'}))
		if len(m.ContextBefore) > 0 {
			snippetBefore += "\n"
		}

		// Find the exact match on the matching line.
		i := re.FindIndex(m.FullLine)

		snippetBefore += string(m.FullLine[:i[0]])
		snippetMatch := string(m.FullLine[i[0]:i[1]])
		snippetAfter := string(m.FullLine[i[1]:])

		if len(m.ContextAfter) > 0 {
			snippetAfter += "\n" + string(bytes.Join(m.ContextAfter, []byte{'\n'}))
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
// The results are unordered.
func (s *Searcher) Search(opts Options) ([]Result, error) {
	// Package index needs a regexp from the syntax package.
	syntaxRe, err := syntax.Parse(opts.Regexp, syntax.POSIX)
	if err != nil {
		return nil, err
	}

	// While package grep wants us to use the regex package.
	re, err := regexp.Compile(opts.Regexp)
	if err != nil {
		return nil, err
	}

	// Find candidate files.
	fileids := s.idx.PostingQuery(index.RegexpQuery(syntaxRe))

	// Grep all the files.
	rChan := make(chan Result, 10)
	var wg sync.WaitGroup
	for _, id := range fileids {
		path := s.idx.Name(id)

		wg.Add(1)
		go func() {
			defer wg.Done()
			f, err := os.Open(path)
			if err != nil {
				return
			}

			g, err := grep.New(f)
			if err != nil {
				return
			}

			// Copy the regexp to avoid lock contention.
			m := g.Search(re.Copy(), opts.Context)
			if len(m) > 0 {
				rChan <- MakeResult(path, re, m)
			}
		}()
	}

	// Collect the results.
	var results []Result
	var wg2 sync.WaitGroup
	wg2.Add(1)
	go func() {
		defer wg2.Done()
		for {
			r, ok := <-rChan
			if !ok {
				return
			}
			results = append(results, r)
		}
	}()

	wg.Wait()
	close(rChan)
	wg2.Wait()

	return results, nil
}
