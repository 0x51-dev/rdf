package rdf

type Graph struct {
	triples []*Triple
}

func NewGraph(ts ...*Triple) *Graph {
	return &Graph{triples: ts}
}

func (g *Graph) Add(s, p, o Node) {
	t := &Triple{s, p, o}
	g.triples = append(g.triples, t)
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
	for i, other := range g.triples {
		if t.Equal(other) {
			g.triples = append(g.triples[:i], g.triples[i+1:]...)
			return
		}
	}
}

func (g *Graph) iter() chan *Triple {
	ch := make(chan *Triple)
	go func() {
		for _, t := range g.triples {
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
