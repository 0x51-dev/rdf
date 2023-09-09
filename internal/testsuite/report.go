package testsuite

import (
	"github.com/0x51-dev/rdf/turtle"
	"time"
)

type Developer struct {
	IRI      turtle.IRI
	Name     string
	Title    string
	MBox     string
	Homepage string
}

func (d Developer) line() turtle.Line {
	var b turtle.Triple
	b.Subject = d.IRI
	b.PredicateObjectList = turtle.PredicateObjectList{
		{
			Verb: new(turtle.A),
			ObjectList: []turtle.Object{
				&turtle.IRI{Prefixed: true, Value: "foaf:Person"},
				&turtle.IRI{Prefixed: true, Value: "earl:Assertor"},
			},
		},
		{
			Verb:       &turtle.IRI{Prefixed: true, Value: "foaf:name"},
			ObjectList: []turtle.Object{&turtle.StringLiteral{Value: d.Name}},
		},
		{
			Verb:       &turtle.IRI{Prefixed: true, Value: "foaf:title"},
			ObjectList: []turtle.Object{&turtle.StringLiteral{Value: d.Title}},
		},
		{
			Verb:       &turtle.IRI{Prefixed: true, Value: "foaf:mbox"},
			ObjectList: []turtle.Object{&turtle.IRI{Value: d.MBox}},
		},
		{
			Verb:       &turtle.IRI{Prefixed: true, Value: "foaf:homepage"},
			ObjectList: []turtle.Object{&turtle.IRI{Value: d.Homepage}},
		},
	}
	return &b
}

// OutcomeValue describes a resulting condition from carrying out the test.
type OutcomeValue string

const (
	// Passed - the subject passed the test.
	Passed OutcomeValue = "earl:passed"
	// Failed - the subject failed the test.
	Failed OutcomeValue = "earl:failed"
	// CantTell - it is unclear if the subject passed or failed the test.
	CantTell OutcomeValue = "earl:cantTell"
	// Inapplicable - the test is not applicable to the subject.
	Inapplicable OutcomeValue = "earl:inapplicable"
	// Untested - the test has not been carried out.
	Untested OutcomeValue = "earl:untested"
)

type Project struct {
	IRI                 turtle.IRI
	Name                string
	Homepage            string
	License             string
	Description         string
	Created             time.Time
	ProgrammingLanguage string
	Implements          []string
	Developer           []Developer
}

func (p Project) line() turtle.Line {
	var t turtle.Triple
	t.Subject = p.IRI
	var implements []turtle.Object
	for _, i := range p.Implements {
		implements = append(implements, &turtle.IRI{Value: i})
	}
	var developers []turtle.Object
	for _, d := range p.Developer {
		developers = append(developers, &d.IRI)
	}
	t.PredicateObjectList = []turtle.PredicateObject{
		{
			Verb:       new(turtle.A),
			ObjectList: []turtle.Object{&turtle.IRI{Prefixed: true, Value: "doap:Project"}},
		},
		{
			Verb:       &turtle.IRI{Prefixed: true, Value: "doap:name"},
			ObjectList: []turtle.Object{&turtle.StringLiteral{Value: p.Name}},
		},
		{
			Verb:       &turtle.IRI{Prefixed: true, Value: "doap:homepage"},
			ObjectList: []turtle.Object{&turtle.IRI{Value: p.Homepage}},
		},
		{
			Verb:       &turtle.IRI{Prefixed: true, Value: "doap:license"},
			ObjectList: []turtle.Object{&turtle.IRI{Value: p.License}},
		},
		{
			Verb:       &turtle.IRI{Prefixed: true, Value: "doap:description"},
			ObjectList: []turtle.Object{&turtle.StringLiteral{Value: p.Description, LanguageTag: "en"}},
		},
		{
			Verb:       &turtle.IRI{Prefixed: true, Value: "doap:created"},
			ObjectList: []turtle.Object{&turtle.StringLiteral{Value: p.Created.Format("2006-01-02-0700"), DatatypeIRI: "xsd:date"}},
		},
		{
			Verb:       &turtle.IRI{Prefixed: true, Value: "doap:programming-language"},
			ObjectList: []turtle.Object{&turtle.IRI{Value: p.ProgrammingLanguage}},
		},
		{
			Verb:       &turtle.IRI{Prefixed: true, Value: "doap:implements"},
			ObjectList: implements,
		},
		{
			Verb:       &turtle.IRI{Prefixed: true, Value: "doap:developer"},
			ObjectList: developers,
		},
	}
	return t
}

type Report struct {
	prefixes  []turtle.Line
	Project   Project
	testcases []TestCase
}

func NewReport() *Report {
	return &Report{
		prefixes: []turtle.Line{
			&turtle.Prefix{Name: "dc:", IRI: "http://purl.org/dc/elements/1.1/"},
			&turtle.Prefix{Name: "rdft:", IRI: "http://www.w3.org/ns/rdftest#"},
			&turtle.Prefix{Name: "earl:", IRI: "http://www.w3.org/ns/earl#"},
			&turtle.Prefix{Name: "foaf:", IRI: "http://xmlns.com/foaf/0.1/"},
			&turtle.Prefix{Name: "turtletest:", IRI: "http://www.w3.org/2013/TurtleTests/manifest.ttl#"},
			&turtle.Prefix{Name: "dct:", IRI: "http://purl.org/dc/terms/"},
			&turtle.Prefix{Name: "xsd:", IRI: "http://www.w3.org/2001/XMLSchema#"},
			&turtle.Prefix{Name: "rdf:", IRI: "http://www.w3.org/1999/02/22-rdf-syntax-ns#"},
			&turtle.Prefix{Name: "doap:", IRI: "http://usefulinc.com/ns/doap#"},
		},
	}
}

func (r *Report) AddTestCase(tc TestCase) {
	r.testcases = append(r.testcases, tc)
}

func (r *Report) String() string {
	var d turtle.Document
	d = append(d, r.prefixes...)
	for _, dev := range r.Project.Developer {
		d = append(d, dev.line())
	}
	d = append(d, r.Project.line())
	for _, tc := range r.testcases {
		d = append(d, tc.line())
	}
	return d.String()
}

type TestCase struct {
	AssertedBy turtle.IRI
	Mode       TestMode
	Result     TestResult
	Subject    turtle.IRI
	Test       turtle.IRI
}

func (tc TestCase) line() turtle.Line {
	var b turtle.Triple
	b.BlankNodePropertyList = turtle.BlankNodePropertyList{
		{
			Verb:       new(turtle.A),
			ObjectList: []turtle.Object{&turtle.IRI{Prefixed: true, Value: "earl:Assertion"}},
		},
		{
			Verb:       &turtle.IRI{Prefixed: true, Value: "earl:assertedBy"},
			ObjectList: []turtle.Object{&tc.AssertedBy},
		},
		{
			Verb:       &turtle.IRI{Prefixed: true, Value: "earl:mode"},
			ObjectList: []turtle.Object{&turtle.IRI{Prefixed: true, Value: string(tc.Mode)}},
		},
		{
			Verb: &turtle.IRI{Prefixed: true, Value: "earl:result"},
			ObjectList: []turtle.Object{turtle.BlankNodePropertyList{
				{
					Verb:       new(turtle.A),
					ObjectList: []turtle.Object{&turtle.IRI{Prefixed: true, Value: "earl:TestResult"}},
				},
				{
					Verb:       &turtle.IRI{Prefixed: true, Value: "dct:date"},
					ObjectList: []turtle.Object{&tc.Result.Date},
				},
				{
					Verb:       &turtle.IRI{Prefixed: true, Value: "earl:outcome"},
					ObjectList: []turtle.Object{&turtle.IRI{Prefixed: true, Value: string(tc.Result.Outcome)}},
				},
			}},
		},
		{
			Verb:       &turtle.IRI{Prefixed: true, Value: "earl:subject"},
			ObjectList: []turtle.Object{&tc.Subject},
		},
		{
			Verb:       &turtle.IRI{Prefixed: true, Value: "earl:test"},
			ObjectList: []turtle.Object{&tc.Test},
		},
	}
	return &b
}

// TestMode describes how a test was carried out.
type TestMode string

const (
	// Automatic - where the test was carried out automatically by the software tool and without any human intervention.
	Automatic TestMode = "earl:automatic"
	// Manual - where the test was carried out by human evaluators. This includes the case where the evaluators are
	// aided by instructions or guidance provided by software tools, but where the evaluators carried out the actual
	// test procedure.
	Manual TestMode = "earl:manual"
	// SemiAuto - where the test was partially carried out by software tools, but where human input or judgment was
	// still required to decide or help decide the outcome of the test.
	SemiAuto TestMode = "earl:semiAuto"
	// Undisclosed - where the exact testing process is undisclosed.
	Undisclosed TestMode = "earl:undisclosed"
	// UnknownMode - where the testing process is unknown or undetermined.
	UnknownMode TestMode = "earl:unknownMode"
)

type TestResult struct {
	Date    turtle.StringLiteral
	Outcome OutcomeValue
}
