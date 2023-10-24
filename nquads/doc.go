package nquads

import (
	"fmt"
	"github.com/0x51-dev/rdf/nquads/grammar"
	nt "github.com/0x51-dev/rdf/ntriples"
	"github.com/0x51-dev/upeg/parser"
	"github.com/0x51-dev/upeg/parser/op"
	"strings"
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
		quad, err := ParseQuad(n)
		if err != nil {
			return nil, err
		}
		quads = append(quads, *quad)
	}
	return quads, nil
}

func (d Document) Equal(other Document) bool {
	if len(d) != len(other) {
		return false
	}
	g0, g1 := d.Graphs(), other.Graphs()
	for g, q0 := range g0 {
		q1, ok := g1[g]
		if !ok {
			return false
		}
		if !q0.Equal(q1) {
			return false
		}
	}
	return true
}

func (d Document) Graphs() map[nt.Subject]nt.Document {
	g := make(map[nt.Subject]nt.Document)
	for _, q := range d {
		g[q.GraphLabel] = append(g[q.GraphLabel], q.Triple)
	}
	return g
}

func (d Document) String() string {
	var s string
	for _, q := range d {
		s += fmt.Sprintf("%s\n", q)
	}
	return s
}

type Quad struct {
	nt.Triple
	GraphLabel nt.Subject
}

func NewQuadFromTriple(t nt.Triple, graphLabel nt.Subject) Quad {
	return Quad{
		Triple:     t,
		GraphLabel: graphLabel,
	}
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
		Triple: nt.Triple{
			Subject:   s,
			Predicate: *p,
			Object:    o,
		},
		GraphLabel: g,
	}, nil
}

func (q Quad) String() string {
	if q.GraphLabel == nil {
		return fmt.Sprintf("%s %s %s .", q.Subject, q.Predicate, q.Object)
	}
	return fmt.Sprintf("%s %s %s %s .", q.Subject, q.Predicate, q.Object, q.GraphLabel)
}
