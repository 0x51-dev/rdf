package grammar_test

import (
	trig "github.com/0x51-dev/rdf/trig/grammar"
	"github.com/0x51-dev/upeg/parser/op"
	"testing"
)

func TestBlock(t *testing.T) {
	for _, test := range []string{
		`{
	<http://example.org/bob> dc:publisher "Bob" .
	<http://example.org/alice> dc:publisher "Alice" .
}`,
		"[] {<http://a.example/s> <http://a.example/p> <http://a.example/o> .}",
		"{_:s:p :o .}",
		"{[:p(:o)].}",
		"<http://a.example/s> <http://a.example/p> p:o#comment\n.",
		`<http://example/graph> {
  <http://a.example/s> <http://a.example/p> <http://a.example/o>#comment
  .
}`,
	} {
		p, err := trig.NewParser([]rune(test))
		if err != nil {
			t.Fatal(err)
		}
		if _, err := p.Parse(op.And{trig.Block, op.EOF{}}); err != nil {
			t.Fatal(err)
		}
	}
}
