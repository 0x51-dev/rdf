package grammar

import (
	"github.com/0x51-dev/upeg/parser/op"
	"testing"
)

func TestStatement(t *testing.T) {
	for _, test := range []string{
		"@prefix : <http://example.org/>",
		"@prefix foaf:<http://xmlns.com/foaf/0.1/>",
		"<http://example.org/spiderman> <http://example.org/enemyOf> <http://example.org/green-goblin>",
		`<http://example.org/spiderman> <http://xmlns.com/foaf/0.1/name> "Spiderman", "Peter Parker"`,
		":spiderman :enemyOf [ id :green-goblin :portrayedBy [ id :willem-dafoe a :Actor ] ] ; :portrayedBy [ id :tobey-maguire a :Actor ]",
		"{ :weather a :Raining } => { :weather a :Cloudy }",
		"{ ?x a :SuperHero } => { ?x a :Imaginary }",
	} {
		p, err := NewParser([]rune(test))
		if err != nil {
			t.Fatal(err)
		}
		if _, err := p.Parse(op.And{Statement, op.EOF{}}); err != nil {
			t.Fatal(err)
		}
	}
}
