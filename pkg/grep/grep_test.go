package grep

import (
	"bytes"
	"reflect"
	"regexp"
	"testing"
)

func TestSearch(t *testing.T) {
	b := bytes.NewBufferString("hello\nworld")
	g, err := New(b)
	if err != nil {
		t.Fatalf("New got err %v want nil", err)
	}

	r := regexp.MustCompile("world")
	m := g.Search(r)
	expect := []Match{
		{
			LineNum:  2,
			FullLine: []byte("world"),
			Match:    []byte("world"),
		},
	}
	if !reflect.DeepEqual(m, expect) {
		t.Errorf("Search got %+v want %+v", m, expect)
	}
}
