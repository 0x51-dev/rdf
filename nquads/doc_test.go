package nquads_test

import (
	"embed"
	_ "embed"
	"fmt"
	"github.com/0x51-dev/rdf/internal/project"
	"github.com/0x51-dev/rdf/internal/testsuite"
	nq "github.com/0x51-dev/rdf/nquads"
	nt "github.com/0x51-dev/rdf/ntriples"
	ttl "github.com/0x51-dev/rdf/turtle"
	"os"
	"testing"
)

var (
	//go:embed testdata/example1.nq
	example1 string

	//go:embed testdata/example2.nq
	example2 string

	//go:embed testdata/example3.nq
	example3 string

	//go:embed testdata/suite/manifest.ttl
	rawManifest string

	//go:embed testdata/suite/*.nq
	suite embed.FS
)

func Example_example1() {
	doc, _ := nq.ParseDocument(example1)
	fmt.Println(doc)
	// Output:
	// <http://one.example/subject1> <http://one.example/predicate1> <http://one.example/object1> <http://example.org/graph3> .
	// _:subject1 <http://an.example/predicate1> "object1" <http://example.org/graph1> .
	// _:subject2 <http://an.example/predicate2> "object2" <http://example.org/graph5> .
}

func TestDocument_Equal(t *testing.T) {
	a := nt.IRIReference("https://example.com/a")
	b := nt.IRIReference("https://example.com/b")
	d := nq.Document{
		nq.NewQuadFromTriple(
			nt.Triple{
				Subject:   a,
				Predicate: b,
				Object:    nt.BlankNode("b1"),
			}, nil,
		),
	}
	if !d.Equal(d) {
		t.Error()
	}
	if !d.Equal(nq.Document{
		nq.NewQuadFromTriple(
			nt.Triple{
				Subject:   a,
				Predicate: b,
				Object:    nt.BlankNode("b2"),
			}, nil,
		),
	}) {
		t.Error()
	}
	if d.Equal(nq.Document{
		nq.NewQuadFromTriple(
			nt.Triple{
				Subject:   a,
				Predicate: b,
				Object:    nt.BlankNode("b1"),
			}, nt.IRIReference("https://example.com/g"), // Different graph.
		),
	}) {
		t.Error()
	}
}

func TestExamples(t *testing.T) {
	for _, test := range []struct {
		doc   string
		quads int
	}{
		{example1, 3},
		{example2, 1},
		{example3, 2},
	} {
		doc, err := nq.ParseDocument(test.doc)
		if err != nil {
			t.Fatal(err)
		}
		if len(doc) != test.quads {
			t.Error(len(doc))
		}

		{ // fmt.Stringer
			doc2, err := nq.ParseDocument(doc.String())
			if err != nil {
				t.Fatal(err)
			}
			if !doc.Equal(doc2) {
				t.Error(doc, doc2)
			}
		}
	}
}

func TestSuite(t *testing.T) {
	manifest, err := testsuite.LoadManifest(rawManifest)
	if err != nil {
		t.Fatal(err)
	}

	report := project.NewReport(ttl.IRI{Value: "http://www.w3.org/2013/N-QuadsTests/manifest.ttl#"})
	for _, k := range manifest.Keys {
		e := manifest.Entries[k]
		raw, err := suite.ReadFile(fmt.Sprintf("testdata/suite/%s", e.Action))
		if err != nil {
			t.Fatal(err)
		}
		doc, err := nq.ParseDocument(string(raw))
		switch e.Type {
		case "rdft:TestNQuadsPositiveSyntax":
			t.Run(e.Name, func(t *testing.T) {
				if err != nil {
					report.AddTest(e.Name, testsuite.Failed)
					t.Fatal(err)
				}

				// fmt.Stringer
				doc2, err := nq.ParseDocument(doc.String())
				if err != nil {
					report.AddTest(e.Name, testsuite.Failed)
					t.Fatal(err)
				}
				if !doc.Equal(doc2) {
					report.AddTest(e.Name, testsuite.Failed)
					t.Fatal(len(doc), len(doc2))
				}

				report.AddTest(e.Name, testsuite.Passed)
			})
		case "rdft:TestNQuadsNegativeSyntax":
			t.Run(e.Name, func(t *testing.T) {
				if err == nil {
					report.AddTest(e.Name, testsuite.Failed)
					t.Fatal("expected error")
				}

				report.AddTest(e.Name, testsuite.Passed)
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
