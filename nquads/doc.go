package nquads

import (
	"fmt"
	"github.com/0x51-dev/rdf/nquads/grammar"
	nt "github.com/0x51-dev/rdf/ntriples"
	"github.com/0x51-dev/upeg/parser"
	"github.com/0x51-dev/upeg/parser/op"
	"strings"
)

var validation = true

// DisableValidation disables validation of IRIs.
func DisableValidation() {
	validation = false
}

func parseBlankNodeLabel(n *parser.Node) (*nt.BlankNode, error) {
	if n.Name != "BlankNodeLabel" {
		return nil, fmt.Errorf("blank-node: unknown %s", n.Name)
	}
	bn := nt.BlankNode(n.Value())
	return &bn, nil
}

func parseGraphLabel(n *parser.Node) (nt.Subject, error) {
	if n.Name != "GraphLabel" {
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

func parseIRIReference(n *parser.Node) (*nt.IRIReference, error) {
	if n.Name != "IRIReference" {
		return nil, fmt.Errorf("iri-reference: unknown %s", n.Name)
	}

	ref := nt.IRIReference(n.Value())
	// IRIs in the RDF abstract syntax must be absolute, and may contain a fragment identifier.
	if validation && !ref.IsValid() {
		return nil, fmt.Errorf("iri-reference: invalid: %s", ref)
	}
	return &ref, nil
}

func parseLiteral(n *parser.Node) (*nt.Literal, error) {
	if n.Name != "Literal" {
		return nil, fmt.Errorf("literal: unknown %s", n.Name)
	}
	var literal nt.Literal
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

func parseObject(n *parser.Node) (nt.Object, error) {
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

func parsePredicate(n *parser.Node) (*nt.IRIReference, error) {
	if n.Name != "Predicate" {
		return nil, fmt.Errorf("predicate: unknown %s", n.Name)
	}
	return parseIRIReference(n.Children()[0])
}

func parseSubject(n *parser.Node) (nt.Subject, error) {
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

type Document []Quad

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
	var quads []Quad
	for _, n := range n.Children() {
		quad, err := parseQuad(n)
		if err != nil {
			return nil, err
		}
		quads = append(quads, *quad)
	}
	return quads, nil
}

func (d Document) String() string {
	var s string
	for _, q := range d {
		s += fmt.Sprintf("%s\n", q)
	}
	return s
}

type Quad struct {
	Subject    nt.Subject
	Predicate  nt.IRIReference
	Object     nt.Object
	GraphLabel nt.Subject
}

func parseQuad(n *parser.Node) (*Quad, error) {
	if n.Name != "Statement" {
		return nil, fmt.Errorf("quad: unknown %s", n.Name)
	}
	if l := len(n.Children()); l != 3 && l != 4 {
		return nil, fmt.Errorf("quad: expected 3 or 4 children")
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
	var g nt.Subject
	if len(children) == 4 {
		g, err = parseGraphLabel(children[3])
		if err != nil {
			return nil, err
		}
	}
	return &Quad{
		Subject:    s,
		Predicate:  *p,
		Object:     o,
		GraphLabel: g,
	}, nil
}

func (q Quad) String() string {
	if q.GraphLabel == nil {
		return fmt.Sprintf("%s %s %s .", q.Subject, q.Predicate, q.Object)
	}
	return fmt.Sprintf("%s %s %s %s .", q.Subject, q.Predicate, q.Object, q.GraphLabel)
}
