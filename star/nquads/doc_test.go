package nquads_test

import (
	nqs "github.com/0x51-dev/rdf/star/nquads"
	"testing"
)

// MORE TESTS: https://w3c.github.io/rdf-star/tests/nt/syntax/manifest.html
func TestParseDocument(t *testing.T) {
	if _, err := nqs.ParseDocument("<< <http://example/s> <http://example/p> <http://example/o> >> <http://example/q> <http://example/z> ."); err != nil {
		t.Fatal(err)
	}
}
