package ntriples_test

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
	doc, _ := nt.ParseDocument(example1)
	fmt.Println(doc)
	// Output:
	// <http://one.example/subject1> <http://one.example/predicate1> <http://one.example/object1> .
	// _:subject1 <http://an.example/predicate1> "object1" .
	// _:subject2 <http://an.example/predicate2> "object2" .
}

func TestDocument_Equal(t *testing.T) {
	a := nt.IRIReference("https://example.com/a")
	b := nt.IRIReference("https://example.com/b")
	d := nt.Document{
		nt.Triple{
			Subject:   a,
			Predicate: b,
			Object:    nt.BlankNode("b1"),
		},
	}
	if !d.Equal(d) {
		t.Error()
	}
	if !d.Equal(nt.Document{
		nt.Triple{
			Subject:   a,
			Predicate: b,
			Object:    nt.BlankNode("b2"),
		},
	}) {
		t.Error()
	}
	t.Run("blank node", func(t *testing.T) {
		t12 := nt.Triple{
			Subject:   nt.BlankNode("b1"),
			Predicate: b,
			Object:    nt.BlankNode("b2"),
		}
		t21 := nt.Triple{
			Subject:   nt.BlankNode("b2"),
			Predicate: b,
			Object:    nt.BlankNode("b1"),
		}
		if (nt.Document{t12, t21}).Equal(nt.Document{t12}) {
			t.Error()
		}
		if (nt.Document{t12, t21}).Equal(nt.Document{t12, t12}) {
			t.Error()
		}
		if !(nt.Document{t12, t21}).Equal(nt.Document{t21, t12}) {
			t.Error()
		}
		if (nt.Document{t12, t21}).Equal(nt.Document{t21, nt.Triple{
			Subject:   nt.BlankNode("b2"),
			Predicate: b,
			Object:    nt.BlankNode("b2"), // Different reference.
		}}) {
			t.Error()
		}
	})
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
		doc, err := nt.ParseDocument(test.doc)
		if err != nil {
			t.Fatal(err)
		}
		if len(doc) != test.triples {
			t.Error(len(doc))
		}

		{ // fmt.Stringer
			doc2, err := nt.ParseDocument(doc.String())
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

	report := project.NewReport(ttl.IRI{Value: "http://www.w3.org/2013/N-TriplesTests/manifest.ttl#"})
	for _, k := range manifest.Keys {
		e := manifest.Entries[k]
		raw, err := suite.ReadFile(fmt.Sprintf("testdata/suite/%s", e.Action))
		if err != nil {
			t.Fatal(err)
		}
		doc, err := nt.ParseDocument(string(raw))
		switch e.Type {
		case "rdft:TestNTriplesPositiveSyntax":
			t.Run(e.Name, func(t *testing.T) {
				if err != nil {
					report.AddTest(e.Name, testsuite.Failed)
					t.Fatal(err)
				}

				// fmt.Stringer
				doc2, err := nt.ParseDocument(doc.String())
				if err != nil {
					report.AddTest(e.Name, testsuite.Failed)
					t.Fatal(err)
				}
				if !doc.Equal(doc2) {
					report.AddTest(e.Name, testsuite.Failed)
					t.Fatal(doc, doc2)
				}

				report.AddTest(e.Name, testsuite.Passed)
			})
		case "rdft:TestNTriplesNegativeSyntax":
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
