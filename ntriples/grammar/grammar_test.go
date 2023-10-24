package grammar_test

import (
	_ "embed"
	. "github.com/0x51-dev/rdf/ntriples/grammar"
	"github.com/0x51-dev/upeg/parser"
	"github.com/0x51-dev/upeg/parser/op"
	"testing"
)

func TestBlankNodeLabel(t *testing.T) {
	for _, test := range []string{
		"_:0",
		"_:a",
		"_:aa",
		"_:a.a",
		"_:a-",
		"_:a.....1",
		"_:anon",
	} {
		p, err := parser.New([]rune(test))
		if err != nil {
			t.Fatal(err)
		}
		if _, err := p.Parse(op.And{BlankNodeLabel, op.EOF{}}); err != nil {
			t.Fatal(test, err)
		}
	}
}

func TestObject(t *testing.T) {
	for _, test := range []string{
		"<http://one.example/object1>",
		"_:object1",
	} {
		p, err := parser.New([]rune(test))
		if err != nil {
			t.Fatal(err)
		}
		if _, err := p.Parse(op.And{Object, op.EOF{}}); err != nil {
			t.Fatal(err)
		}
	}
}

func TestPN_CHARS(t *testing.T) {
	for _, char := range []rune{
		'A', 'Z', 'a', 'z', '0', '9',
	} {
		p, err := parser.New([]rune{char})
		if err != nil {
			t.Fatal(err)
		}
		if _, err := p.Parse(op.And{PN_CHARS, op.EOF{}}); err != nil {
			t.Fatal(err)
		}
	}
	for _, char := range []rune{
		'.', ':',
	} {
		p, err := parser.New([]rune{char})
		if err != nil {
			t.Fatal(err)
		}
		if _, err := p.Parse(op.And{PN_CHARS, op.EOF{}}); err == nil {
			t.Fatal(char)
		}
	}
}

func TestPN_CHARS_BASE(t *testing.T) {
	for _, char := range []rune{
		'A', 'Z', 'a', 'z',
	} {
		p, err := parser.New([]rune{char})
		if err != nil {
			t.Fatal(err)
		}
		if _, err := p.Parse(op.And{PN_CHARS_BASE, op.EOF{}}); err != nil {
			t.Fatal(err)
		}
	}
	for _, char := range []rune{
		':',
	} {
		p, err := parser.New([]rune{char})
		if err != nil {
			t.Fatal(err)
		}
		if _, err := p.Parse(op.And{PN_CHARS_BASE, op.EOF{}}); err == nil {
			t.Fatal(char)
		}
	}
}

func TestPredicate(t *testing.T) {
	for _, test := range []string{
		"<http://one.example/predicate1>",
	} {
		p, err := parser.New([]rune(test))
		if err != nil {
			t.Fatal(err)
		}
		if _, err := p.Parse(op.And{Predicate, op.EOF{}}); err != nil {
			t.Fatal(err)
		}
	}
}

func TestSubject(t *testing.T) {
	for _, test := range []string{
		"<http://one.example/subject1>",
		"_:subject1",
		"_:subject2",
	} {
		p, err := parser.New([]rune(test))
		if err != nil {
			t.Fatal(err)
		}
		if _, err := p.Parse(op.And{Subject, op.EOF{}}); err != nil {
			t.Fatal(err)
		}
	}
}

func TestTriple(t *testing.T) {
	for _, test := range []string{
		"<http://one.example/subject1> <http://one.example/predicate1> <http://one.example/object1> . # comments here",
		"_:subject1 <http://an.example/predicate1> \"object1\" .",
		"_:subject2 <http://an.example/predicate2> \"object2\" .",
		"<http://example.org/#spiderman> <http://www.perceive.net/schemas/relationship/enemyOf> <http://example.org/#green-goblin> .",
		"<http://example.org/resource15> <http://example.org/property> _:anon.",
	} {
		p, err := parser.New([]rune(test))
		if err != nil {
			t.Fatal(err)
		}
		if _, err := p.Parse(op.And{Triple, op.EOF{}}); err != nil {
			t.Fatal(err)
		}
	}
}
