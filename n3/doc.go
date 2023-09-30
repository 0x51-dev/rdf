package n3

import (
	"fmt"
	"github.com/0x51-dev/rdf/n3/grammar"
	ttl "github.com/0x51-dev/rdf/turtle"
	"github.com/0x51-dev/upeg/parser"
	"github.com/0x51-dev/upeg/parser/op"
	"strings"
)

type A struct{}

func (a A) String() string {
	return "a"
}

func (a A) verb() {}

type Base ttl.Base

func (b Base) String() string {
	return fmt.Sprintf("@base <%s>", string(b))
}

func (b Base) statement() {}

type BlankNode ttl.BlankNode

func (b BlankNode) String() string {
	return ttl.BlankNode(b).String()
}

func (b BlankNode) object() {}

func (b BlankNode) subject() {}

func (b BlankNode) verb() {}

type BlankNodePropertyList PredicateObjectList

func parseBlankNodePropertyList(n *parser.Node) (BlankNodePropertyList, error) {
	if n.Name != "BlankNodePropertyList" {
		return nil, fmt.Errorf("blank node property list: unknown %s", n.Name)
	}
	pol, err := parsePredicateObjectList(n.Children()[0])
	if err != nil {
		return nil, err
	}
	return BlankNodePropertyList(pol), nil
}

func (bnpl BlankNodePropertyList) String() string {
	return fmt.Sprintf("[ %s ]", PredicateObjectList(bnpl))
}

func (bnpl BlankNodePropertyList) object() {}

func (bnpl BlankNodePropertyList) subject() {}

func (bnpl BlankNodePropertyList) verb() {}

type BooleanLiteral ttl.BooleanLiteral

func (b BooleanLiteral) String() string {
	return ttl.BooleanLiteral(b).String()
}

func (b BooleanLiteral) literal() {}

func (b BooleanLiteral) object() {}

func (b BooleanLiteral) subject() {}

func (b BooleanLiteral) verb() {}

type Collection []Object

func parseCollection(n *parser.Node) (Collection, error) {
	if n.Name != "Collection" {
		return nil, fmt.Errorf("collection: unknown %s", n.Name)
	}
	var c Collection
	for _, n := range n.Children() {
		switch n.Name {
		case "Object":
			o, err := parseObject(n)
			if err != nil {
				return nil, err
			}
			c = append(c, o)
		default:
			return nil, fmt.Errorf("collection: unknown %s", n.Name)
		}
	}
	return c, nil
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

func (c Collection) verb() {}

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
		case "Statement":
			s, err := parseStatement(n)
			if err != nil {
				return nil, err
			}
			doc = append(doc, s)
		default:
			return nil, fmt.Errorf("document: unknown %s", n.Name)
		}
	}
	return doc, nil
}

func (d Document) String() string {
	var s string
	for _, l := range d {
		s += l.String()
		s += " .\n"
	}
	return s
}

type Equal struct{}

func (e Equal) String() string {
	return "="
}

func (e Equal) verb() {}

type EqualOrGreaterThan struct{}

func (eog EqualOrGreaterThan) String() string {
	return "=>"
}

func (eog EqualOrGreaterThan) verb() {}

type EqualOrLessThan struct{}

func (eol EqualOrLessThan) String() string {
	return "<="
}

func (eol EqualOrLessThan) verb() {}

type Formula FormulaContent

func parseFormula(n *parser.Node) (*Formula, error) {
	if n.Name != "Formula" {
		return nil, fmt.Errorf("formula: unknown %s", n.Name)
	}
	fc, err := parseFormulaContent(n.Children()[0])
	if err != nil {
		return nil, err
	}
	return (*Formula)(fc), nil
}

func (f Formula) String() string {
	return fmt.Sprintf("{ %s }", FormulaContent(f))
}

func (f Formula) object() {}

func (f Formula) subject() {}

func (f Formula) verb() {}

type FormulaContent struct {
	Statement Statement
	Formula   *FormulaContent
}

func parseFormulaContent(n *parser.Node) (*FormulaContent, error) {
	if n.Name != "FormulaContent" {
		return nil, fmt.Errorf("formula content: unknown %s", n.Name)
	}
	var fc FormulaContent
	switch n := n.Children()[0]; n.Name {
	case "Statement":
		s, err := parseStatement(n)
		if err != nil {
			return nil, err
		}
		fc.Statement = s
	default:
		return nil, fmt.Errorf("formula content: unknown %s", n.Name)
	}
	if len(n.Children()) == 2 {
		switch n := n.Children()[1]; n.Name {
		case "FormulaContent":
			f, err := parseFormulaContent(n)
			if err != nil {
				return nil, err
			}
			fc.Formula = f
		default:
			return nil, fmt.Errorf("formula content: unknown %s", n.Name)
		}
	}
	return &fc, nil
}

func (fc FormulaContent) String() string {
	switch fc.Statement.(type) {
	case *Triple:
		if fc.Formula != nil {
			return fmt.Sprintf("%s . %s", fc.Statement, fc.Formula)
		}
		return fc.Statement.String()
	default:
		if fc.Formula != nil {
			return fmt.Sprintf("%s%s", fc.Statement, fc.Formula)
		}
		return fc.Statement.String()
	}
}

type Has Path

func (h Has) String() string {
	return fmt.Sprintf("has %s", Path(h))
}

func (h Has) verb() {}

type IRI ttl.IRI

func (i IRI) String() string {
	return ttl.IRI(i).String()
}

func (i IRI) object() {}

func (i IRI) subject() {}

func (i IRI) verb() {}

type Inverse Path

func (i Inverse) String() string {
	return fmt.Sprintf("<-%s", Path(i))
}

func (i Inverse) verb() {}

type IsOf Path

func (i IsOf) String() string {
	return fmt.Sprintf("is %s of", Path(i))
}

func (i IsOf) verb() {}

type Literal interface {
	literal()

	Object

	fmt.Stringer
}

type NumericLiteral ttl.NumericLiteral

func (n NumericLiteral) String() string {
	return ttl.NumericLiteral(n).String()
}

func (n NumericLiteral) literal() {}

func (n NumericLiteral) object() {}

func (n NumericLiteral) subject() {}

func (n NumericLiteral) verb() {}

type Object interface {
	object()

	fmt.Stringer
}

func parseObject(n *parser.Node) (Object, error) {
	if n.Name != "Object" {
		return nil, fmt.Errorf("object: unknown %s", n.Name)
	}
	switch n := n.Children()[0]; n.Name {
	case "Path":
		return parsePath(n)
	default:
		return nil, fmt.Errorf("object: unknown %s", n.Name)
	}
}

type ObjectList []Object

func parseObjectList(n *parser.Node) (ObjectList, error) {
	if n.Name != "ObjectList" {
		return nil, fmt.Errorf("object list: unknown %s", n.Name)
	}
	var ol ObjectList
	for _, n := range n.Children() {
		if n.Name == "Object" {
			object, err := parseObject(n)
			if err != nil {
				return nil, err
			}
			ol = append(ol, object)
		} else {
			for _, n := range n.Children() {
				object, err := parseObject(n)
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

type Path struct {
	PathItem PathItem
	Forward  *Path
	Reverse  *Path
}

func parsePath(n *parser.Node) (*Path, error) {
	if n.Name != "Path" {
		return nil, fmt.Errorf("path: unknown %s", n.Name)
	}
	var p Path
	switch n := n.Children()[0]; n.Name {
	case "PathItem":
		pi, err := parsePathItem(n)
		if err != nil {
			return nil, err
		}
		p.PathItem = pi
	default:
		return nil, fmt.Errorf("path: unknown %s", n.Name)
	}
	return &p, nil
}

func (p Path) String() string {
	if p.Forward != nil {
		return fmt.Sprintf("%s!%s", p.PathItem, p.Forward)
	}
	if p.Reverse != nil {
		return fmt.Sprintf("%s^%s", p.PathItem, p.Reverse)
	}
	return p.PathItem.String()
}

func (p Path) object() {}

func (p Path) simple() (PathItem, bool) {
	return p.PathItem, p.Forward == nil && p.Reverse == nil
}

func (p Path) subject() {}

func (p Path) verb() {}

type PathItem interface {
	Subject
	Object
	Verb

	fmt.Stringer
}

func parsePathItem(n *parser.Node) (PathItem, error) {
	if n.Name != "PathItem" {
		return nil, fmt.Errorf("path item: unknown %s", n.Name)
	}
	switch n := n.Children()[0]; n.Name {
	case "IRI":
		i, err := ttl.ParseIRI(n)
		if err != nil {
			return nil, err
		}
		return (*IRI)(i), nil
	case "BlankNode":
		b, err := ttl.ParseBlankNode(n)
		if err != nil {
			return nil, err
		}
		return (*BlankNode)(b), nil
	case "QuickVar":
		return parseQuickVar(n)
	case "Collection":
		return parseCollection(n)
	case "BlankNodePropertyList":
		return parseBlankNodePropertyList(n)
	case "Literal":
		l, err := ttl.ParseLiteral(n)
		if err != nil {
			return nil, err
		}
		switch l := l.(type) {
		case *ttl.BooleanLiteral:
			return (*BooleanLiteral)(l), nil
		case *ttl.NumericLiteral:
			return (*NumericLiteral)(l), nil
		case *ttl.StringLiteral:
			return (*StringLiteral)(l), nil
		default:
			return nil, fmt.Errorf("path item: unknown %s", n.Name)
		}
	case "Formula":
		return parseFormula(n)
	default:
		return nil, fmt.Errorf("path item: unknown %s", n.Name)
	}
}

type PredicateObject struct {
	Verb       Verb
	ObjectList ObjectList
}

func parsePredicateObject(n *parser.Node) (*PredicateObject, error) {
	if n.Name != "PredicateObject" {
		return nil, fmt.Errorf("predicate object: unknown %s", n.Name)
	}
	v, err := parseVerb(n.Children()[0])
	if err != nil {
		return nil, err
	}
	ol, err := parseObjectList(n.Children()[1])
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

type PredicateObjectList []PredicateObject

func parsePredicateObjectList(n *parser.Node) (PredicateObjectList, error) {
	if n.Name != "PredicateObjectList" {
		return nil, fmt.Errorf("predicate object list: unknown %s", n.Name)
	}
	var pol PredicateObjectList
	for _, n := range n.Children() {
		switch n.Name {
		case "PredicateObject":
			po, err := parsePredicateObject(n)
			if err != nil {
				return nil, err
			}
			pol = append(pol, *po)
		default:
			return nil, fmt.Errorf("predicate object list: unknown %s", n.Name)
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

type Prefix ttl.Prefix

func (p Prefix) String() string {
	return fmt.Sprintf("@prefix %s <%s>", p.Name, p.IRI)
}

func (p Prefix) statement() {}

type QuickVar string

func parseQuickVar(n *parser.Node) (*QuickVar, error) {
	if n.Name != "QuickVar" {
		return nil, fmt.Errorf("quick var: unknown %s", n.Name)
	}
	qv := QuickVar(n.Children()[0].Value())
	return &qv, nil
}

func (qv QuickVar) String() string {
	return fmt.Sprintf("?%s", string(qv))
}

func (qv QuickVar) object() {}

func (qv QuickVar) subject() {}

func (qv QuickVar) verb() {}

type Statement interface {
	statement()

	fmt.Stringer
}

func parseStatement(n *parser.Node) (Statement, error) {
	if n.Name != "Statement" {
		return nil, fmt.Errorf("statement: unknown %s", n.Name)
	}
	switch n := n.Children()[0]; n.Name {
	case "Triples":
		return parseTriples(n)
	case "Prefix":
		p, err := ttl.ParsePrefix(n)
		if err != nil {
			return nil, err
		}
		return (*Prefix)(p), nil
	case "Base":
		b, err := ttl.ParseBase(n)
		if err != nil {
			return nil, err
		}
		return (*Base)(b), nil
	default:
		return nil, fmt.Errorf("statement: unknown %s", n.Name)
	}
}

type StringLiteral ttl.StringLiteral

func (s StringLiteral) String() string {
	return ttl.StringLiteral(s).String()
}

func (s StringLiteral) literal() {}

func (s StringLiteral) object() {}

func (s StringLiteral) subject() {}

func (s StringLiteral) verb() {}

type Subject interface {
	subject()

	fmt.Stringer
}

func parseSubject(n *parser.Node) (Subject, error) {
	if n.Name != "Subject" {
		return nil, fmt.Errorf("subject: unknown %s", n.Name)
	}
	switch n := n.Children()[0]; n.Name {
	case "Path":
		p, err := parsePath(n)
		if err != nil {
			return nil, err
		}
		if s, ok := p.simple(); ok {
			return s, nil
		}
		return p, nil
	default:
		return nil, fmt.Errorf("subject: unknown %s", n.Name)
	}
}

type Triple struct {
	Subject             Subject
	PredicateObjectList PredicateObjectList
}

func parseTriples(n *parser.Node) (*Triple, error) {
	if n.Name != "Triples" {
		return nil, fmt.Errorf("triples: unknown %s", n.Name)
	}
	var t Triple
	for _, n := range n.Children() {
		switch n.Name {
		case "Subject":
			s, err := parseSubject(n)
			if err != nil {
				return nil, err
			}
			t.Subject = s
		case "PredicateObjectList":
			pol, err := parsePredicateObjectList(n)
			if err != nil {
				return nil, err
			}
			t.PredicateObjectList = pol
		default:
			return nil, fmt.Errorf("triples: unknown %s", n.Name)
		}
	}
	return &t, nil
}

func (t *Triple) String() string {
	return fmt.Sprintf("%s %s", t.Subject, t.PredicateObjectList)
}

func (t *Triple) statement() {}

type Verb interface {
	verb()

	fmt.Stringer
}

func parseVerb(n *parser.Node) (Verb, error) {
	if n.Name != "Verb" {
		return nil, fmt.Errorf("verb: unknown %s", n.Name)
	}
	switch n := n.Children()[0]; n.Name {
	case "Path":
		p, err := parsePath(n)
		if err != nil {
			return nil, err
		}
		if s, ok := p.simple(); ok {
			return s, nil
		}
		return p, nil
	case "a":
		return &A{}, nil
	case "has":
		p, err := parsePath(n.Children()[0])
		if err != nil {
			return nil, err
		}
		return (*Has)(p), nil
	case "isOf":
		p, err := parsePath(n.Children()[0])
		if err != nil {
			return nil, err
		}
		return (*IsOf)(p), nil
	case "eq":
		return &Equal{}, nil
	case "eqg":
		return &EqualOrGreaterThan{}, nil
	case "eql":
		return &EqualOrLessThan{}, nil
	case "inverse":
		p, err := parsePath(n.Children()[0])
		if err != nil {
			return nil, err
		}
		return (*Inverse)(p), nil
	default:
		return nil, fmt.Errorf("verb: unknown %s", n.Name)
	}
}
