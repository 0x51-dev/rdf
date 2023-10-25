package turtle

import (
	"fmt"
	nt "github.com/0x51-dev/rdf/ntriples"
	"github.com/0x51-dev/rdf/turtle/grammar"
	"github.com/0x51-dev/upeg/parser"
	"github.com/0x51-dev/upeg/parser/op"
	"sort"
	"strings"
)

func EvaluateDocument(doc Document, cwd string) (nt.Document, error) {
	return NewContext().evaluateDocument(doc, cwd)
}

func ValidateDocument(doc Document) bool {
	return NewContext().validateDocument(doc)
}

type A struct{}

func (a A) Equal(v any) bool {
	if _, ok := v.(A); ok {
		return true
	}
	if v, ok := v.(*A); ok && v != nil {
		return true
	}
	return false
}

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

func (b Base) Equal(v any) bool {
	if v, ok := v.(Base); ok {
		return b == v
	}
	if v, ok := v.(*Base); ok && v != nil {
		return b == *v
	}
	return false
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

func (b BlankNode) Equal(v any) bool {
	if v, ok := v.(BlankNode); ok {
		return b.equal(v)
	}
	if v, ok := v.(*BlankNode); ok && v != nil {
		return b.equal(*v)
	}
	return false
}

func (b BlankNode) String() string {
	if b == "[]" {
		return string(b)
	}
	return fmt.Sprintf("_:%s", string(b))
}

func (b BlankNode) equal(other BlankNode) bool {
	return b == other
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
	sort.Sort(pol)
	return BlankNodePropertyList(pol), nil
}

func (b BlankNodePropertyList) Equal(v any) bool {
	if v, ok := v.(BlankNodePropertyList); ok {
		return b.equal(v)
	}
	if v, ok := v.(*BlankNodePropertyList); ok && v != nil {
		return b.equal(*v)
	}
	return false
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

func (b BlankNodePropertyList) equal(other BlankNodePropertyList) bool {
	return PredicateObjectList(b).equal(PredicateObjectList(other))
}

func (b BlankNodePropertyList) object() {}

type BooleanLiteral bool

func (bl BooleanLiteral) Equal(v any) bool {
	if v, ok := v.(BooleanLiteral); ok {
		return bl == v
	}
	if v, ok := v.(*BooleanLiteral); ok && v != nil {
		return bl == *v
	}
	return false
}

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

func (c Collection) Equal(v any) bool {
	if v, ok := v.(Collection); ok {
		return c.equal(v, false)
	}
	if v, ok := v.(*Collection); ok && v != nil {
		return c.equal(*v, false)
	}
	return false
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

func (c Collection) equal(other Collection, checkBlankNode bool) bool {
	if len(c) != len(other) {
		return false
	}
	for i, o := range c {
		switch o.(type) {
		case BlankNode, *BlankNode:
			if !checkBlankNode {
				_, isBn := other[i].(BlankNode)
				_, isBnPtr := other[i].(*BlankNode)
				if !isBn && !isBnPtr {
					return false
				}
			} else {
				if !o.Equal(other[i]) {
					return false
				}
			}
		default:
			if !o.Equal(other[i]) {
				return false
			}
		}
	}
	return true
}

func (c Collection) normalizeBlankNodes(f func(b fmt.Stringer) *BlankNode) (n Collection) {
	for _, v := range c {
		switch o := v.(type) {
		case BlankNode, *BlankNode:
			n = append(n, f(o))
		case BlankNodePropertyList:
			n = append(n, BlankNodePropertyList(
				PredicateObjectList(o).normalizeBlankNodes(f),
			))
		case Collection:
			n = append(n, o.normalizeBlankNodes(f))
		default:
			n = append(n, o)
		}
	}
	return
}

func (c Collection) object() {}

func (c Collection) subject() {}

type Directive interface {
	directive()

	Equal(v any) bool
	fmt.Stringer
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
	var document Document
	for _, n := range n.Children() {
		switch n.Name {
		case "Directive":
			d, err := ParseDirective(n)
			if err != nil {
				return nil, err
			}
			switch d := d.(type) {
			case *Base:
				document = append(document, d)
			case *Prefix:
				document = append(document, d)
			default:
				return nil, fmt.Errorf("document: unknown directive: %T", d)
			}
		case "Triples":
			t, err := ParseTriples(n)
			if err != nil {
				return nil, err
			}
			document = append(document, t)
		default:
			return nil, fmt.Errorf("document: unknown: %s", n.Name)
		}
	}
	sort.Sort(document)
	return document, nil
}

func (d Document) Equal(other Document) bool {
	if len(d) != len(other) {
		return false
	}
	d, o := d.normalizeBlankNodes(), other.normalizeBlankNodes()
	for i, s0 := range d {
		s1 := o[i]
		if !s0.Equal(s1) {
			return false
		}
	}
	return true
}

func (d Document) Len() int {
	return len(d)
}

func (d Document) Less(i, j int) bool {
	switch d0 := d[i].(type) {
	case Triple, *Triple:
		switch t1 := d[j].(type) {
		case Triple, *Triple:
			return d0.String() < t1.String()
		default:
			return false
		}
	case Prefix:
		switch p1 := d[j].(type) {
		case Prefix, *Prefix:
			return d0.Less(p1)
		default:
			return false
		}
	case *Prefix:
		switch p1 := d[j].(type) {
		case Prefix, *Prefix:
			return d0.Less(p1)
		default:
			return false
		}
	default:
		return false
	}
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

func (d Document) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

func (d Document) normalizeBlankNodes() (n Document) {
	var i int
	mapping := make(map[string]string)
	f := func(b fmt.Stringer) *BlankNode {
		if bn, ok := mapping[b.String()]; ok {
			return (*BlankNode)(&bn)
		}
		bn := fmt.Sprintf("%d", i)
		mapping[b.String()] = bn
		i++
		return (*BlankNode)(&bn)
	}
	for _, s := range d {
		switch t := s.(type) {
		case *Triple:
			if t.Subject != nil {
				subject := t.Subject
				switch s := t.Subject.(type) {
				case BlankNode, *BlankNode:
					subject = f(s)
				case Collection:
					subject = s.normalizeBlankNodes(f)
				case *Collection:
					if s != nil {
						subject = s.normalizeBlankNodes(f)
					}
				}
				n = append(n, &Triple{
					Subject:             subject,
					PredicateObjectList: t.PredicateObjectList.normalizeBlankNodes(f),
				})
			} else {
				var predicateObjectList PredicateObjectList
				if t.PredicateObjectList != nil {
					predicateObjectList = t.PredicateObjectList.normalizeBlankNodes(f)
				}
				n = append(n, &Triple{
					BlankNodePropertyList: BlankNodePropertyList(
						PredicateObjectList(t.BlankNodePropertyList).normalizeBlankNodes(f),
					),
					PredicateObjectList: predicateObjectList,
				})
			}
		default:
			n = append(n, s)
		}
	}
	return
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

func (i IRI) Equal(v any) bool {
	if v, ok := v.(IRI); ok {
		return i.equal(v)
	}
	if v, ok := v.(*IRI); ok && v != nil {
		return i.equal(*v)
	}
	return false
}

func (i IRI) String() string {
	if i.Prefixed {
		return i.Value
	}
	return fmt.Sprintf("<%s>", i.Value)
}

func (i IRI) equal(other IRI) bool {
	return i.Prefixed == other.Prefixed && i.Value == other.Value
}

func (i IRI) object() {}

func (i IRI) subject() {}

func (i IRI) verb() {}

// Literal is either StringLiteral, NumericLiteral or BooleanLiteral.
type Literal interface {
	literal()

	Object

	Equal(v any) bool
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
			iri, err := ParseIRI(n)
			if err != nil {
				return nil, err
			}
			v.DatatypeIRI = iri
		default:
			return nil, fmt.Errorf("rdf literal: unknown: %s", n.Name)
		}
	}
	return v, nil
}

type NumericLiteral struct {
	Type  NumericType
	Value string
}

func (nl NumericLiteral) Equal(v any) bool {
	if v, ok := v.(NumericLiteral); ok {
		return nl.equal(v)
	}
	if v, ok := v.(*NumericLiteral); ok && v != nil {
		return nl.equal(*v)
	}
	return false
}
func (nl NumericLiteral) String() string {
	return nl.Value
}

func (nl NumericLiteral) equal(other NumericLiteral) bool {
	return nl.Type == other.Type && nl.Value == other.Value
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

	Equal(v any) bool
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
type ObjectList Collection

func ParseObjectList(n *parser.Node) (ObjectList, error) {
	if n.Name != "ObjectList" {
		return nil, fmt.Errorf("object list: unknown %s", n.Name)
	}
	n.Name = "Collection"
	c, err := ParseCollection(n)
	if err != nil {
		return nil, err
	}
	ol := ObjectList(c)
	sort.Sort(ol)
	return ol, nil
}

func (ol ObjectList) Equal(v any) bool {
	if v, ok := v.(ObjectList); ok {
		return ol.equal(v)
	}
	if v, ok := v.(*ObjectList); ok && v != nil {
		return ol.equal(*v)
	}
	return false
}

func (ol ObjectList) Len() int {
	return len(ol)
}

func (ol ObjectList) Less(i, j int) bool {
	return ol[i].String() < ol[j].String()
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

func (ol ObjectList) Swap(i, j int) {
	ol[i], ol[j] = ol[j], ol[i]
}

func (ol ObjectList) equal(other ObjectList) bool {
	return Collection(ol).Equal(Collection(other))
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

func (po PredicateObject) Equal(v any) bool {
	if v, ok := v.(PredicateObject); ok {
		return po.equal(v)
	}
	if v, ok := v.(*PredicateObject); ok && v != nil {
		return po.equal(*v)
	}
	return false
}

func (po PredicateObject) String() string {
	var s string
	s += po.Verb.String()
	if po.ObjectList != nil {
		s += fmt.Sprintf(" %s", po.ObjectList)
	}
	return s
}

func (po PredicateObject) equal(other PredicateObject) bool {
	return po.Verb.Equal(other.Verb) && po.ObjectList.Equal(other.ObjectList)
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
	sort.Sort(pol)
	return pol, nil
}

func (pol PredicateObjectList) Equal(v any) bool {
	if v, ok := v.(PredicateObjectList); ok {
		return pol.equal(v)
	}
	if v, ok := v.(*PredicateObjectList); ok && v != nil {
		return pol.equal(*v)
	}
	return false
}

func (pol PredicateObjectList) Len() int {
	return len(pol)
}

func (pol PredicateObjectList) Less(i, j int) bool {
	return pol[i].String() < pol[j].String()
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

func (pol PredicateObjectList) Swap(i, j int) {
	pol[i], pol[j] = pol[j], pol[i]
}

func (pol PredicateObjectList) equal(other PredicateObjectList) bool {
	if len(pol) != len(other) {
		return false
	}
	for i, po := range pol {
		if !po.Equal(other[i]) {
			return false
		}
	}
	return true
}

func (pol PredicateObjectList) normalizeBlankNodes(f func(b fmt.Stringer) *BlankNode) (n PredicateObjectList) {
	for _, v := range pol {
		n = append(n, PredicateObject{
			Verb:       v.Verb,
			ObjectList: ObjectList(Collection(v.ObjectList).normalizeBlankNodes(f)),
		})
	}
	return
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

func (p Prefix) Equal(v any) bool {
	if v, ok := v.(Prefix); ok {
		return p.equal(v)
	}
	if v, ok := v.(*Prefix); ok && v != nil {
		return p.equal(*v)
	}
	return false
}

func (p Prefix) Less(other any) bool {
	if other, ok := other.(Prefix); ok {
		return p.less(other)
	}
	if other, ok := other.(*Prefix); ok && other != nil {
		return p.less(*other)
	}
	return false
}

func (p Prefix) String() string {
	return fmt.Sprintf("@prefix %s <%s> .", p.Name, p.IRI)
}

func (p Prefix) directive() {}

func (p Prefix) equal(other Prefix) bool {
	return p.Name == other.Name && p.IRI == other.IRI
}

func (p Prefix) less(other Prefix) bool {
	if p.Name == other.Name {
		return false
	}
	return p.String() < other.String()
}

func (p Prefix) statement() {}

type Statement interface {
	statement()

	Equal(v any) bool
	fmt.Stringer
}

type StringLiteral struct {
	Value       string
	Multiline   bool
	SingleQuote bool
	LanguageTag string
	DatatypeIRI *IRI
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

func (sl StringLiteral) Equal(v any) bool {
	if v, ok := v.(StringLiteral); ok {
		return sl.equal(v)
	}
	if v, ok := v.(*StringLiteral); ok && v != nil {
		return sl.equal(*v)
	}
	return false
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
	if sl.DatatypeIRI != nil {
		s += fmt.Sprintf("^^%s", sl.DatatypeIRI)
	}
	return s
}

func (sl StringLiteral) equal(other StringLiteral) bool {
	if sl.DatatypeIRI != nil && !sl.DatatypeIRI.Equal(other.DatatypeIRI) {
		return false
	}
	return sl.Value == other.Value && sl.Multiline == other.Multiline &&
		sl.SingleQuote == other.SingleQuote && sl.LanguageTag == other.LanguageTag
}

func (sl StringLiteral) literal() {}

func (sl StringLiteral) object() {}

type Subject interface {
	subject()

	Equal(v any) bool
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

func (t Triple) Equal(v any) bool {
	if v, ok := v.(Triple); ok {
		return t.equal(v)
	}
	if v, ok := v.(*Triple); ok && v != nil {
		return t.equal(*v)
	}
	return false
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

func (t Triple) equal(other Triple) bool {
	if t.Subject != nil {
		return t.Subject.Equal(other.Subject) && t.PredicateObjectList.Equal(other.PredicateObjectList)
	}
	if !t.BlankNodePropertyList.Equal(other.BlankNodePropertyList) {
		return false
	}
	if t.PredicateObjectList != nil {
		return t.PredicateObjectList.Equal(other.PredicateObjectList)
	}
	return true
}

func (t Triple) statement() {}

type Verb interface {
	verb()

	Equal(v any) bool
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
