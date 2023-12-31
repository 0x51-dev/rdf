package grammar_test

import (
	. "github.com/0x51-dev/rdf/turtle/grammar"
	"github.com/0x51-dev/upeg/parser/op"
	"testing"
)

func TestBase(t *testing.T) {
	for _, test := range []string{
		"@base <http://example.org/> .",
	} {
		p, err := NewParser([]rune(test))
		if err != nil {
			t.Fatal(err)
		}
		if _, err := p.Parse(op.And{Base, op.EOF{}}); err != nil {
			t.Fatal(err)
		}
	}
}

func TestBlankNodePropertyList(t *testing.T) {
	for _, test := range []string{
		`[ foaf:name "Bob" ]`,
		`[
         foaf:name "Eve" ]`,
		`[
    	 foaf:name "Bob" ;
    	 foaf:knows [
		 foaf:name "Eve" ] ;
    	 foaf:mbox <bob@example.com> ]`,
		"[ :a :b ]",
		"[:p(<http://example/o>)]",
	} {
		p, err := NewParser([]rune(test))
		if err != nil {
			t.Fatal(err)
		}
		if _, err := p.Parse(op.And{BlankNodePropertyList, op.EOF{}}); err != nil {
			t.Fatal(err)
		}
	}
}

func TestCollection(t *testing.T) {
	for _, test := range []string{
		"()", "(())", "(()())", "(()(()))", "(()(()(())))",
		"( _:a )",
	} {
		p, err := NewParser([]rune(test))
		if err != nil {
			t.Fatal(err)
		}
		if _, err := p.Parse(op.And{Collection, op.EOF{}}); err != nil {
			t.Fatal(err)
		}
	}
}

func TestDocument(t *testing.T) {
	for _, test := range []string{
		`<http://example.org/#spiderman> <http://www.perceive.net/schemas/relationship/enemyOf> <http://example.org/#green-goblin> .`,
		`[] foaf:knows [ foaf:name "Bob" ] .`,
		`:subject :predicate ( _:a _:b _:c ) .`,
		`<>  rdf:type mf:Manifest ;
			 mf:name "N-Triples tests" ;
			 mf:entries
			 (
			 <#lantag_with_subtag>
			 <#minimal_whitespace>
			 ) .`,
		`<a> <b> <c>.`,
		"[ :a :b ] :c :d .",
		"<s> <p> 123.E+1 .",
		"[ :p (:o) ] .",
	} {
		p, err := NewParser([]rune(test))
		if err != nil {
			t.Fatal(err)
		}
		if _, err := p.Parse(op.And{Document, op.EOF{}}); err != nil {
			t.Fatal(err)
		}
	}
}

func TestExponent(t *testing.T) {
	for _, test := range []string{
		"e0", "e+0", "e-0", "E0", "E+0", "E-0", "E+1",
	} {
		p, err := NewParser([]rune(test))
		if err != nil {
			t.Fatal(err)
		}
		if _, err := p.Parse(op.And{Exponent, op.EOF{}}); err != nil {
			t.Fatal(err)
		}
	}
}

func TestIRI(t *testing.T) {
	for _, test := range []string{
		"<http://a.example/s>",
		"prefix:",
	} {
		p, err := NewParser([]rune(test))
		if err != nil {
			t.Fatal(err)
		}
		if _, err := p.Parse(op.And{IRI, op.EOF{}}); err != nil {
			t.Fatal(err)
		}
	}
}

func TestInteger(t *testing.T) {
	for _, test := range []string{
		"0", "99", "-1", "-99", "+1", "+99",
	} {
		p, err := NewParser([]rune(test))
		if err != nil {
			t.Fatal(err)
		}
		if _, err := p.Parse(op.And{Integer, op.EOF{}}); err != nil {
			t.Fatal(err)
		}
	}
	for _, test := range []string{
		"+", "-1.0",
	} {
		p, err := NewParser([]rune(test))
		if err != nil {
			t.Fatal(err)
		}
		if _, err := p.Parse(op.And{Integer, op.EOF{}}); err == nil {
			t.Fatal(test)
		}
	}
}

func TestLiteral(t *testing.T) {
	for _, test := range []string{
		`"Green Goblin"`,
		`"-1.0"^^<http://www.w3.org/2001/XMLSchema#decimal>`,
		"-1.0", "-123", "123.0", "123.E+1",
	} {
		p, err := NewParser([]rune(test))
		if err != nil {
			t.Fatal(err)
		}
		if _, err := p.Parse(op.And{Literal, op.EOF{}}); err != nil {
			t.Fatal(err)
		}
	}
}

func TestNumericLiteral(t *testing.T) {
	for _, test := range []string{
		"1", "+1", "-1", "+99",
		".1", "+.1", "-.1", "+.99", "1.1", "+1.1", "-1.1", "+1.99",
		"1e0", "+1e0", "-1e0", "+99e0", "0.E0", "+0.E0", "-0.E0", "+99.E0",
		"1e-1", "+1e+1", // etc
	} {
		p, err := NewParser([]rune(test))
		if err != nil {
			t.Fatal(err)
		}
		if _, err := p.Parse(op.And{Literal, op.EOF{}}); err != nil {
			t.Fatal(err)
		}
	}
}

func TestObjectList(t *testing.T) {
	for _, test := range []string{
		`"Green Goblin"`,
		`"Spiderman", "Человек-паук"@ru`,
		"<http://example.org/#green-goblin>",
		`[ foaf:name "Bob" ]`,
		`:a, :b, :c`,
		`<http://norman.walsh.name/knows/who/dan-brickley> ,
		 [ :mbox <mailto:timbl@w3.org> ] ,
		 <http://getopenid.com/amyvdh>`,
		"-1.0",
	} {
		p, err := NewParser([]rune(test))
		if err != nil {
			t.Fatal(err)
		}
		if _, err := p.Parse(op.And{ObjectList, op.EOF{}}); err != nil {
			t.Fatal(err)
		}
	}
}

func TestPNAME_NS(t *testing.T) {
	for _, test := range []string{
		"rdf:", "rdfs:", "foaf:", "rel:",
		"a:", "a.a:", "a::", ":", "a:a::a:",
		"a·̀ͯ‿.⁀:",
	} {
		p, err := NewParser([]rune(test))
		if err != nil {
			t.Fatal(err)
		}
		if _, err := p.Parse(op.And{PNAME_NS, op.EOF{}}); err != nil {
			t.Fatal(err)
		}
	}
	for _, test := range []string{
		"rdf", "rdfs", "foaf", "rel", "rdf:a",
		"invalid.",
	} {
		p, err := NewParser([]rune(test))
		if err != nil {
			t.Fatal(err)
		}
		if _, err := p.Parse(op.And{PNAME_NS, op.EOF{}}); err == nil {
			t.Fatal(test)
		}
	}
}

func TestPN_LOCAL(t *testing.T) {
	for _, test := range []string{
		"a·̀ͯ‿.⁀", ":", "a", "aaa", "enemyOf", "name", "2..0",
	} {
		p, err := NewParser([]rune(test))
		if err != nil {
			t.Fatal(err)
		}
		if _, err := p.Parse(op.And{PN_LOCAL, op.EOF{}}); err != nil {
			t.Fatal(err)
		}
	}
	for _, test := range []string{
		"2..0.",
	} {
		p, err := NewParser([]rune(test))
		if err != nil {
			t.Fatal(err)
		}
		if _, err := p.Parse(op.And{PN_LOCAL, op.EOF{}}); err == nil {
			t.Fatal(test)
		}
	}
}

func TestPredicateObjectList(t *testing.T) {
	for _, test := range []string{
		`foaf:name "Green Goblin"`,
		`rel:enemyOf <#spiderman> ;
		 a foaf:Person ;    # in the context of the Marvel universe
    	 foaf:name "Green Goblin"`,
		`rel:enemyOf <#green-goblin> ;
    	 a foaf:Person ;
    	 foaf:name "Spiderman", "Человек-паук"@ru`,
		`foaf:name "Bob" ;
		 foaf:knows [
         foaf:name "Eve" ] ;
    	 foaf:mbox <bob@example.com>`,
	} {
		p, err := NewParser([]rune(test))
		if err != nil {
			t.Fatal(err)
		}
		if _, err := p.Parse(op.And{PredicateObjectList, op.EOF{}}); err != nil {
			t.Fatal(err)
		}
	}
}

func TestPrefix(t *testing.T) {
	for _, test := range []string{
		"@prefix rdf: <http://www.w3.org/1999/02/22-rdf-syntax-ns#> .",
		"@prefix rdfs: <http://www.w3.org/2000/01/rdf-schema#> .",
		"@prefix foaf: <http://xmlns.com/foaf/0.1/> .",
		"@prefix rel: <http://www.perceive.net/schemas/relationship/> .",
		"@prefix a·̀ͯ‿.⁀: <http://a.example/>.",
	} {
		p, err := NewParser([]rune(test))
		if err != nil {
			t.Fatal(err)
		}
		if _, err := p.Parse(op.And{PrefixID, op.EOF{}}); err != nil {
			t.Fatal(err)
		}
	}
}

func TestPrefixedName(t *testing.T) {
	for _, test := range []string{
		"p:a·̀ͯ‿.⁀", "p:", "p:p", ":0.1", ":0..2",
	} {
		p, err := NewParser([]rune(test))
		if err != nil {
			t.Fatal(err)
		}
		if _, err := p.Parse(op.And{PrefixedName, op.EOF{}}); err != nil {
			t.Fatal(err)
		}
	}
	for _, test := range []string{
		"invalid.:o",
	} {
		p, err := NewParser([]rune(test))
		if err != nil {
			t.Fatal(err)
		}
		if _, err := p.Parse(op.And{PrefixedName, op.EOF{}}); err == nil {
			t.Fatal(test)
		}
	}
}

func TestString(t *testing.T) {
	for _, test := range []string{
		`'''TEST'''`, `'"TEST"'`, `'''('')'''`,
		"'''TEST\nTEST # TEST\nTEST'''",
	} {
		p, err := NewParser([]rune(test))
		if err != nil {
			t.Fatal(err)
		}
		if _, err := p.Parse(op.And{String, op.EOF{}}); err != nil {
			t.Fatal(err)
		}
	}
}

func TestVerb(t *testing.T) {
	for _, test := range []string{
		"rel:enemyOf", "a", "foaf:name",
	} {
		p, err := NewParser([]rune(test))
		if err != nil {
			t.Fatal(err)
		}
		if _, err := p.Parse(op.And{Verb, op.EOF{}}); err != nil {
			t.Fatal(err)
		}
	}
}
