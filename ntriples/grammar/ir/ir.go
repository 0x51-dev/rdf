package ir

import (
	"fmt"
	"github.com/0x51-dev/upeg/parser"
)

type BlankNode string

func (n BlankNode) object() {}

func (n BlankNode) subject() {}

type IRIReference string

func ParseIRIReference(n *parser.Node) (*IRIReference, error) {
	if n.Name != "IRIReference" {
		return nil, fmt.Errorf("iri-reference: unknown %s", n.Name)
	}
	ref := IRIReference(n.Value())
	return &ref, nil
}

func (r IRIReference) object() {}

func (r IRIReference) subject() {}

type Literal string

func (l Literal) object() {}

type Object interface {
	object()
}

func ParseObject(n *parser.Node) (Object, error) {
	switch n.Name {
	case "IRIReference":
		ref := IRIReference(n.Value())
		return &ref, nil
	case "BlankNodeLabel":
		bn := BlankNode(n.Value())
		return &bn, nil
	case "Literal":
		l := Literal(n.Value())
		return &l, nil
	default:
		return nil, fmt.Errorf("object: unknown: %s", n.Name)
	}
}

type Subject interface {
	subject()
}

func ParseSubject(n *parser.Node) (Subject, error) {
	switch n.Name {
	case "IRIReference":
		ref := IRIReference(n.Value())
		return &ref, nil
	case "BlankNodeLabel":
		bn := BlankNode(n.Value())
		return &bn, nil
	default:
		return nil, fmt.Errorf("subject: unknown: %s", n.Name)
	}
}

type Triple struct {
	Subject   Subject
	Predicate IRIReference
	Object    Object
}

func ParseDocument(n *parser.Node) ([]Triple, error) {
	var triples []Triple
	for _, n := range n.Children() {
		t, err := ParseTriple(n)
		if err != nil {
			return nil, err
		}
		triples = append(triples, *t)
	}
	return triples, nil
}

func ParseTriple(n *parser.Node) (*Triple, error) {
	if len(n.Children()) != 3 {
		return nil, fmt.Errorf("triple: expected 3 children")
	}
	children := n.Children()
	s, err := ParseSubject(children[0])
	if err != nil {
		return nil, err
	}
	p, err := ParseIRIReference(children[1])
	if err != nil {
		return nil, err
	}
	o, err := ParseObject(children[2])
	if err != nil {

	}
	return &Triple{
		Subject:   s,
		Predicate: *p,
		Object:    o,
	}, nil
}
