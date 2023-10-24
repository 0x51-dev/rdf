package ntriples

import (
	"fmt"
	nt "github.com/0x51-dev/rdf/ntriples"
	"github.com/0x51-dev/rdf/star/ntriples/grammar"
	"github.com/0x51-dev/upeg/parser"
	"github.com/0x51-dev/upeg/parser/op"
	"strings"
)

type BlankNode nt.BlankNode

func (b BlankNode) Equal(v any) bool {
	return nt.BlankNode(b).Equal(v)
}

func (b BlankNode) String() string {
	return nt.BlankNode(b).String()
}

func (b BlankNode) object() {}

func (b BlankNode) subject() {}

type Document []Triple

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

// Equal returns true if the document is equal to the given value. Both document will be assumed sorted.
// NOTE: blank nodes will be compared, not by value, but by relation in the document.
func (d Document) Equal(other Document) bool {
	if len(d) != len(other) {
		return false
	}
	d, o := d.normalizeBlankNodes(), other.normalizeBlankNodes()
	for i := range d {
		t0, t1 := d[i], o[i]
		if !t0.equal(t1, true) {
			return false
		}
	}
	return true
}

func (d Document) String() string {
	var s string
	for _, t := range d {
		s += fmt.Sprintf("%s\n", t)
	}
	return s
}

func (d Document) normalizeBlankNodes() (n Document) {
	var i int
	mapping := make(map[string]string)
	for _, t := range d {
		subject := t.Subject
		if b, ok := t.Subject.(BlankNode); ok {
			if bn, ok := mapping[string(b)]; !ok {
				bn = fmt.Sprintf("%d", i)
				mapping[string(b)] = bn
				subject = (*BlankNode)(&bn)
				i++
			} else {
				subject = (*BlankNode)(&bn)
			}
		}
		predicate := t.Predicate
		object := t.Object
		if b, ok := t.Object.(BlankNode); ok {
			if bn, ok := mapping[string(b)]; !ok {
				bn = fmt.Sprintf("%d", i)
				mapping[string(b)] = bn
				object = (*BlankNode)(&bn)
				i++
			} else {
				object = (*BlankNode)(&bn)
			}
		}
		n = append(n, Triple{
			Subject:   subject,
			Predicate: predicate,
			Object:    object,
		})
	}
	return
}

type IRIReference nt.IRIReference

func (i IRIReference) Equal(v any) bool {
	return nt.IRIReference(i).Equal(v)
}

func (i IRIReference) String() string {
	return nt.IRIReference(i).String()
}

func (i IRIReference) object() {}

func (i IRIReference) subject() {}

type Literal nt.Literal

func (l Literal) Equal(v any) bool {
	return nt.Literal(l).Equal(v)
}

func (l Literal) String() string {
	return nt.Literal(l).String()
}

func (l Literal) object() {}

// Object is either an IRI, a blank node, literal or quoted triple.
type Object interface {
	object()

	Equal(v any) bool
	fmt.Stringer
}

func ParseObject(n *parser.Node) (Object, error) {
	if n.Name != "Object" {
		return nil, fmt.Errorf("object: unknown: %s", n.Name)
	}
	switch n = n.Children()[0]; n.Name {
	case "IRIReference":
		iri, err := nt.ParseIRIReference(n)
		if err != nil {
			return nil, err
		}
		return (*IRIReference)(iri), nil
	case "BlankNodeLabel":
		bn, err := nt.ParseBlankNodeLabel(n)
		if err != nil {
			return nil, err
		}
		return (*BlankNode)(bn), nil
	case "Literal":
		l, err := nt.ParseLiteral(n)
		if err != nil {
			return nil, err
		}
		return (*Literal)(l), nil
	case "QuotedTriple":
		return ParseQuotedTriple(n)
	default:
		return nil, fmt.Errorf("object: unknown: %s", n.Name)
	}
}

type QuotedTriple struct {
	Triple
}

func ParseQuotedTriple(n *parser.Node) (*QuotedTriple, error) {
	if n.Name != "QuotedTriple" {
		return nil, fmt.Errorf("quoted triple: unknown: %s", n.Name)
	}
	n.Name = "Triple"
	t, err := ParseTriple(n)
	if err != nil {
		return nil, err
	}
	return &QuotedTriple{Triple: *t}, nil
}

func (t QuotedTriple) Equal(v any) bool {
	if o, ok := v.(QuotedTriple); ok {
		return t.equal(o.Triple, true)
	}
	if o, ok := v.(*QuotedTriple); ok && o != nil {
		return t.equal(o.Triple, true)
	}
	return false
}

func (t QuotedTriple) String() string {
	return fmt.Sprintf("<<%s %s %s>>", t.Subject, t.Predicate, t.Object)
}

func (t QuotedTriple) object() {}

func (t QuotedTriple) subject() {}

// Subject is either an IRI, blank node or a quoted triple.
type Subject interface {
	subject()

	Equal(v any) bool
	fmt.Stringer
}

func ParseSubject(n *parser.Node) (Subject, error) {
	if n.Name != "Subject" {
		return nil, fmt.Errorf("subject: unknown: %s", n.Name)
	}
	switch n = n.Children()[0]; n.Name {
	case "IRIReference":
		iri, err := nt.ParseIRIReference(n)
		if err != nil {
			return nil, err
		}
		return (*IRIReference)(iri), nil
	case "BlankNodeLabel":
		bn, err := nt.ParseBlankNodeLabel(n)
		if err != nil {
			return nil, err
		}
		return (*BlankNode)(bn), nil
	case "QuotedTriple":
		return ParseQuotedTriple(n)
	default:
		return nil, fmt.Errorf("subject: unknown: %s", n.Name)
	}
}

type Triple struct {
	Subject   Subject
	Predicate nt.IRIReference
	Object    Object
}

func ParseTriple(n *parser.Node) (*Triple, error) {
	if n.Name != "Triple" {
		return nil, fmt.Errorf("triple: unknown %s", n.Name)
	}
	if len(n.Children()) != 3 {
		return nil, fmt.Errorf("triple: expected 3 children")
	}
	children := n.Children()
	s, err := ParseSubject(children[0])
	if err != nil {
		return nil, err
	}
	p, err := nt.ParsePredicate(children[1])
	if err != nil {
		return nil, err
	}
	o, err := ParseObject(children[2])
	if err != nil {
		return nil, err
	}
	return &Triple{
		Subject:   s,
		Predicate: *p,
		Object:    o,
	}, nil
}

// Equal returns true if the triple is equal to the given value.
// NOTE: comparing blank nodes does not really make sense, as they are not globally unique.
func (t Triple) Equal(v any) bool {
	if o, ok := v.(Triple); ok {
		return t.equal(o, false)
	}
	if o, ok := v.(*Triple); ok && o != nil {
		return t.equal(*o, false)
	}
	return false
}

func (t Triple) String() string {
	return fmt.Sprintf("%s %s %s .", t.Subject, t.Predicate, t.Object)
}

func (t Triple) equal(other Triple, checkBlankNode bool) bool {
	switch t.Subject.(type) {
	case BlankNode:
		if !checkBlankNode {
			_, isBn := other.Subject.(BlankNode)
			_, isBnPtr := other.Subject.(*BlankNode)
			if !isBn && !isBnPtr {
				return false
			}
		} else {
			if !t.Subject.Equal(other.Subject) {
				return false
			}
		}
	default:
		if !t.Subject.Equal(other.Subject) {
			return false
		}
	}

	if !t.Predicate.Equal(other.Predicate) {
		return false
	}

	switch t.Object.(type) {
	case BlankNode:
		if !checkBlankNode {
			_, isBn := other.Subject.(BlankNode)
			_, isBnPtr := other.Subject.(*BlankNode)
			if !isBn && !isBnPtr {
				return false
			}
		} else {
			if !t.Subject.Equal(other.Subject) {
				return false
			}
		}
	default:
		if !t.Object.Equal(other.Object) {
			return false
		}
	}
	return true
}
