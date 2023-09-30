package ntriples

import (
	"fmt"
	"github.com/0x51-dev/rdf/ntriples/grammar"
	"github.com/0x51-dev/rids/iri"
	"github.com/0x51-dev/upeg/parser"
	"github.com/0x51-dev/upeg/parser/op"
	"strconv"
	"strings"
)

var validation = true

// ToggleValidation enables/disables validation of IRIs.
func ToggleValidation(enabled bool) {
	validation = enabled
}

type BlankNode string

func parseBlankNodeLabel(n *parser.Node) (*BlankNode, error) {
	if n.Name != "BlankNodeLabel" {
		return nil, fmt.Errorf("blank-node: unknown %s", n.Name)
	}
	bn := BlankNode(n.Value())
	return &bn, nil
}

func (n *BlankNode) String() string {
	return string(*n)
}

func (n *BlankNode) object() {}

func (n *BlankNode) subject() {}

type Document []Triple

func ParseDocument(doc string) (Document, error) {
	if len(doc) == 0 {
		return nil, nil
	}
	if !strings.HasSuffix(doc, "\n") {
		doc += "\n"
	}
	p, err := parser.New([]rune(doc))
	if err != nil {
		return nil, err
	}
	n, err := p.Parse(op.And{grammar.Document, op.EOF{}})
	if err != nil {
		return nil, err
	}
	return parseDocument(n)
}

func parseDocument(n *parser.Node) (Document, error) {
	if n.Name != "Document" {
		return nil, fmt.Errorf("document: unknown %s", n.Name)
	}
	var triples []Triple
	for _, n := range n.Children() {
		t, err := parseTriple(n)
		if err != nil {
			return nil, err
		}
		triples = append(triples, *t)
	}
	return triples, nil
}

func (d Document) String() string {
	var s string
	for _, t := range d {
		s += fmt.Sprintf("%s\n", t)
	}
	return s
}

type IRIReference string

func parseIRIReference(n *parser.Node) (*IRIReference, error) {
	if n.Name != "IRIReference" {
		return nil, fmt.Errorf("iri-reference: unknown %s", n.Name)
	}

	ref := IRIReference(n.Value())
	// IRIs in the RDF abstract syntax must be absolute, and may contain a fragment identifier.
	if validation && !ref.IsValid() {
		return nil, fmt.Errorf("iri-reference: invalid: %s", ref)
	}
	return &ref, nil
}

func parsePredicate(n *parser.Node) (*IRIReference, error) {
	if n.Name != "Predicate" {
		return nil, fmt.Errorf("predicate: unknown %s", n.Name)
	}
	return parseIRIReference(n.Children()[0])
}

func (r IRIReference) IsValid() bool {
	v := string(r)
	if strings.Contains(v, "\\u") || strings.Contains(v, "\\U") {
		// Unescape unicode characters.
		if v_, err := strconv.Unquote(`"` + v + `"`); err == nil {
			v = v_
		}
	}
	p, err := parser.New([]rune(v))
	if err != nil {
		return false
	}
	if _, err := p.Match(op.And{iri.AbsoluteIRI, op.Optional{Value: op.And{'#', iri.Ifragment}}, op.EOF{}}); err != nil {
		return false
	}
	return true
}

func (r IRIReference) String() string {
	return fmt.Sprintf("<%s>", string(r))
}

func (r IRIReference) object() {}

func (r IRIReference) subject() {}

type Literal struct {
	Value     string
	Reference *IRIReference
	Language  string
}

func parseLiteral(n *parser.Node) (*Literal, error) {
	if n.Name != "Literal" {
		return nil, fmt.Errorf("literal: unknown %s", n.Name)
	}
	var literal Literal
	for _, n := range n.Children() {
		switch n.Name {
		case "StringLiteral":
			literal.Value = n.Value()
		case "IRIReference":
			ref, err := parseIRIReference(n)
			if err != nil {
				return nil, err
			}
			literal.Reference = ref
		case "LanguageTag":
			literal.Language = n.Value()
		default:
			return nil, fmt.Errorf("literal: unknown child: %s", n.Name)
		}
	}
	return &literal, nil
}

func (l Literal) String() string {
	if l.Reference != nil {
		return fmt.Sprintf(`"%s"^^%s`, l.Value, l.Reference)
	}
	if len(l.Language) > 0 {
		return fmt.Sprintf(`"%s"@%s`, l.Value, l.Language)
	}
	return fmt.Sprintf(`"%s"`, l.Value)
}

func (l Literal) object() {}

type Object interface {
	object()

	fmt.Stringer
}

func parseObject(n *parser.Node) (Object, error) {
	if n.Name != "Object" {
		return nil, fmt.Errorf("object: unknown: %s", n.Name)
	}
	switch n = n.Children()[0]; n.Name {
	case "IRIReference":
		return parseIRIReference(n)
	case "BlankNodeLabel":
		return parseBlankNodeLabel(n)
	case "Literal":
		return parseLiteral(n)
	default:
		return nil, fmt.Errorf("object: unknown: %s", n.Name)
	}
}

type Subject interface {
	subject()
}

func parseSubject(n *parser.Node) (Subject, error) {
	if n.Name != "Subject" {
		return nil, fmt.Errorf("subject: unknown: %s", n.Name)
	}
	switch n = n.Children()[0]; n.Name {
	case "IRIReference":
		return parseIRIReference(n)
	case "BlankNodeLabel":
		return parseBlankNodeLabel(n)
	default:
		return nil, fmt.Errorf("subject: unknown: %s", n.Name)
	}
}

type Triple struct {
	Subject   Subject
	Predicate IRIReference
	Object    Object
}

func parseTriple(n *parser.Node) (*Triple, error) {
	if n.Name != "Triple" {
		return nil, fmt.Errorf("triple: unknown %s", n.Name)
	}
	if len(n.Children()) != 3 {
		return nil, fmt.Errorf("triple: expected 3 children")
	}
	children := n.Children()
	s, err := parseSubject(children[0])
	if err != nil {
		return nil, err
	}
	p, err := parsePredicate(children[1])
	if err != nil {
		return nil, err
	}
	o, err := parseObject(children[2])
	if err != nil {
		return nil, err
	}
	return &Triple{
		Subject:   s,
		Predicate: *p,
		Object:    o,
	}, nil
}

func (t Triple) String() string {
	return fmt.Sprintf("%s %s %s .", t.Subject, t.Predicate, t.Object)
}
