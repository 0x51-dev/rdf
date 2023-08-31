package ir_test

import (
	"fmt"
	"github.com/0x51-dev/rdf/ntriples/grammar"
	"github.com/0x51-dev/rdf/ntriples/grammar/ir"
	"github.com/0x51-dev/upeg/parser"
	"github.com/0x51-dev/upeg/parser/op"
	"testing"
)

func TestParseBlankNodeLabel(t *testing.T) {
	blank := "_:blank"
	n := getNode(t, blank, grammar.BlankNodeLabel)
	b, err := ir.ParseBlankNodeLabel(n)
	if err != nil {
		t.Error(err)
	}
	if b != ir.BlankNode(blank) {
		t.Error("invalid blank node", b)
	}
}

func TestParseIRIReference(t *testing.T) {
	ref := "http://example.com"
	n := getNode(t, fmt.Sprintf("<%s>", ref), grammar.IRIReference)
	r, err := ir.ParseIRIReference(n)
	if err != nil {
		t.Error(err)
	}
	if *r != ir.IRIReference(ref) {
		t.Error("invalid IRI reference", r)
	}
}

func TestParseLiteral(t *testing.T) {
	{
		lit := "literal"
		n := getNode(t, fmt.Sprintf(`"%s"`, lit), grammar.Literal)
		l, err := ir.ParseLiteral(n)
		if err != nil {
			t.Error(err)
		}
		if l.Value != lit {
			t.Error("invalid literal", l)
		}
	}
	{
		lit := "literal"
		ref := "http://example.com"
		n := getNode(t, fmt.Sprintf(`"%s"^^<%s>`, lit, ref), grammar.Literal)
		l, err := ir.ParseLiteral(n)
		if err != nil {
			t.Error(err)
		}
		if l.Value != lit {
			t.Error("invalid literal", l)
		}
		if *l.Reference != ir.IRIReference(ref) {
			t.Error("invalid literal", l)
		}
	}
	{
		lit := "literal"
		lang := "en"
		n := getNode(t, fmt.Sprintf(`"%s"@%s`, lit, lang), grammar.Literal)
		l, err := ir.ParseLiteral(n)
		if err != nil {
			t.Error(err)
		}
		if l.Value != lit {
			t.Error("invalid literal", l)
		}
		if l.Language != lang {
			t.Error("invalid literal", l)
		}
	}
}

func getNode(t *testing.T, i string, o any) *parser.Node {
	p, err := parser.New([]rune(i))
	if err != nil {
		t.Fatal(err)
	}
	n, err := p.Parse(op.And{o, op.EOF{}})
	if err != nil {
		t.Fatal(err)
	}
	return n
}
