package rdf

import "testing"

func TestGraph(t *testing.T) {
	var (
		a = &BlankNode{"a"}
		b = &BlankNode{"b"}
		c = &BlankNode{"c"}
	)

	abc := NewTriple(a, b, c)
	abb := NewTriple(a, b, b)
	aaa := NewTriple(a, a, a)

	g := NewGraph(abc, abb, aaa)

	if triple := g.Find(nil, nil, nil); triple != abc {
		t.Error("expected abc")
	}
	if triple := g.Find(a, b, c); triple != abc {
		t.Error("expected abc")
	}
	if triple := g.Find(a, a, a); t == nil || triple != aaa {
		t.Error("expected aaa")
	}
	if triple := g.Find(nil, a, nil); triple != aaa {
		t.Error("expected aaa")
	}

	if len(g.FindAll(a, nil, nil)) != 3 {
		t.Error("expected 3 triples")
	}
	if len(g.FindAll(nil, b, nil)) != 2 {
		t.Error("expected 2 triples")
	}
}
