package ir

import (
	"fmt"
	nt "github.com/0x51-dev/rdf/ntriples/grammar/ir"
	"github.com/0x51-dev/upeg/parser"
)

func ParseGraphLabel(n *parser.Node) (nt.Subject, error) {
	if n.Name != "GraphLabel" {
		return nil, fmt.Errorf("subject: unknown: %s", n.Name)
	}
	switch n = n.Children()[0]; n.Name {
	case "IRIReference":
		return nt.ParseIRIReference(n)
	case "BlankNodeLabel":
		return nt.ParseBlankNodeLabel(n)
	default:
		return nil, fmt.Errorf("subject: unknown: %s", n.Name)
	}
}

type Document []Quad

func ParseDocument(n *parser.Node) (Document, error) {
	if n.Name != "Document" {
		return nil, fmt.Errorf("document: unknown %s", n.Name)
	}
	var quads []Quad
	for _, n := range n.Children() {
		quad, err := ParseQuad(n)
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

func ParseQuad(n *parser.Node) (*Quad, error) {
	if n.Name != "Statement" {
		return nil, fmt.Errorf("quad: unknown %s", n.Name)
	}
	if l := len(n.Children()); l != 3 && l != 4 {
		return nil, fmt.Errorf("quad: expected 3 or 4 children")
	}
	children := n.Children()
	s, err := nt.ParseSubject(children[0])
	if err != nil {
		return nil, err
	}
	p, err := nt.ParsePredicate(children[1])
	if err != nil {
		return nil, err
	}
	o, err := nt.ParseObject(children[2])
	if err != nil {
		return nil, err
	}
	var g nt.Subject
	if len(children) == 4 {
		g, err = ParseGraphLabel(children[3])
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
