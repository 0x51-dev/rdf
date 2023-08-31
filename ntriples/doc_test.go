package ntriples_test

import (
	"embed"
	_ "embed"
	"fmt"
	"github.com/0x51-dev/rdf/internal/testsuite"
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

	//go:embed testdata/suite/manifest.ttl
	rawManifest string

	//go:embed testdata/suite/*.nt
	suite embed.FS
)

func Example_example1() {
	doc, _ := ntriples.ParseDocument(example1)
	fmt.Println(doc)
	// Output:
	// <http://one.example/subject1> <http://one.example/predicate1> <http://one.example/object1> .
	// _:subject1 <http://an.example/predicate1> "object1" .
	// _:subject2 <http://an.example/predicate2> "object2" .
}

func TestExamples(t *testing.T) {
	for _, test := range []struct {
		doc     string
		triples int
	}{
		{example1, 3},
		{example2, 1},
		{example3, 7},
		{example4, 2},
	} {
		doc, err := ntriples.ParseDocument(test.doc)
		if err != nil {
			t.Fatal(err)
		}
		if len(doc) != test.triples {
			t.Error(len(doc))
		}

		{ // fmt.Stringer
			doc, err := ntriples.ParseDocument(doc.String())
			if err != nil {
				t.Fatal(err)
			}
			if len(doc) != test.triples {
				t.Error(len(doc))
			}
		}
	}
}

func TestSuite(t *testing.T) {
	manifest, err := testsuite.LoadManifest(rawManifest)
	if err != nil {
		t.Fatal(err)
	}
	for _, k := range manifest.Keys {
		e := manifest.Entries[k]
		raw, err := suite.ReadFile(fmt.Sprintf("testdata/suite/%s", e.Action))
		if err != nil {
			t.Fatal(err)
		}
		doc, err := ntriples.ParseDocument(string(raw))
		switch e.Type {
		case "rdft:TestNTriplesPositiveSyntax":
			t.Run(e.Name, func(t *testing.T) {
				if err != nil {
					t.Fatal(err)
				}

				// fmt.Stringer
				doc2, err := ntriples.ParseDocument(doc.String())
				if err != nil {
					t.Fatal(err)
				}
				if len(doc) != len(doc2) {
					t.Error(len(doc), len(doc2))
				}
			})
		case "rdft:TestNTriplesNegativeSyntax":
			t.Run(e.Name, func(t *testing.T) {
				if err == nil {
					t.Fatal("expected error")
				}
			})
		default:
			t.Fatal("unknown test type", e.Type)
		}
	}
}
