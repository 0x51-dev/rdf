package grammar_test

import (
	. "github.com/0x51-dev/rdf/nquads/grammar"
	"github.com/0x51-dev/upeg/parser"
	"github.com/0x51-dev/upeg/parser/op"
	"testing"
)

func TestQuad(t *testing.T) {
	for _, test := range []string{
		"<http://one.example/subject1> <http://one.example/predicate1> <http://one.example/object1> <http://example.org/graph3> . # comments here",
		"_:subject1 <http://an.example/predicate1> \"object1\" <http://example.org/graph1> .",
		"_:subject2 <http://an.example/predicate2> \"object2\" <http://example.org/graph5> .",
		"<http://example.org/#spiderman> <http://www.perceive.net/schemas/relationship/enemyOf> <http://example.org/#green-goblin> <http://example.org/graphs/spiderman> .",
		"_:alice <http://xmlns.com/foaf/0.1/knows> _:bob <http://example.org/graphs/john> .",
		"_:bob <http://xmlns.com/foaf/0.1/knows> _:alice <http://example.org/graphs/james> .",
	} {
		p, err := parser.New([]rune(test))
		if err != nil {
			t.Fatal(err)
		}
		if _, err := p.Parse(op.And{Statement, op.EOF{}}); err != nil {
			t.Fatal(err)
		}
	}
}
