package turtle_test

import (
	"embed"
	_ "embed"
	"fmt"
	"github.com/0x51-dev/rdf/internal/project"
	"github.com/0x51-dev/rdf/internal/testsuite"
	nt "github.com/0x51-dev/rdf/ntriples"
	ttl "github.com/0x51-dev/rdf/turtle"
	"os"
	"testing"
)

var (
	//go:embed testdata/example1.ttl
	example1 string

	//go:embed testdata/*.ttl
	examples embed.FS

	//go:embed testdata/suite/manifest.ttl
	rawManifest string

	//go:embed testdata/suite/*
	suite embed.FS
)

func Example_example1() {
	doc, _ := ttl.ParseDocument(example1)
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
	// Amount of Triples in each example (manually counted).
	triples := []int{
		7, 1, 1, 2, 1, 2, 2, 2, 14, 3,
		10, 2, 2, 3, 2, 3, 6, 3, 4, 2,
		3, 3, 2, 5, 2, 8, 3,
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
			doc, err := ttl.ParseDocument(raw)
			if err != nil {
				t.Fatal(err)
			}
			if len(doc) != n {
				t.Fatal(n, len(doc))
			}

			{ // fmt.Stringer
				doc2, err := ttl.ParseDocument(doc.String())
				if err != nil {
					t.Fatal(doc.String())
				}
				if len(doc2) != n {
					t.Fatal(n, len(doc2))
				}
			}
		})
	}
}

func TestSuite(t *testing.T) {
	nt.ToggleValidation(false)

	manifest, err := testsuite.LoadManifest(rawManifest)
	if err != nil {
		t.Fatal(err)
	}

	report := project.NewReport(ttl.IRI{Prefixed: true, Value: "turtletest:"})
	for _, k := range manifest.Keys {
		e := manifest.Entries[k]
		raw, err := suite.ReadFile(fmt.Sprintf("testdata/suite/%s", e.Action))
		if err != nil {
			t.Fatal(err)
		}
		doc, err := ttl.ParseDocument(string(raw))
		switch e.Type {
		case "rdft:TestTurtlePositiveSyntax":
			t.Run(e.Name, func(t *testing.T) {
				if err != nil {
					report.AddTest(e.Name, testsuite.Failed)
					t.Fatal(err)
				}

				if ntr, err := ttl.EvaluateDocument(doc); err != nil {
					report.AddTest(e.Name, testsuite.Failed)
					t.Fatal(err)
				} else {
					if _, err := nt.ParseDocument(ntr.String()); err != nil {
						report.AddTest(e.Name, testsuite.Failed)
						t.Fatal(ntr.String())
					}
				}

				// fmt.Stringer
				doc2, err := ttl.ParseDocument(doc.String())
				if err != nil {
					report.AddTest(e.Name, testsuite.Failed)
					t.Fatal(err)
				}
				if len(doc) != len(doc2) {
					report.AddTest(e.Name, testsuite.Failed)
					t.Fatal(len(doc), len(doc2))
				}

				report.AddTest(e.Name, testsuite.Passed)
			})
		case "rdft:TestTurtleNegativeSyntax":
			t.Run(e.Name, func(t *testing.T) {
				if err == nil && ttl.ValidateDocument(doc) {
					report.AddTest(e.Name, testsuite.Failed)
					t.Fatal("expected error")
				}

				report.AddTest(e.Name, testsuite.Passed)
			})
		case "rdft:TestTurtleEval":
			t.Run(e.Name, func(t *testing.T) {
				if err != nil {
					report.AddTest(e.Name, testsuite.Failed)
					t.Fatal(err)
				}

				// fmt.Stringer
				doc2, err := ttl.ParseDocument(doc.String())
				if err != nil {
					report.AddTest(e.Name, testsuite.Failed)
					t.Fatal(err)
				}
				if len(doc) != len(doc2) {
					report.AddTest(e.Name, testsuite.Failed)
					t.Fatal(len(doc), len(doc2))
				}

				raw, err := suite.ReadFile(fmt.Sprintf("testdata/suite/%s", e.Result))
				if err != nil {
					t.Fatal(err)
				}
				ntr, err := nt.ParseDocument(string(raw))
				if err != nil {
					t.Fatal(err)
				}
				ntr2, err := ttl.EvaluateDocument(doc)
				if err != nil {
					report.AddTest(e.Name, testsuite.Failed)
					t.Fatal(err)
				}
				if len(ntr) != len(ntr2) {
					report.AddTest(e.Name, testsuite.Failed)
					t.Fatal(len(ntr), len(ntr2))
				}
				if _, err := nt.ParseDocument(ntr2.String()); err != nil {
					report.AddTest(e.Name, testsuite.Failed)
					t.Fatal(ntr2.String())
				}
				report.AddTest(e.Name, testsuite.Passed)
			})
		case "rdft:TestTurtleNegativeEval":
			t.Run(e.Name, func(t *testing.T) {
				if err != nil {
					// If we can not parse the document, we can not evaluate it...
					report.AddTest(e.Name, testsuite.Passed)
					return
				}

				// fmt.Stringer
				doc2, err := ttl.ParseDocument(doc.String())
				if err != nil {
					report.AddTest(e.Name, testsuite.Failed)
					t.Fatal(err)
				}
				if len(doc) != len(doc2) {
					report.AddTest(e.Name, testsuite.Failed)
					t.Fatal(len(doc), len(doc2))
				}

				if _, err := ttl.EvaluateDocument(doc); err == nil {
					report.AddTest(e.Name, testsuite.Failed)
					t.Fatal(doc)
				}
			})
		default:
			t.Fatal("unknown test type", e.Type)
		}
	}

	t.Log("Total tests:", report.Len())
	if os.Getenv("TEST_SUITE_REPORT") == "true" {
		_ = os.WriteFile("testdata/suite/report.ttl", []byte(report.String()), 0644)
	}
}
