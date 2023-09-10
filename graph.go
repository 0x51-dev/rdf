package rdf

type Graph struct {
	triples map[*Triple]struct{}
}

func NewGraph(ts ...*Triple) *Graph {
	triples := make(map[*Triple]struct{})
	for _, t := range ts {
		triples[t] = struct{}{}
	}
	return &Graph{triples}
}

func (g *Graph) Add(s, p, o Node) {
	t := &Triple{s, p, o}
	g.triples[t] = struct{}{}
}

// Find returns the first triple matching the given pattern.
func (g *Graph) Find(s, p, o Node) *Triple {
	for t := range g.iter() {
		if (s == nil || t.Subject == s) &&
			(p == nil || t.Predicate == p) &&
			(o == nil || t.Object == o) {
			return t
		}
	}
	return nil
}

// FindAll returns all triples matching the given pattern.
func (g *Graph) FindAll(s, p, o Node) []*Triple {
	var triples []*Triple
	for t := range g.iter() {
		if (s == nil || t.Subject == s) &&
			(p == nil || t.Predicate == p) &&
			(o == nil || t.Object == o) {
			triples = append(triples, t)
		}
	}
	return triples
}

func (g *Graph) Remove(t *Triple) {
	delete(g.triples, t)
}

func (g *Graph) iter() chan *Triple {
	ch := make(chan *Triple)
	go func() {
		for t := range g.triples {
			ch <- t
		}
		close(ch)
	}()
	return ch
}

type Triple struct {
	Subject   Node
	Predicate Node
	Object    Node
}

func NewTriple(s, p, o Node) *Triple {
	return &Triple{s, p, o}
}

func (t *Triple) Equal(other *Triple) bool {
	return t.Subject.Equal(other.Subject) && t.Predicate.Equal(other.Predicate) && t.Object.Equal(other.Object)
}
