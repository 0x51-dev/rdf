package testsuite

import (
	ttl "github.com/0x51-dev/rdf/turtle"
	"time"
)

type Developer struct {
	IRI      ttl.IRI
	Name     string
	Title    string
	MBox     string
	Homepage string
}

func (d Developer) line() ttl.Statement {
	var b ttl.Triple
	b.Subject = d.IRI
	b.PredicateObjectList = ttl.PredicateObjectList{
		{
			Verb: new(ttl.A),
			ObjectList: []ttl.Object{
				&ttl.IRI{Prefixed: true, Value: "foaf:Person"},
				&ttl.IRI{Prefixed: true, Value: "earl:Assertor"},
			},
		},
		{
			Verb:       &ttl.IRI{Prefixed: true, Value: "foaf:name"},
			ObjectList: []ttl.Object{&ttl.StringLiteral{Value: d.Name}},
		},
		{
			Verb:       &ttl.IRI{Prefixed: true, Value: "foaf:title"},
			ObjectList: []ttl.Object{&ttl.StringLiteral{Value: d.Title}},
		},
		{
			Verb:       &ttl.IRI{Prefixed: true, Value: "foaf:mbox"},
			ObjectList: []ttl.Object{&ttl.IRI{Value: d.MBox}},
		},
		{
			Verb:       &ttl.IRI{Prefixed: true, Value: "foaf:homepage"},
			ObjectList: []ttl.Object{&ttl.IRI{Value: d.Homepage}},
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
	IRI                 ttl.IRI
	Name                string
	Homepage            string
	License             string
	Description         string
	Created             time.Time
	ProgrammingLanguage string
	Implements          []string
	Developer           []Developer
}

func (p Project) line() ttl.Statement {
	var t ttl.Triple
	t.Subject = p.IRI
	var implements []ttl.Object
	for _, i := range p.Implements {
		implements = append(implements, &ttl.IRI{Value: i})
	}
	var developers []ttl.Object
	for _, d := range p.Developer {
		developers = append(developers, &d.IRI)
	}
	t.PredicateObjectList = []ttl.PredicateObject{
		{
			Verb:       new(ttl.A),
			ObjectList: []ttl.Object{&ttl.IRI{Prefixed: true, Value: "doap:Project"}},
		},
		{
			Verb:       &ttl.IRI{Prefixed: true, Value: "doap:name"},
			ObjectList: []ttl.Object{&ttl.StringLiteral{Value: p.Name}},
		},
		{
			Verb:       &ttl.IRI{Prefixed: true, Value: "doap:homepage"},
			ObjectList: []ttl.Object{&ttl.IRI{Value: p.Homepage}},
		},
		{
			Verb:       &ttl.IRI{Prefixed: true, Value: "doap:license"},
			ObjectList: []ttl.Object{&ttl.IRI{Value: p.License}},
		},
		{
			Verb:       &ttl.IRI{Prefixed: true, Value: "doap:description"},
			ObjectList: []ttl.Object{&ttl.StringLiteral{Value: p.Description, LanguageTag: "en"}},
		},
		{
			Verb:       &ttl.IRI{Prefixed: true, Value: "doap:created"},
			ObjectList: []ttl.Object{&ttl.StringLiteral{Value: p.Created.Format("2006-01-02-0700"), DatatypeIRI: "xsd:date"}},
		},
		{
			Verb:       &ttl.IRI{Prefixed: true, Value: "doap:programming-language"},
			ObjectList: []ttl.Object{&ttl.IRI{Value: p.ProgrammingLanguage}},
		},
		{
			Verb:       &ttl.IRI{Prefixed: true, Value: "doap:implements"},
			ObjectList: implements,
		},
		{
			Verb:       &ttl.IRI{Prefixed: true, Value: "doap:developer"},
			ObjectList: developers,
		},
	}
	return t
}

type Report struct {
	prefixes  []ttl.Statement
	Project   Project
	testcases []TestCase
}

func NewReport() *Report {
	return &Report{
		prefixes: []ttl.Statement{
			&ttl.Prefix{Name: "dc:", IRI: "http://purl.org/dc/elements/1.1/"},
			&ttl.Prefix{Name: "rdft:", IRI: "http://www.w3.org/ns/rdftest#"},
			&ttl.Prefix{Name: "earl:", IRI: "http://www.w3.org/ns/earl#"},
			&ttl.Prefix{Name: "foaf:", IRI: "http://xmlns.com/foaf/0.1/"},
			&ttl.Prefix{Name: "turtletest:", IRI: "http://www.w3.org/2013/TurtleTests/manifest.ttl#"},
			&ttl.Prefix{Name: "dct:", IRI: "http://purl.org/dc/terms/"},
			&ttl.Prefix{Name: "xsd:", IRI: "http://www.w3.org/2001/XMLSchema#"},
			&ttl.Prefix{Name: "rdf:", IRI: "http://www.w3.org/1999/02/22-rdf-syntax-ns#"},
			&ttl.Prefix{Name: "doap:", IRI: "http://usefulinc.com/ns/doap#"},
		},
	}
}

func (r *Report) AddTestCase(tc TestCase) {
	r.testcases = append(r.testcases, tc)
}

func (r *Report) Len() int {
	return len(r.testcases)
}

func (r *Report) String() string {
	var d ttl.Document
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
	AssertedBy ttl.IRI
	Mode       TestMode
	Result     TestResult
	Subject    ttl.IRI
	Test       ttl.IRI
}

func (tc TestCase) line() ttl.Statement {
	var b ttl.Triple
	b.BlankNodePropertyList = ttl.BlankNodePropertyList{
		{
			Verb:       new(ttl.A),
			ObjectList: []ttl.Object{&ttl.IRI{Prefixed: true, Value: "earl:Assertion"}},
		},
		{
			Verb:       &ttl.IRI{Prefixed: true, Value: "earl:assertedBy"},
			ObjectList: []ttl.Object{&tc.AssertedBy},
		},
		{
			Verb:       &ttl.IRI{Prefixed: true, Value: "earl:mode"},
			ObjectList: []ttl.Object{&ttl.IRI{Prefixed: true, Value: string(tc.Mode)}},
		},
		{
			Verb: &ttl.IRI{Prefixed: true, Value: "earl:result"},
			ObjectList: []ttl.Object{ttl.BlankNodePropertyList{
				{
					Verb:       new(ttl.A),
					ObjectList: []ttl.Object{&ttl.IRI{Prefixed: true, Value: "earl:TestResult"}},
				},
				{
					Verb:       &ttl.IRI{Prefixed: true, Value: "dct:date"},
					ObjectList: []ttl.Object{&tc.Result.Date},
				},
				{
					Verb:       &ttl.IRI{Prefixed: true, Value: "earl:outcome"},
					ObjectList: []ttl.Object{&ttl.IRI{Prefixed: true, Value: string(tc.Result.Outcome)}},
				},
			}},
		},
		{
			Verb:       &ttl.IRI{Prefixed: true, Value: "earl:subject"},
			ObjectList: []ttl.Object{&tc.Subject},
		},
		{
			Verb:       &ttl.IRI{Prefixed: true, Value: "earl:test"},
			ObjectList: []ttl.Object{&tc.Test},
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
	Date    ttl.StringLiteral
	Outcome OutcomeValue
}
