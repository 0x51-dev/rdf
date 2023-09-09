package project

import (
	"fmt"
	"github.com/0x51-dev/rdf/internal/testsuite"
	"github.com/0x51-dev/rdf/turtle"
	"time"
)

var (
	assertedBy = turtle.IRI{Value: "https://github.com/q-uint"}
	subject    = turtle.IRI{Value: "https://github.com/0x51-dev/rdf"}
)

type Report struct {
	test turtle.IRI
	r    *testsuite.Report
}

func NewReport(test turtle.IRI) *Report {
	r := testsuite.NewReport()
	r.Project = testsuite.Project{
		IRI:                 turtle.IRI{Value: "https://github.com/0x51-dev/rdf"},
		Name:                "RDF",
		Homepage:            "https://github.com/0x51-dev/rdf",
		License:             "https://www.apache.org/licenses/LICENSE-2.0",
		Description:         "RDF is a Go library for working with RDF data.",
		Created:             time.Date(2023, 7, 15, 0, 0, 0, 0, time.UTC),
		ProgrammingLanguage: "Go",
		Implements: []string{
			"https://www.w3.org/TR/n-triples/",
			"https://www.w3.org/TR/n-quads/",
			"https://www.w3.org/TR/turtle/",
		},
		Developer: []testsuite.Developer{
			{
				IRI:      turtle.IRI{Value: "https://github.com/q-uint"},
				Name:     "Quint Daenen",
				Title:    "Implementor",
				MBox:     "mailto:quint@0x51.dev",
				Homepage: "https://0x51.dev",
			},
		},
	}
	return &Report{test: test, r: r}
}

func (r *Report) AddTest(name string, outcome testsuite.OutcomeValue) {
	test := turtle.IRI{Prefixed: r.test.Prefixed, Value: fmt.Sprintf("%s%s", r.test.Value, name)}
	r.r.AddTestCase(testsuite.TestCase{
		AssertedBy: assertedBy,
		Mode:       testsuite.Automatic,
		Result: testsuite.TestResult{
			Date: turtle.StringLiteral{
				Value:       time.Now().In(time.UTC).Format("2006-01-02-0700"),
				DatatypeIRI: "xsd:date",
			},
			Outcome: outcome,
		},
		Subject: subject,
		Test:    test,
	})
}

func (r *Report) String() string {
	return r.r.String()
}
