package turtle

import (
	"fmt"
	nt "github.com/0x51-dev/rdf/ntriples"
	"github.com/0x51-dev/rdf/turtle/grammar"
	"github.com/0x51-dev/upeg/parser"
	"github.com/0x51-dev/upeg/parser/op"
	"strings"
)

func EvaluateDocument(doc Document) (nt.Document, error) {
	return NewContext().evaluateDocument(doc)
}

func ValidateDocument(doc Document) bool {
	return NewContext().validateDocument(doc)
}

type A struct{}

func (a A) String() string {
	return "a"
}

func (a A) verb() {}

type Base string

func ParseBase(n *parser.Node) (*Base, error) {
	if n.Name != "Base" {
		return nil, fmt.Errorf("Base: unknown %s", n.Name)
	}
	base := Base(n.Children()[0].Value())
	return &base, nil
}

func (b Base) String() string {
	return fmt.Sprintf("@base <%s> .", string(b))
}

func (b Base) directive() {}

func (b Base) statement() {}

type BlankNode string

func ParseBlankNode(n *parser.Node) (*BlankNode, error) {
	if n.Name != "BlankNode" {
		return nil, fmt.Errorf("blank node: unknown %s", n.Name)
	}
	switch n = n.Children()[0]; n.Name {
	case "BlankNodeLabel":
		bn := BlankNode(n.Value())
		return &bn, nil
	case "Anon":
		bn := BlankNode("[]")
		return &bn, nil
	default:
		return nil, fmt.Errorf("blank node: unknown: %s", n.Name)
	}
}

func (b BlankNode) String() string {
	if b == "[]" {
		return string(b)
	}
	return fmt.Sprintf("_:%s", string(b))
}

func (b BlankNode) object() {}

func (b BlankNode) subject() {}

// BlankNodePropertyList may appear in the subject or object position of a triple. That subject or object is a fresh RDF
// blank node. This blank node also serves as the subject of the triples produced by matching the predicateObjectList
// production embedded in a blankNodePropertyList.
type BlankNodePropertyList PredicateObjectList

func ParseBlankNodePropertyList(n *parser.Node) (BlankNodePropertyList, error) {
	if n.Name != "BlankNodePropertyList" {
		return nil, fmt.Errorf("blank node property list: unknown %s", n.Name)
	}
	pol, err := ParsePredicateObjectList(n.Children()[0])
	if err != nil {
		return nil, err
	}
	return BlankNodePropertyList(pol), nil
}

func (b BlankNodePropertyList) String() string {
	var s string
	s += "[ "
	for i, po := range b {
		if i > 0 {
			s += " ; "
		}
		s += po.String()
	}
	s += " ]"
	return s
}

func (b BlankNodePropertyList) object() {}

type BooleanLiteral bool

func (bl BooleanLiteral) String() string {
	return fmt.Sprintf("%t", bl)
}

func (bl BooleanLiteral) literal() {}

func (bl BooleanLiteral) object() {}

// Collection represents a rdf:first/rdf:rest list structure with the sequence of objects of the rdf:first statements
// being the order of the terms enclosed by (). It appears in the subject or object position of a triple. The blank node
// at the head of the list is the subject or object of the containing triple.
type Collection []Object

func ParseCollection(n *parser.Node) (Collection, error) {
	if n.Name != "Collection" {
		return nil, fmt.Errorf("collection: unknown %s", n.Name)
	}
	var collection Collection
	for _, n := range n.Children() {
		object, err := ParseObject(n)
		if err != nil {
			return nil, err
		}
		collection = append(collection, object)
	}
	return collection, nil
}

func (c Collection) String() string {
	var s string
	s += "("
	for i, o := range c {
		if i > 0 {
			s += " "
		}
		s += o.String()
	}
	s += ")"
	return s
}

func (c Collection) object() {}

func (c Collection) subject() {}

type Directive interface {
	directive()
}

func ParseDirective(n *parser.Node) (Directive, error) {
	if n.Name != "Directive" {
		return nil, fmt.Errorf("directive: unknown %s", n.Name)
	}
	switch n = n.Children()[0]; n.Name {
	case "Base":
		return ParseBase(n)
	case "Prefix":
		return ParsePrefix(n)
	default:
		return nil, fmt.Errorf("directive: unknown: %s", n.Name)
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
		return nil, fmt.Errorf("document: unknown %s", n.Name)
	}
	var doc Document
	for _, n := range n.Children() {
		switch n.Name {
		case "Directive":
			d, err := ParseDirective(n)
			if err != nil {
				return nil, err
			}
			switch d := d.(type) {
			case *Base:
				doc = append(doc, d)
			case *Prefix:
				doc = append(doc, d)
			default:
				return nil, fmt.Errorf("document: unknown directive: %T", d)
			}
		case "Triples":
			t, err := ParseTriples(n)
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
	var s string
	for _, l := range d {
		s += l.String()
		s += "\n"
	}
	return s
}

func (d Document) SubjectMap() (map[string]*Triple, error) {
	m := make(map[string]*Triple)
	for _, t := range d {
		switch t := t.(type) {
		case *Triple:
			var name string
			switch t := t.Subject.(type) {
			case *IRI:
				name = t.Value
			default:
				return nil, fmt.Errorf("document: subject not an iri: %T", t)
			}
			if _, ok := m[name]; ok {
				return nil, fmt.Errorf("document: duplicate subject: %s", t.Subject)
			}
			m[name] = t
		}
	}
	return m, nil
}

// IRI can be written as a relative/absolute IRI or prefixed name.
type IRI struct {
	Prefixed bool
	Value    string
}

func ParseIRI(n *parser.Node) (*IRI, error) {
	if n.Name != "IRI" {
		return nil, fmt.Errorf("iri: unknown %s", n.Name)
	}
	n = n.Children()[0]
	var prefixed bool
	switch n.Name {
	case "PrefixedName":
		prefixed = true
	case "IRIReference":
	default:
		return nil, fmt.Errorf("iri: unknown: %s", n.Name)
	}
	return &IRI{
		Prefixed: prefixed,
		Value:    n.Value(),
	}, nil
}

func (i IRI) String() string {
	if i.Prefixed {
		return i.Value
	}
	return fmt.Sprintf("<%s>", i.Value)
}

func (i IRI) object() {}

func (i IRI) subject() {}

func (i IRI) verb() {}

// Literal is either StringLiteral, NumericLiteral or BooleanLiteral.
type Literal interface {
	literal()

	Object

	fmt.Stringer
}

func ParseBooleanLiteral(n *parser.Node) (Literal, error) {
	if n.Name != "BooleanLiteral" {
		return nil, fmt.Errorf("boolean literal: unknown %s", n.Name)
	}
	switch n.Value() {
	case "true":
		b := BooleanLiteral(true)
		return &b, nil
	case "false":
		b := BooleanLiteral(false)
		return &b, nil
	default:
		return nil, fmt.Errorf("boolean literal: unknown: %s", n.Name)
	}
}

func ParseDecimal(n *parser.Node) (Literal, error) {
	if n.Name != "Decimal" {
		return nil, fmt.Errorf("decimal: unknown %s", n.Name)
	}
	return &NumericLiteral{
		Type:  Decimal,
		Value: n.Value(),
	}, nil
}

func ParseDouble(n *parser.Node) (Literal, error) {
	if n.Name != "Double" {
		return nil, fmt.Errorf("double: unknown %s", n.Name)
	}
	return &NumericLiteral{
		Type:  Double,
		Value: n.Value(),
	}, nil
}
func ParseInteger(n *parser.Node) (Literal, error) {
	if n.Name != "Integer" {
		return nil, fmt.Errorf("integer: unknown %s", n.Name)
	}
	return &NumericLiteral{
		Type:  Integer,
		Value: n.Value(),
	}, nil
}

func ParseLiteral(n *parser.Node) (Literal, error) {
	if n.Name != "Literal" {
		return nil, fmt.Errorf("literal: unknown %s", n.Name)
	}
	switch n = n.Children()[0]; n.Name {
	case "RDFLiteral":
		return ParseRDFLiteral(n)
	case "NumericLiteral":
		return ParseNumericLiteral(n)
	case "BooleanLiteral":
		return ParseBooleanLiteral(n)
	default:
		return nil, fmt.Errorf("literal: unknown: %s", n.Name)
	}
}

func ParseNumericLiteral(n *parser.Node) (Literal, error) {
	if n.Name != "NumericLiteral" {
		return nil, fmt.Errorf("numeric literal: unknown %s", n.Name)
	}
	switch n = n.Children()[0]; n.Name {
	case "Decimal":
		return ParseDecimal(n)
	case "Double":
		return ParseDouble(n)
	case "Integer":
		return ParseInteger(n)
	default:
		return nil, fmt.Errorf("numeric literal: unknown: %s", n.Name)
	}
}

func ParseRDFLiteral(n *parser.Node) (Literal, error) {
	if n.Name != "RDFLiteral" {
		return nil, fmt.Errorf("rdf literal: unknown %s", n.Name)
	}
	v, err := ParseStringLiteral(n.Children()[0])
	if err != nil {
		return nil, err
	}
	if len(n.Children()) > 1 {
		switch n := n.Children()[1]; n.Name {
		case "LanguageTag":
			v.LanguageTag = n.Value()
		case "IRI":
			v.DatatypeIRI = n.Value()
		}
	}
	return v, nil
}

type NumericLiteral struct {
	Type  NumericType
	Value string
}

func (nl NumericLiteral) String() string {
	return nl.Value
}

func (nl NumericLiteral) literal() {}

func (nl NumericLiteral) object() {}

type NumericType int

const (
	Decimal NumericType = iota
	Double
	Integer
)

type Object interface {
	object()

	fmt.Stringer
}

func ParseObject(n *parser.Node) (Object, error) {
	if n.Name != "Object" {
		return nil, fmt.Errorf("object: unknown %s", n.Name)
	}
	switch n = n.Children()[0]; n.Name {
	case "IRI":
		return ParseIRI(n)
	case "BlankNode":
		return ParseBlankNode(n)
	case "Collection":
		return ParseCollection(n)
	case "BlankNodePropertyList":
		return ParseBlankNodePropertyList(n)
	case "Literal":
		return ParseLiteral(n)
	default:
		return nil, fmt.Errorf("object: unknown: %s", n.Name)
	}
}

// ObjectList matches a series of objects separated by ',' following a predicate. This expresses a series of RDF Triples
// with the corresponding subject and predicate and each object allocated to one triple.
type ObjectList []Object

func ParseObjectList(n *parser.Node) (ObjectList, error) {
	if n.Name != "ObjectList" {
		return nil, fmt.Errorf("object list: unknown %s", n.Name)
	}
	var ol ObjectList
	for _, n := range n.Children() {
		if n.Name == "Object" {
			object, err := ParseObject(n)
			if err != nil {
				return nil, err
			}
			ol = append(ol, object)
		} else {
			for _, n := range n.Children() {
				object, err := ParseObject(n)
				if err != nil {
					return nil, err
				}
				ol = append(ol, object)
			}
		}
	}
	return ol, nil
}

func (ol ObjectList) String() string {
	var s string
	for i, o := range ol {
		if i > 0 {
			s += ", "
		}
		s += o.String()
	}
	return s
}

type PredicateObject struct {
	Verb       Verb
	ObjectList ObjectList
}

func ParsePredicateObject(n *parser.Node) (*PredicateObject, error) {
	if n.Name != "PredicateObject" {
		return nil, fmt.Errorf("predicate object: unknown %s", n.Name)
	}
	v, err := ParseVerb(n.Children()[0])
	if err != nil {
		return nil, err
	}
	ol, err := ParseObjectList(n.Children()[1])
	if err != nil {
		return nil, err
	}
	return &PredicateObject{
		Verb:       v,
		ObjectList: ol,
	}, nil
}

func (po PredicateObject) String() string {
	var s string
	s += po.Verb.String()
	if po.ObjectList != nil {
		s += fmt.Sprintf(" %s", po.ObjectList)
	}
	return s
}

// PredicateObjectList matches a series of predicates and objects, separated by ';', following a subject. This expresses
// a series of RDF Triples with that subject and each predicate and object allocated to one triple.
type PredicateObjectList []PredicateObject

func ParsePredicateObjectList(n *parser.Node) (PredicateObjectList, error) {
	if n.Name != "PredicateObjectList" {
		return nil, fmt.Errorf("predicate object list: unknown %s", n.Name)
	}
	var pol PredicateObjectList
	for _, n := range n.Children() {
		if n.Name == "PredicateObject" {
			po, err := ParsePredicateObject(n)
			if err != nil {
				return nil, err
			}
			pol = append(pol, *po)
			continue
		} else {
			for _, n := range n.Children() {
				po, err := ParsePredicateObject(n)
				if err != nil {
					return nil, err
				}
				pol = append(pol, *po)
			}
		}
	}
	return pol, nil
}

func (pol PredicateObjectList) String() string {
	var s string
	for i, po := range pol {
		if i > 0 {
			s += " ; "
		}
		s += po.String()
	}
	return s
}

type Prefix struct {
	Name string
	IRI  string
}

func ParsePrefix(n *parser.Node) (*Prefix, error) {
	if n.Name != "Prefix" {
		return nil, fmt.Errorf("prefix: unknown %s", n.Name)
	}
	return &Prefix{
		Name: n.Children()[0].Value(),
		IRI:  n.Children()[1].Value(),
	}, nil
}

func (p Prefix) String() string {
	return fmt.Sprintf("@prefix %s <%s> .", p.Name, p.IRI)
}

func (p Prefix) directive() {}

func (p Prefix) statement() {}

type Statement interface {
	statement()

	fmt.Stringer
}

type StringLiteral struct {
	Value       string
	Multiline   bool
	SingleQuote bool
	LanguageTag string
	DatatypeIRI string
}

func ParseStringLiteral(n *parser.Node) (*StringLiteral, error) {
	switch n.Name {
	case "StringLiteral":
		return &StringLiteral{
			Value: n.Value(),
		}, nil
	case "StringLiteralSQ":
		return &StringLiteral{
			Value:       n.Value(),
			SingleQuote: true,
		}, nil
	case "StringLiteralLQ":
		return &StringLiteral{
			Value:     n.Value(),
			Multiline: true,
		}, nil
	case "StringLiteralLSQ":
		return &StringLiteral{
			Value:       n.Value(),
			SingleQuote: true,
			Multiline:   true,
		}, nil
	default:
		return nil, fmt.Errorf("string literal: unknown %s", n.Name)
	}
}

func (sl StringLiteral) String() string {
	var s string
	if sl.Multiline {
		if sl.SingleQuote {
			s = fmt.Sprintf(`'''%s'''`, sl.Value)
		} else {
			s = fmt.Sprintf(`"""%s"""`, sl.Value)
		}
	} else {
		if sl.SingleQuote {
			s = fmt.Sprintf(`'%s'`, sl.Value)
		} else {
			s = fmt.Sprintf(`"%s"`, sl.Value)
		}
	}
	if sl.LanguageTag != "" {
		s += fmt.Sprintf("@%s", sl.LanguageTag)
	}
	if sl.DatatypeIRI != "" {
		s += fmt.Sprintf("^^%s", sl.DatatypeIRI)
	}
	return s
}

func (sl StringLiteral) literal() {}

func (sl StringLiteral) object() {}

type Subject interface {
	subject()

	fmt.Stringer
}

func ParseSubject(n *parser.Node) (Subject, error) {
	if n.Name != "Subject" {
		return nil, fmt.Errorf("subject: unknown %s", n.Name)
	}
	switch n = n.Children()[0]; n.Name {
	case "IRI":
		return ParseIRI(n)
	case "BlankNode":
		return ParseBlankNode(n)
	case "Collection":
		return ParseCollection(n)
	default:
		return nil, fmt.Errorf("subject: unknown: %s", n.Name)
	}
}

type Triple struct {
	Subject               Subject
	BlankNodePropertyList BlankNodePropertyList
	PredicateObjectList   PredicateObjectList
}

func ParseTripleBlankNodePropertyList(n *parser.Node) (*Triple, error) {
	if n.Name != "TripleBlankNodePropertyList" {
		return nil, fmt.Errorf("triple blank node property list: unknown %s", n.Name)
	}
	bnpl, err := ParseBlankNodePropertyList(n.Children()[0])
	if err != nil {
		return nil, err
	}
	if len(n.Children()) == 2 {
		pol, err := ParsePredicateObjectList(n.Children()[1])
		if err != nil {
			return nil, err
		}
		return &Triple{
			BlankNodePropertyList: bnpl,
			PredicateObjectList:   pol,
		}, nil
	}
	return &Triple{BlankNodePropertyList: bnpl}, nil
}

func ParseTripleSubject(n *parser.Node) (*Triple, error) {
	if n.Name != "TripleSubject" {
		return nil, fmt.Errorf("triple subject: unknown %s", n.Name)
	}
	var triple Triple
	s, err := ParseSubject(n.Children()[0])
	if err != nil {
		return nil, err
	}
	triple.Subject = s
	pl, err := ParsePredicateObjectList(n.Children()[1])
	if err != nil {
		return nil, err
	}
	triple.PredicateObjectList = pl
	return &triple, nil
}

func ParseTriples(n *parser.Node) (*Triple, error) {
	if n.Name != "Triples" {
		return nil, fmt.Errorf("Triples: unknown %s", n.Name)
	}
	switch n = n.Children()[0]; n.Name {
	case "TripleSubject":
		return ParseTripleSubject(n)
	case "TripleBlankNodePropertyList":
		return ParseTripleBlankNodePropertyList(n)
	default:
		return nil, fmt.Errorf("Triples: unknown: %s", n.Name)
	}
}

func (t Triple) PredicateObjectMap() (map[string]ObjectList, error) {
	m := make(map[string]ObjectList)
	for _, po := range t.PredicateObjectList {
		name := po.Verb.String()
		m[name] = po.ObjectList
	}
	return m, nil
}

func (t Triple) String() string {
	var s string
	if t.Subject != nil {
		s += t.Subject.String()
	}
	if t.BlankNodePropertyList != nil {
		s += t.BlankNodePropertyList.String()
	}
	if t.PredicateObjectList != nil {
		s += fmt.Sprintf(" %s", t.PredicateObjectList)
	}
	s += " ."
	return s
}

func (t Triple) statement() {}

type Verb interface {
	verb()

	fmt.Stringer
}

func ParseVerb(n *parser.Node) (Verb, error) {
	if n.Name != "Verb" {
		return nil, fmt.Errorf("verb: unknown %s", n.Name)
	}
	switch n = n.Children()[0]; n.Name {
	case "IRI":
		iri, err := ParseIRI(n)
		if err != nil {
			return nil, err
		}
		return iri, nil
	case "a":
		return &A{}, nil
	default:
		return nil, fmt.Errorf("verb: unknown: %s", n.Name)
	}
}
