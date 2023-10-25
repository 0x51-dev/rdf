package nquads

import (
	"fmt"
	nq "github.com/0x51-dev/rdf/nquads"
	nt "github.com/0x51-dev/rdf/ntriples"
	"github.com/0x51-dev/rdf/star/nquads/grammar"
	nts "github.com/0x51-dev/rdf/star/ntriples"
	"github.com/0x51-dev/upeg/parser"
	"github.com/0x51-dev/upeg/parser/op"
	"strings"
)

type Document []Quad

func ParseDocument(doc string) (Document, error) {
	if len(doc) == 0 {
		return nil, nil
	}
	if !strings.HasSuffix(doc, "\n") {
		doc += "\n"
	}
	p, err := grammar.NewParser([]rune(doc))
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
		quad, err := ParseQuad(n)
		if err != nil {
			return nil, err
		}
		quads = append(quads, *quad)
	}
	return quads, nil
}

type Quad struct {
	nts.Triple
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
	s, err := nts.ParseSubject(children[0])
	if err != nil {
		return nil, err
	}
	p, err := nt.ParsePredicate(children[1])
	if err != nil {
		return nil, err
	}
	o, err := nts.ParseObject(children[2])
	if err != nil {
		return nil, err
	}
	var g nt.Subject
	if len(children) == 4 {
		g, err = nq.ParseGraphLabel(children[3])
		if err != nil {
			return nil, err
		}
	}
	return &Quad{
		Triple: nts.Triple{
			Subject:   s,
			Predicate: *p,
			Object:    o,
		},
		GraphLabel: g,
	}, nil
}
