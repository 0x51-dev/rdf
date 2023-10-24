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

func ParseBlankNodeLabel(n *parser.Node) (*BlankNode, error) {
	if n.Name != "BlankNodeLabel" {
		return nil, fmt.Errorf("blank-node: unknown %s", n.Name)
	}
	bn := BlankNode(n.Value())
	return &bn, nil
}

// Equal returns true if the blank node is equal to the given value.
func (n BlankNode) Equal(v any) bool {
	if v, ok := v.(BlankNode); ok {
		return n.equal(v)
	}
	if v, ok := v.(*BlankNode); ok && v != nil {
		return n.equal(*v)
	}
	return false
}

func (n BlankNode) String() string {
	return fmt.Sprintf("_:%s", string(n))
}

func (n BlankNode) equal(other BlankNode) bool {
	return n == other
}

func (n BlankNode) object() {}

func (n BlankNode) subject() {}

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

type IRIReference string

func ParseIRIReference(n *parser.Node) (*IRIReference, error) {
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

func ParsePredicate(n *parser.Node) (*IRIReference, error) {
	if n.Name != "Predicate" {
		return nil, fmt.Errorf("predicate: unknown %s", n.Name)
	}
	return ParseIRIReference(n.Children()[0])
}

// Equal returns true if the IRI reference is equal to the given value.
func (r IRIReference) Equal(v any) bool {
	if v, ok := v.(IRIReference); ok {
		return r.equal(v)
	}
	if v, ok := v.(*IRIReference); ok && v != nil {
		return r.equal(*v)
	}
	return false
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

func (r IRIReference) equal(other IRIReference) bool {
	return r == other
}

func (r IRIReference) object() {}

func (r IRIReference) subject() {}

type Literal struct {
	Value     string
	Reference *IRIReference
	Language  string
}

func ParseLiteral(n *parser.Node) (*Literal, error) {
	if n.Name != "Literal" {
		return nil, fmt.Errorf("literal: unknown %s", n.Name)
	}
	var literal Literal
	for _, n := range n.Children() {
		switch n.Name {
		case "StringLiteral":
			literal.Value = n.Value()
		case "IRIReference":
			ref, err := ParseIRIReference(n)
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

// Equal returns true if the literal is equal to the given value.
func (l Literal) Equal(v any) bool {
	if v, ok := v.(Literal); ok {
		return l.equal(v)
	}
	if v, ok := v.(*Literal); ok && v != nil {
		return l.equal(*v)
	}
	return false
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

func (l Literal) equal(other Literal) bool {
	if l.Reference != nil && !l.Reference.Equal(other.Reference) {
		return false
	}
	if l.Reference == nil && other.Reference != nil {
		return false
	}
	return l.Value == other.Value && l.Language == other.Language
}

func (l Literal) object() {}

// Object is either an IRI, a blank node, or a literal.
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
		return ParseIRIReference(n)
	case "BlankNodeLabel":
		return ParseBlankNodeLabel(n)
	case "Literal":
		return ParseLiteral(n)
	default:
		return nil, fmt.Errorf("object: unknown: %s", n.Name)
	}
}

// Subject is either an IRI or a blank node.
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
		return ParseIRIReference(n)
	case "BlankNodeLabel":
		return ParseBlankNodeLabel(n)
	default:
		return nil, fmt.Errorf("subject: unknown: %s", n.Name)
	}
}

type Triple struct {
	Subject   Subject
	Predicate IRIReference
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
	p, err := ParsePredicate(children[1])
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
