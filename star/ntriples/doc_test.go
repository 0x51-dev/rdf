package ntriples_test

import (
	nts "github.com/0x51-dev/rdf/star/ntriples"
	"testing"
)

// MORE TESTS: https://w3c.github.io/rdf-star/tests/nt/syntax/manifest.html
func TestParseDocument(t *testing.T) {
	if _, err := nts.ParseDocument("<< <http://example/s> <http://example/p> <http://example/o> >> <http://example/q> <http://example/z> ."); err != nil {
		t.Fatal(err)
	}
}
