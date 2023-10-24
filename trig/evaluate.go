package trig

import (
	"fmt"
	nq "github.com/0x51-dev/rdf/nquads"
	nt "github.com/0x51-dev/rdf/ntriples"
	ttl "github.com/0x51-dev/rdf/turtle"
)

func (ctx *Context) evaluateDocument(d Document) (nq.Document, error) {
	var triples []nq.Quad
	for _, t := range d {
		switch t := t.(type) {
		case *Base:
			ctx.Context.Base = string(*t)
		case *Prefix:
			ctx.Context.Prefixes[t.Name] = t.IRI
		case *TriplesOrGraph:
			switch los := t.LabelOrSubject.(type) {
			case *IRI:
				if len(t.WrappedGraph) != 0 {
					graphLabel, err := ctx.EvaluateIRI((*ttl.IRI)(los))
					if err != nil {
						return nil, err
					}
					for _, t := range t.WrappedGraph {
						ts, err := ctx.EvaluateTriple(&t)
						if err != nil {
							return nil, err
						}
						for _, t := range ts {
							triples = append(triples, nq.NewQuadFromTriple(t, graphLabel))
						}
					}
				} else {
					subject, err := ctx.EvaluateIRI((*ttl.IRI)(los))
					if err != nil {
						return nil, err
					}
					for _, po := range t.PredicateObjectList {
						p, os, ts, err := ctx.EvaluatePredicateObject(po)
						if err != nil {
							return nil, err
						}
						for _, t := range ts {
							triples = append(triples, nq.NewQuadFromTriple(t, nil))
						}
						for _, o := range os {
							triples = append(triples, nq.NewQuadFromTriple(
								nt.Triple{
									Subject:   subject,
									Predicate: p,
									Object:    o,
								}, nil,
							))
						}
					}
				}
			case *BlankNode:
				if *los == "[]" {
					bn := ctx.bn()
					if len(t.WrappedGraph) != 0 {
						for _, t := range t.WrappedGraph {
							ts, err := ctx.EvaluateTriple(&t)
							if err != nil {
								return nil, err
							}
							for _, t := range ts {
								triples = append(triples, nq.NewQuadFromTriple(t, &bn))
							}
						}
					} else {
						for _, po := range t.PredicateObjectList {
							p, os, ts, err := ctx.EvaluatePredicateObject(po)
							if err != nil {
								return nil, err
							}
							for _, t := range ts {
								triples = append(triples, nq.NewQuadFromTriple(t, nil))
							}
							for _, o := range os {
								triples = append(triples, nq.NewQuadFromTriple(
									nt.Triple{
										Subject:   &bn,
										Predicate: p,
										Object:    o,
									}, nil,
								))
							}
						}
					}
				} else {
					if len(t.WrappedGraph) != 0 {
						for _, t := range t.WrappedGraph {
							ts, err := ctx.EvaluateTriple(&t)
							if err != nil {
								return nil, err
							}
							for _, t := range ts {
								triples = append(triples, nq.NewQuadFromTriple(t, (*nt.BlankNode)(los)))
							}
						}
					} else {
						for _, po := range t.PredicateObjectList {
							p, os, ts, err := ctx.EvaluatePredicateObject(po)
							if err != nil {
								return nil, err
							}
							for _, t := range ts {
								triples = append(triples, nq.NewQuadFromTriple(t, nil))
							}
							for _, o := range os {
								triples = append(triples, nq.NewQuadFromTriple(
									nt.Triple{
										Subject:   (*nt.BlankNode)(los),
										Predicate: p,
										Object:    o,
									}, nil,
								))
							}
						}
					}
				}
			default:
				return nil, fmt.Errorf("unknown label or subject type %T", los)
			}
		case WrappedGraph:
			for _, t := range t {
				ts, err := ctx.EvaluateTriple(&t)
				if err != nil {
					return nil, err
				}
				for _, t := range ts {
					triples = append(triples, nq.NewQuadFromTriple(t, nil))
				}
			}
		case *Triple2:
			if len(t.BlankNodePropertyList) != 0 {
				bn := ctx.bn()
				for _, po := range t.BlankNodePropertyList {
					p, os, ts, err := ctx.EvaluatePredicateObject(po)
					if err != nil {
						return nil, err
					}
					for _, t := range ts {
						triples = append(triples, nq.NewQuadFromTriple(t, nil))
					}
					for _, o := range os {
						triples = append(triples, nq.NewQuadFromTriple(
							nt.Triple{
								Subject:   &bn,
								Predicate: p,
								Object:    o,
							}, nil,
						))
					}
				}
				for _, po := range t.PredicateObjectList {
					p, os, ts, err := ctx.EvaluatePredicateObject(po)
					if err != nil {
						return nil, err
					}
					for _, t := range ts {
						triples = append(triples, nq.NewQuadFromTriple(t, nil))
					}
					for _, o := range os {
						triples = append(triples, nq.NewQuadFromTriple(
							nt.Triple{
								Subject:   &bn,
								Predicate: p,
								Object:    o,
							}, nil,
						))
					}
				}
			} else {
				var subject nt.Subject
				var objects []nt.Object
				for _, o := range t.Collection {
					os, ts, err := ctx.EvaluateObject(o)
					if err != nil {
						return nil, err
					}
					objects = append(objects, os...)
					for _, t := range ts {
						triples = append(triples, nq.NewQuadFromTriple(t, nil))
					}
				}
				if len(objects) == 0 {
					return triples, nil
				}

				var el nt.BlankNode
				for i, o := range objects {
					e := ctx.el()
					if i == 0 {
						subject = e
					}
					triples = append(triples, nq.NewQuadFromTriple(
						nt.Triple{
							Subject:   &e,
							Predicate: "http://www.w3.org/1999/02/22-rdf-syntax-ns#first",
							Object:    o,
						}, nil,
					))
					if i+1 != len(objects) {
						el = ctx.el()
						triples = append(triples, nq.NewQuadFromTriple(
							nt.Triple{
								Subject:   &e,
								Predicate: "http://www.w3.org/1999/02/22-rdf-syntax-ns#rest",
								Object:    &el,
							}, nil,
						))
					} else {
						el = e
					}
				}
				triples = append(triples, nq.NewQuadFromTriple(
					nt.Triple{
						Subject:   &el,
						Predicate: "http://www.w3.org/1999/02/22-rdf-syntax-ns#rest",
						Object:    nt.IRIReference("http://www.w3.org/1999/02/22-rdf-syntax-ns#nil"),
					}, nil,
				))

				for _, po := range t.PredicateObjectList {
					p, os, ts, err := ctx.EvaluatePredicateObject(po)
					if err != nil {
						return nil, err
					}
					for _, t := range ts {
						triples = append(triples, nq.NewQuadFromTriple(t, nil))
					}
					for _, o := range os {
						triples = append(triples, nq.NewQuadFromTriple(
							nt.Triple{
								Subject:   subject,
								Predicate: p,
								Object:    o,
							}, nil,
						))
					}
				}
			}
		default:
			return nil, fmt.Errorf("unknown document type %T", t)
		}
	}
	return triples, nil
}
