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
		e    []Match
	}{
		{
			name: "simple",
			in:   bytes.NewBufferString("hello\nworld"),
			r:    regexp.MustCompile("world"),
			e: []Match{
				{
					LineNum:  2,
					FullLine: []byte("world"),
					Match:    []byte("world"),
				},
			},
		},
		{
			name: "no-matches",
			in:   bytes.NewBufferString("hello\nworld"),
			r:    regexp.MustCompile("foo"),
			e:    []Match{},
		},
		{
			name: "multi-match",
			in:   bytes.NewBufferString("foo\nbar\nfabulous"),
			r:    regexp.MustCompile("^f"),
			e: []Match{
				{
					LineNum:  1,
					FullLine: []byte("foo"),
					Match:    []byte("f"),
				},
				{
					LineNum:  3,
					FullLine: []byte("fabulous"),
					Match:    []byte("f"),
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

			m := g.Search(tc.r)
			if !reflect.DeepEqual(m, tc.e) {
				t.Errorf("Search got %+v want %+v", m, tc.e)
			}
		})
	}
}
