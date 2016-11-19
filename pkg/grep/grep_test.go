package grep

import (
	"bytes"
	"io"
	"reflect"
	"regexp"
	"testing"
)

func TestSearch(t *testing.T) {
	testCases := []struct {
		name string
		in   io.Reader
		r    *regexp.Regexp
		c    int
		e    []Match
	}{
		{
			name: "simple",
			in:   bytes.NewBufferString("hello\nworld"),
			r:    regexp.MustCompile("world"),
			c:    0,
			e: []Match{
				{
					LineNum:       2,
					Match:         []byte("world"),
					FullLine:      []byte("world"),
					ContextBefore: [][]byte{},
					ContextAfter:  [][]byte{},
				},
			},
		},
		{
			name: "no-matches",
			in:   bytes.NewBufferString("hello\nworld"),
			r:    regexp.MustCompile("foo"),
			c:    0,
			e:    []Match{},
		},
		{
			name: "multi-match",
			in:   bytes.NewBufferString("foo\nbar\nfabulous"),
			r:    regexp.MustCompile("^f"),
			c:    0,
			e: []Match{
				{
					LineNum:       1,
					Match:         []byte("f"),
					FullLine:      []byte("foo"),
					ContextBefore: [][]byte{},
					ContextAfter:  [][]byte{},
				},
				{
					LineNum:       3,
					Match:         []byte("f"),
					FullLine:      []byte("fabulous"),
					ContextBefore: [][]byte{},
					ContextAfter:  [][]byte{},
				},
			},
		},
		{
			name: "complete-context",
			in:   bytes.NewBufferString("one\ntwo\nthree\nfour"),
			r:    regexp.MustCompile("three"),
			c:    1,
			e: []Match{
				{
					LineNum:  3,
					Match:    []byte("three"),
					FullLine: []byte("three"),
					ContextBefore: [][]byte{
						[]byte("two"),
					},
					ContextAfter: [][]byte{
						[]byte("four"),
					},
				},
			},
		},
		{
			name: "partial-context",
			in:   bytes.NewBufferString("one\ntwo\nthree\nfour"),
			r:    regexp.MustCompile("three"),
			c:    3,
			e: []Match{
				{
					LineNum:  3,
					Match:    []byte("three"),
					FullLine: []byte("three"),
					ContextBefore: [][]byte{
						[]byte("one"),
						[]byte("two"),
					},
					ContextAfter: [][]byte{
						[]byte("four"),
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g, err := New(tc.in)
			if err != nil {
				t.Fatalf("New got err %v want nil", err)
			}

			m := g.Search(tc.r, tc.c)
			if !reflect.DeepEqual(m, tc.e) {
				t.Errorf("Search got %+v want %+v", m, tc.e)
			}
		})
	}
}
