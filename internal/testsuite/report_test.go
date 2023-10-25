package testsuite_test

import (
	"fmt"
	"github.com/0x51-dev/rdf/internal/testsuite"
	ttl "github.com/0x51-dev/rdf/turtle"
	"time"
)

func ExampleNewReport() {
	r := testsuite.NewReport()
	r.AddTestCase(testsuite.TestCase{
		AssertedBy: ttl.IRI{Value: "https://github.com/q-uint"},
		Mode:       testsuite.Automatic,
		Result: testsuite.TestResult{
			Date: ttl.StringLiteral{
				Value:       "2023-09-09+00:00",
				DatatypeIRI: &ttl.IRI{Prefixed: true, Value: "xsd:date"},
			},
			Outcome: testsuite.Passed,
		},
		Subject: ttl.IRI{Value: "https://github.com/0x51-dev/rdf"},
		Test:    ttl.IRI{Value: "http://example.com/test"},
	})
	r.Project = testsuite.Project{
		IRI:                 ttl.IRI{Value: "https://github.com/0x51-dev/rdf"},
		Name:                "RDF",
		Homepage:            "https://github.com/0x51-dev/rdf",
		License:             "https://www.apache.org/licenses/LICENSE-2.0",
		Description:         "RDF is a Go library for working with RDF data.",
		Created:             time.Date(2023, 7, 15, 0, 0, 0, 0, time.UTC),
		ProgrammingLanguage: "Go",
		Implements: []string{
			"https://www.w3.org/TR/turtle/",
		},
		Developer: []testsuite.Developer{
			{
				IRI:      ttl.IRI{Value: "https://github.com/q-uint"},
				Name:     "Quint Daenen",
				Title:    "Implementor",
				MBox:     "mailto:quint@0x51.dev",
				Homepage: "https://0x51.dev",
			},
		},
	}
	fmt.Println(r)
	// Output:
	// @prefix dc: <http://purl.org/dc/elements/1.1/> .
	// @prefix rdft: <http://www.w3.org/ns/rdftest#> .
	// @prefix earl: <http://www.w3.org/ns/earl#> .
	// @prefix foaf: <http://xmlns.com/foaf/0.1/> .
	// @prefix turtletest: <http://www.w3.org/2013/TurtleTests/manifest.ttl#> .
	// @prefix dct: <http://purl.org/dc/terms/> .
	// @prefix xsd: <http://www.w3.org/2001/XMLSchema#> .
	// @prefix rdf: <http://www.w3.org/1999/02/22-rdf-syntax-ns#> .
	// @prefix doap: <http://usefulinc.com/ns/doap#> .
	// <https://github.com/q-uint> a foaf:Person, earl:Assertor ; foaf:name "Quint Daenen" ; foaf:title "Implementor" ; foaf:mbox <mailto:quint@0x51.dev> ; foaf:homepage <https://0x51.dev> .
	// <https://github.com/0x51-dev/rdf> a doap:Project ; doap:name "RDF" ; doap:homepage <https://github.com/0x51-dev/rdf> ; doap:license <https://www.apache.org/licenses/LICENSE-2.0> ; doap:description "RDF is a Go library for working with RDF data."@en ; doap:created "2023-07-15+0000"^^xsd:date ; doap:programming-language <Go> ; doap:implements <https://www.w3.org/TR/turtle/> ; doap:developer <https://github.com/q-uint> .
	// [ a earl:Assertion ; earl:assertedBy <https://github.com/q-uint> ; earl:mode earl:automatic ; earl:result [ a earl:TestResult ; dct:date "2023-09-09+00:00"^^xsd:date ; earl:outcome earl:passed ] ; earl:subject <https://github.com/0x51-dev/rdf> ; earl:test <http://example.com/test> ] .
}
