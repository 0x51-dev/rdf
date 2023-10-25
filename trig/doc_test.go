package trig_test

import (
	"embed"
	_ "embed"
	"fmt"
	"github.com/0x51-dev/rdf/internal/project"
	"github.com/0x51-dev/rdf/internal/testsuite"
	nq "github.com/0x51-dev/rdf/nquads"
	nt "github.com/0x51-dev/rdf/ntriples"
	"github.com/0x51-dev/rdf/trig"
	ttl "github.com/0x51-dev/rdf/turtle"
	"os"
	"testing"
)

var (
	//go:embed testdata/example1.trig
	example1 string

	//go:embed testdata/*.trig
	examples embed.FS

	//go:embed testdata/suite/manifest.ttl
	rawManifest string

	//go:embed testdata/suite/*
	suite embed.FS
)

func Example_example1() {
	doc, _ := trig.ParseDocument(example1)
	fmt.Println(doc)
	// Output:
	// @prefix ex: <http://www.example.org/vocabulary#> .
	// @prefix : <http://www.example.org/exampleDocument#> .
	// :G1 { :Monica a ex:Person ; ex:email <mailto:monica@monicamurphy.org> ; ex:hasSkill ex:Management, ex:Programming ; ex:homepage <http://www.monicamurphy.org> ; ex:name "Monica Murphy" . }
}

func TestExamples(t *testing.T) {
	// Amount of triples in each example (manually counted).
	triples := []int{
		3, 6, 7,
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
			raw := em[fmt.Sprintf("example%d.trig", i+1)]
			doc, err := trig.ParseDocument(raw)
			if err != nil {
				t.Fatal(err)
			}
			if len(doc) != n {
				t.Fatal(n, len(doc))
			}

			{ // fmt.Stringer
				doc2, err := trig.ParseDocument(doc.String())
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

	report := project.NewReport(ttl.IRI{Prefixed: true, Value: "http://www.w3.org/2013/TrigTests/manifest.ttl#"})
	for _, k := range manifest.Keys {
		e := manifest.Entries[k]
		raw, err := suite.ReadFile(fmt.Sprintf("testdata/suite/%s", e.Action))
		if err != nil {
			t.Fatal(err)
		}
		doc, err := trig.ParseDocument(string(raw))
		switch e.Type {
		case "rdft:TestTrigPositiveSyntax":
			t.Run(e.Name, func(t *testing.T) {
				if err != nil {
					report.AddTest(e.Name, testsuite.Failed)
					t.Fatal(err)
				}

				if ntr, err := trig.EvaluateDocument(doc); err != nil {
					report.AddTest(e.Name, testsuite.Failed)
					t.Fatal(err)
				} else {
					if _, err := nq.ParseDocument(ntr.String()); err != nil {
						report.AddTest(e.Name, testsuite.Failed)
						t.Fatal(ntr.String())
					}
				}

				// fmt.Stringer
				doc2, err := trig.ParseDocument(doc.String())
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
		case "rdft:TestTrigNegativeSyntax":
			t.Run(e.Name, func(t *testing.T) {
				if err == nil && trig.ValidateDocument(doc) {
					report.AddTest(e.Name, testsuite.Failed)
					t.Fatal("expected error")
				}

				report.AddTest(e.Name, testsuite.Passed)
			})
		case "rdft:TestTrigEval":
			t.Run(e.Name, func(t *testing.T) {
				if err != nil {
					report.AddTest(e.Name, testsuite.Failed)
					t.Fatal(err)
				}

				// fmt.Stringer
				doc2, err := trig.ParseDocument(doc.String())
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
				ntr, err := nq.ParseDocument(string(raw))
				if _ = ntr; err != nil {
					t.Fatal(err)
				}
				ntr2, err := trig.EvaluateDocument(doc)
				if err != nil {
					report.AddTest(e.Name, testsuite.Failed)
					t.Fatal(err)
				}
				if len(ntr) != len(ntr2) {
					report.AddTest(e.Name, testsuite.Failed)
					t.Fatal(len(ntr), len(ntr2))
				}
				if _, err := nq.ParseDocument(ntr2.String()); err != nil {
					report.AddTest(e.Name, testsuite.Failed)
					t.Fatal(ntr2.String())
				}
				report.AddTest(e.Name, testsuite.Passed)
			})
		case "rdft:TestTrigNegativeEval":
			t.Run(e.Name, func(t *testing.T) {
				if err != nil {
					// If we can not parse the document, we can not evaluate it...
					report.AddTest(e.Name, testsuite.Passed)
					return
				}

				// fmt.Stringer
				doc2, err := trig.ParseDocument(doc.String())
				if err != nil {
					report.AddTest(e.Name, testsuite.Failed)
					t.Fatal(err)
				}
				if len(doc) != len(doc2) {
					report.AddTest(e.Name, testsuite.Failed)
					t.Fatal(len(doc), len(doc2))
				}

				if _, err := trig.EvaluateDocument(doc); err == nil {
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
