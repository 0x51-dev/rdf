package ntriples_test

import (
	_ "embed"
	"github.com/0x51-dev/rdf/ntriples"
	"testing"
)

var (
	//go:embed testdata/example1.nt
	example1 string

	//go:embed testdata/example2.nt
	example2 string

	//go:embed testdata/example3.nt
	example3 string

	//go:embed testdata/example4.nt
	example4 string

	//go:embed testdata/test.nt
	ntest string
)

func TestExamples(t *testing.T) {
	for _, test := range []struct {
		doc     string
		triples int
	}{
		{example1, 3},
		{example2, 1},
		{example3, 7},
		{example4, 2},
		{ntest, 30},
	} {
		doc, err := ntriples.ParseDocument(test.doc)
		if err != nil {
			t.Fatal(err)
		}
		if len(doc) != test.triples {
			t.Error(len(doc))
		}
	}
}

func TestSuite(t *testing.T) {
	t.Skip("TODO")
}
