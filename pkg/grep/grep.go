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

	// FullLine is the entire line containing the match.
	FullLine []byte

	// Match is the portion of the line that matched.
	Match []byte
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

// Search finds all lines matching regexp.
func (g *Grep) Search(r *regexp.Regexp) []Match {
	matches := make([]Match, 0)

	for i, l := range g.lines {
		m := r.Find(l)
		if m == nil {
			continue
		}

		matches = append(matches, Match{
			LineNum:  i + 1,
			FullLine: l,
			Match:    m,
		})
	}

	return matches
}
