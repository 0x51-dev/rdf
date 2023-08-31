package turtle_test

import (
	"embed"
	_ "embed"
	"fmt"
	"github.com/0x51-dev/rdf/internal/testsuite"
	"github.com/0x51-dev/rdf/turtle"
	"testing"
)

var (
	//go:embed testdata/example1.ttl
	example1 string

	//go:embed testdata/*.ttl
	examples embed.FS

	//go:embed testdata/suite/manifest.ttl
	rawManifest string

	//go:embed testdata/suite/*.ttl
	suite embed.FS
)

func Example_example1() {
	doc, _ := turtle.ParseDocument(example1)
	fmt.Println(doc)
	// Output:
	// @base <http://example.org/> .
	// @prefix rdf: <http://www.w3.org/1999/02/22-rdf-syntax-ns#> .
	// @prefix rdfs: <http://www.w3.org/2000/01/rdf-schema#> .
	// @prefix foaf: <http://xmlns.com/foaf/0.1/> .
	// @prefix rel: <http://www.perceive.net/schemas/relationship/> .
	// <#green-goblin> rel:enemyOf <#spiderman> ; a foaf:Person ; foaf:name "Green Goblin" .
	// <#spiderman> rel:enemyOf <#green-goblin> ; a foaf:Person ; foaf:name "Spiderman", "Человек-паук"@ru .
}

func TestExamples(t *testing.T) {
	// Amount of triples in each example (manually counted).
	triples := []int{
		2, 1, 1, 2, 1, 2, 1, 1, 9, 2,
		7, 1, 1, 2, 1, 2, 6, 2, 1, 1,
		1, 2, 1, 4, 1, 7, 1,
	}
	entries, err := examples.ReadDir("testdata")
	if err != nil {
		t.Fatal(err)
	}

	em := make(map[string]string)
	for _, f := range entries {
		b, err := examples.ReadFile("testdata/" + f.Name())
		if err != nil {
			t.Fatal(err)
		}
		em[f.Name()] = string(b)
	}
	for i, n := range triples {
		t.Run(fmt.Sprintf("example%d", i+1), func(t *testing.T) {
			raw := em[fmt.Sprintf("example%d.ttl", i+1)]
			doc, err := turtle.ParseDocument(raw)
			if err != nil {
				t.Fatal(err)
			}
			if len(doc.Triples) != n {
				t.Fatal(len(doc.Triples))
			}

			{ // fmt.Stringer
				doc2, err := turtle.ParseDocument(doc.String())
				if err != nil {
					t.Fatal(doc.String())
				}
				if len(doc2.Triples) != n {
					t.Fatal(len(doc2.Triples))
				}
			}
		})
	}
}

func TestSuite(t *testing.T) {
	// TODO: enable this test
	t.Skip("skipping test suite")

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
		doc, err := turtle.ParseDocument(string(raw))
		switch e.Type {
		case "rdft:TestTurtlePositiveSyntax", "rdft:TestTurtleEval":
			t.Run(e.Name, func(t *testing.T) {
				if err != nil {
					t.Fatal(err)
				}

				// fmt.Stringer
				doc2, err := turtle.ParseDocument(doc.String())
				if err != nil {
					t.Fatal(err)
				}
				if len(doc.Triples) != len(doc2.Triples) {
					t.Error(len(doc.Triples), len(doc2.Triples))
				}
			})
		case "rdft:TestTurtleNegativeSyntax", "rdft:TestTurtleNegativeEval":
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
