package trig

import (
	"fmt"
	nq "github.com/0x51-dev/rdf/nquads"
	"github.com/0x51-dev/rdf/trig/grammar"
	ttl "github.com/0x51-dev/rdf/turtle"
	"github.com/0x51-dev/upeg/parser"
	"github.com/0x51-dev/upeg/parser/op"
	"strings"
)

func EvaluateDocument(doc Document) (nq.Document, error) {
	return NewContext().evaluateDocument(doc)
}

func ValidateDocument(doc Document) bool {
	return NewContext().validateDocument(doc)
}

type Base ttl.Base

func (b Base) String() string {
	return ttl.Base(b).String()
}

func (b Base) statement() {}

type BlankNode ttl.BlankNode

func (b BlankNode) String() string {
	return ttl.BlankNode(b).String()
}

func (b BlankNode) labelOrSubject() {}

type Block interface {
	block()

	Statement

	fmt.Stringer
}

func ParseBlock(n *parser.Node) (Block, error) {
	if n.Name != "Block" {
		return nil, fmt.Errorf("block: unknown: %s", n.Name)
	}
	switch n := n.Children()[0]; n.Name {
	case "TriplesOrGraph":
		return ParseTriplesOrGraph(n)
	case "WrappedGraph":
		return ParseWrappedGraph(n)
	case "Triples2":
		return ParseTriples(n)
	case "Graph":
		var t TriplesOrGraph
		switch n := n.Children()[0]; n.Name {
		case "LabelOrSubject":
			los, err := ParseLabelOrSubject(n)
			if err != nil {
				return nil, err
			}
			t.LabelOrSubject = los
		default:
			return nil, fmt.Errorf("triples or graph: unknown: %s", n.Name)
		}
		switch n := n.Children()[1]; n.Name {
		case "WrappedGraph":
			wg, err := ParseWrappedGraph(n)
			if err != nil {
				return nil, err
			}
			t.WrappedGraph = wg
		default:
			return nil, fmt.Errorf("triples or graph: unknown: %s", n.Name)
		}
		return &t, nil
	default:
		return nil, fmt.Errorf("block: unknown: %s", n.Name)
	}
}

type Document []Statement

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
		return nil, fmt.Errorf("document: unknown: %s", n.Name)
	}
	var doc Document
	for _, n := range n.Children() {
		switch n.Name {
		case "Directive":
			d, err := ttl.ParseDirective(n)
			if err != nil {
				return nil, err
			}
			switch d := d.(type) {
			case *ttl.Base:
				doc = append(doc, (*Base)(d))
			case *ttl.Prefix:
				doc = append(doc, (*Prefix)(d))
			default:
				return nil, fmt.Errorf("unknown directive type: %T", d)
			}
		case "Block":
			t, err := ParseBlock(n)
			if err != nil {
				return nil, err
			}
			doc = append(doc, t)
		default:
			return nil, fmt.Errorf("document: unknown: %s", n.Name)
		}
	}
	return doc, nil
}

func (d Document) String() string {
	var b strings.Builder
	for _, s := range d {
		b.WriteString(s.String())
		b.WriteString("\n")
	}
	return b.String()
}

type IRI ttl.IRI

func (i IRI) String() string {
	return ttl.IRI(i).String()
}

func (i IRI) labelOrSubject() {}

// LabelOrSubject is either an IRI or a BlankNode.
type LabelOrSubject interface {
	labelOrSubject()

	fmt.Stringer
}

func ParseLabelOrSubject(n *parser.Node) (LabelOrSubject, error) {
	if n.Name != "LabelOrSubject" {
		return nil, fmt.Errorf("label or subject: unknown: %s", n.Name)
	}
	switch n := n.Children()[0]; n.Name {
	case "IRI":
		iri, err := ttl.ParseIRI(n)
		if err != nil {
			return nil, err
		}
		return (*IRI)(iri), nil
	case "BlankNode":
		bn, err := ttl.ParseBlankNode(n)
		if err != nil {
			return nil, err
		}
		return (*BlankNode)(bn), nil
	default:
		return nil, fmt.Errorf("label or subject: unknown: %s", n.Name)
	}
}

type Prefix ttl.Prefix

func (p Prefix) String() string {
	return ttl.Prefix(p).String()
}

func (p Prefix) statement() {}

type Statement interface {
	statement()

	fmt.Stringer
}

type Triple2 struct {
	BlankNodePropertyList ttl.BlankNodePropertyList
	Collection            ttl.Collection
	PredicateObjectList   ttl.PredicateObjectList
}

func ParseTriples(n *parser.Node) (*Triple2, error) {
	if n.Name != "Triples2" {
		return nil, fmt.Errorf("triples: unknown: %s", n.Name)
	}

	var t Triple2
	switch n := n.Children()[0]; n.Name {
	case "Triples2BlankNodePropertyList":
		bnpl, err := ttl.ParseBlankNodePropertyList(n.Children()[0])
		if err != nil {
			return nil, err
		}
		t.BlankNodePropertyList = bnpl
		if len(n.Children()) == 2 {
			pol, err := ttl.ParsePredicateObjectList(n.Children()[1])
			if err != nil {
				return nil, err
			}
			t.PredicateObjectList = pol
		}
	case "Triples2Collection":
		c, err := ttl.ParseCollection(n.Children()[0])
		if err != nil {
			return nil, err
		}
		t.Collection = c
		pol, err := ttl.ParsePredicateObjectList(n.Children()[1])
		if err != nil {
			return nil, err
		}
		t.PredicateObjectList = pol
	}
	return &t, nil
}

func (t Triple2) String() string {
	if len(t.BlankNodePropertyList) != 0 {
		if len(t.PredicateObjectList) == 0 {
			return fmt.Sprintf("%s .", t.BlankNodePropertyList)
		}
		return fmt.Sprintf("%s %s .", t.BlankNodePropertyList, t.PredicateObjectList)
	}
	return fmt.Sprintf("%s %s .", t.Collection, t.PredicateObjectList)
}

func (t Triple2) block() {}

func (t Triple2) statement() {}

type TriplesBlock []ttl.Triple

func ParseTriplesBlock(n *parser.Node) (TriplesBlock, error) {
	if n.Name != "TriplesBlock" {
		return nil, fmt.Errorf("triples block: unknown: %s", n.Name)
	}
	var tb TriplesBlock
	for _, n := range n.Children() {
		switch n.Name {
		case "Triples":
			t, err := ttl.ParseTriples(n)
			if err != nil {
				return nil, err
			}
			tb = append(tb, *t)
		case "TriplesBlock":
			b, err := ParseTriplesBlock(n)
			if err != nil {
				return nil, err
			}
			tb = append(tb, b...)
		default:
			return nil, fmt.Errorf("triples block: unknown: %s", n.Name)
		}
	}
	return tb, nil
}

func (t TriplesBlock) String() string {
	var b strings.Builder
	for _, t := range t {
		b.WriteString(t.String())
	}
	return b.String()
}

type TriplesOrGraph struct {
	LabelOrSubject      LabelOrSubject
	WrappedGraph        WrappedGraph
	PredicateObjectList ttl.PredicateObjectList
}

func ParseTriplesOrGraph(n *parser.Node) (*TriplesOrGraph, error) {
	if n.Name != "TriplesOrGraph" {
		return nil, fmt.Errorf("triples or graph: unknown: %s", n.Name)
	}
	var t TriplesOrGraph
	switch n := n.Children()[0]; n.Name {
	case "LabelOrSubject":
		los, err := ParseLabelOrSubject(n)
		if err != nil {
			return nil, err
		}
		t.LabelOrSubject = los
	default:
		return nil, fmt.Errorf("triples or graph: unknown: %s", n.Name)
	}
	switch n := n.Children()[1]; n.Name {
	case "WrappedGraph":
		wg, err := ParseWrappedGraph(n)
		if err != nil {
			return nil, err
		}
		t.WrappedGraph = wg
	case "PredicateObjectList":
		pol, err := ttl.ParsePredicateObjectList(n)
		if err != nil {
			return nil, err
		}
		t.PredicateObjectList = pol
	default:
		return nil, fmt.Errorf("triples or graph: unknown: %s", n.Name)
	}
	return &t, nil
}

func (t *TriplesOrGraph) String() string {
	if t.WrappedGraph != nil {
		return fmt.Sprintf("%s %s", t.LabelOrSubject, t.WrappedGraph)
	}
	return fmt.Sprintf("%s %s .", t.LabelOrSubject, t.PredicateObjectList)
}

func (t *TriplesOrGraph) block() {}

func (t *TriplesOrGraph) statement() {}

type WrappedGraph TriplesBlock

func ParseWrappedGraph(n *parser.Node) (WrappedGraph, error) {
	if n.Name != "WrappedGraph" {
		return nil, fmt.Errorf("wrapped graph: unknown: %s", n.Name)
	}
	if len(n.Children()) == 1 {
		wg, err := ParseTriplesBlock(n.Children()[0])
		if err != nil {
			return nil, err
		}

		return WrappedGraph(wg), nil
	}
	return WrappedGraph{}, nil
}

func (w WrappedGraph) String() string {
	return fmt.Sprintf("{ %s }", TriplesBlock(w))
}

func (w WrappedGraph) block() {}

func (w WrappedGraph) statement() {}
