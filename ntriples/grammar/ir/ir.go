package ir

import (
	"fmt"
	"github.com/0x51-dev/ri/i"
	"github.com/0x51-dev/upeg/parser"
	"github.com/0x51-dev/upeg/parser/op"
	"strconv"
	"strings"
)

var validation = true

// DisableValidation disables validation of IRIs.
func DisableValidation() {
	validation = false
}

type BlankNode string

func ParseBlankNodeLabel(n *parser.Node) (BlankNode, error) {
	if n.Name != "BlankNodeLabel" {
		return "", fmt.Errorf("blank-node: unknown %s", n.Name)
	}
	return BlankNode(n.Value()), nil
}

func (n BlankNode) String() string {
	return string(n)
}

func (n BlankNode) object() {}

func (n BlankNode) subject() {}

type Document []Triple

func ParseDocument(n *parser.Node) (Document, error) {
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

func (d Document) String() string {
	var s string
	for _, t := range d {
		s += fmt.Sprintf("%s\n", t)
	}
	return s
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
	if _, err := p.Match(op.And{i.AbsoluteIRI, op.Optional{Value: op.And{'#', i.Ifragment}}, op.EOF{}}); err != nil {
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

func ParseLiteral(n *parser.Node) (Literal, error) {
	if n.Name != "Literal" {
		return Literal{}, fmt.Errorf("literal: unknown %s", n.Name)
	}
	var literal Literal
	for _, n := range n.Children() {
		switch n.Name {
		case "StringLiteral":
			literal.Value = n.Value()
		case "IRIReference":
			ref, err := ParseIRIReference(n)
			if err != nil {
				return literal, err
			}
			literal.Reference = ref
		case "LanguageTag":
			literal.Language = n.Value()
		default:
			return literal, fmt.Errorf("literal: unknown child: %s", n.Name)
		}
	}
	return literal, nil
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

type Subject interface {
	subject()
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

func (t Triple) String() string {
	return fmt.Sprintf("%s %s %s .", t.Subject, t.Predicate, t.Object)
}
