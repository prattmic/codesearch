// Package grep provides a file search implementation.
package grep

import (
	"bytes"
	"io"
	"io/ioutil"
	"regexp"
)

type Match struct {
	// LineNum is the line number of the match.
	LineNum int

	// Match is the portion of the line that matched.
	Match []byte

	// FullLine is the entire line containing the match.
	FullLine []byte

	// ContextBefore is the set of context lines before FullLine.
	ContextBefore [][]byte

	// ContextBefore is the set of context lines after FullLine.
	ContextAfter [][]byte
}

// Grep searches a file.
//
// Limitations:
//  * Stores entire file in memory and splits by line.
//  * Only single-line search.
type Grep struct {
	lines [][]byte
}

// New returns a new Grep.
func New(r io.Reader) (*Grep, error) {
	c, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return &Grep{
		lines: bytes.Split(c, []byte{'\n'}),
	}, nil
}

// context returns valid lines in the range [start, end). start and end may be
// outside the valid range [0, len(g.lines)).
func (g *Grep) context(start, end int) [][]byte {
	c := make([][]byte, 0, end-start)
	for ; start < end && start < len(g.lines); start++ {
		if start < 0 {
			continue
		}
		c = append(c, g.lines[start])
	}

	return c
}

// Search finds all lines matching regexp. context is the number of context
// lines to include before and after the match.
func (g *Grep) Search(r *regexp.Regexp, context int) []Match {
	matches := make([]Match, 0)

	for i, l := range g.lines {
		m := r.Find(l)
		if m == nil {
			continue
		}

		matches = append(matches, Match{
			LineNum:       i + 1,
			Match:         m,
			FullLine:      l,
			ContextBefore: g.context(i-context, i),
			ContextAfter:  g.context(i+1, i+1+context),
		})
	}

	return matches
}
