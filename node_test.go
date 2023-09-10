package rdf

import (
	"encoding/json"
	"fmt"
	"testing"
)

const (
	exampleURI       = "https://example.org/"
	exampleBlankNode = "_:example"
)

func TestFromObject_blankNode(t *testing.T) {
	for _, test := range []string{
		fmt.Sprintf("%q", exampleBlankNode),
		fmt.Sprintf(`{ "@id": %q }`, exampleBlankNode),
	} {
		var m any
		if err := json.Unmarshal([]byte(test), &m); err != nil {
			t.Fatal(err)
		}
		n, err := fromObject(m, false)
		if err != nil {
			t.Fatal(err)
		}
		bn, ok := n.(*BlankNode)
		if !ok {
			t.Errorf("expected BlankNode, got %T", n)
		}
		if bn.Attribute != exampleBlankNode {
			t.Errorf("expected %q, got %q", exampleBlankNode, bn.Attribute)
		}
	}
}

func TestFromObject_iriReference(t *testing.T) {
	for _, test := range []string{
		fmt.Sprintf("%q", exampleURI),
		fmt.Sprintf(`{ "@id": %q }`, exampleURI),
	} {
		var m any
		if err := json.Unmarshal([]byte(test), &m); err != nil {
			t.Fatal(err)
		}
		n, err := fromObject(m, false)
		if err != nil {
			t.Fatal(err)
		}
		r, ok := n.(*IRIReference)
		if !ok {
			t.Errorf("expected IRIReference, got %T", n)
		}
		if r.Value != exampleURI {
			t.Errorf("expected %q, got %q", exampleURI, r.Value)
		}
	}
}
