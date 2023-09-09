package turtle

import (
	"fmt"
	"github.com/0x51-dev/rdf/ntriples"
	"github.com/0x51-dev/rdf/ntriples/grammar"
	"github.com/0x51-dev/upeg/parser"
	"github.com/0x51-dev/upeg/parser/op"
	"strconv"
	"strings"
)

func (ctx *context) evaluateDocument(d Document) (ntriples.Document, error) {
	for _, t := range d {
		switch t := t.(type) {
		case *Base:
			ctx.base = string(*t)
		case *Prefix:
			ctx.prefixes[t.Name] = t.IRI
		case *Triple:
			var subject ntriples.Subject
			if t.Subject != nil {
				switch t := t.Subject.(type) {
				case *IRI:
					s, err := ctx.evaluateIRI(t)
					if err != nil {
						return nil, err
					}
					subject = s
				case *BlankNode:
					bn := ntriples.BlankNode(*t)
					subject = &bn
				case Collection:
					var objects []ntriples.Object
					for _, o := range t {
						o, err := ctx.evaluateObject(o)
						if err != nil {
							return nil, err
						}
						objects = append(objects, o...)
					}
					if len(objects) == 0 {
						return nil, nil
					}

					var el ntriples.BlankNode
					for i, o := range objects {
						e := ctx.el()
						ctx.triples = append(ctx.triples, ntriples.Triple{
							Subject:   &e,
							Predicate: "http://www.w3.org/1999/02/22-rdf-syntax-ns#first",
							Object:    o,
						})
						if i+1 != len(objects) {
							el = ctx.el()
							ctx.triples = append(ctx.triples, ntriples.Triple{
								Subject:   &e,
								Predicate: "http://www.w3.org/1999/02/22-rdf-syntax-ns#rest",
								Object:    &el,
							})
						} else {
							el = e
						}
					}
					ctx.triples = append(ctx.triples, ntriples.Triple{
						Subject:   &el,
						Predicate: "http://www.w3.org/1999/02/22-rdf-syntax-ns#rest",
						Object:    ntriples.IRIReference("http://www.w3.org/1999/02/22-rdf-syntax-ns#nil"),
					})
				default:
					panic(fmt.Errorf("unknown subject type %T", t))
				}

				for _, po := range t.PredicateObjectList {
					p, os, err := ctx.evaluatePredicateObject(po)
					if err != nil {
						return nil, err
					}
					for _, o := range os {
						ctx.triples = append(ctx.triples, ntriples.Triple{
							Subject:   subject,
							Predicate: p,
							Object:    o,
						})
					}
				}
			} else {
				bn := ctx.bn()
				for _, po := range t.BlankNodePropertyList {
					p, os, err := ctx.evaluatePredicateObject(po)
					if err != nil {
						return nil, err
					}
					for _, o := range os {
						ctx.triples = append(ctx.triples, ntriples.Triple{
							Subject:   &bn,
							Predicate: p,
							Object:    o,
						})
					}
				}
				for _, po := range t.PredicateObjectList {
					p, os, err := ctx.evaluatePredicateObject(po)
					if err != nil {
						return nil, err
					}
					for _, o := range os {
						ctx.triples = append(ctx.triples, ntriples.Triple{
							Subject:   &bn,
							Predicate: p,
							Object:    o,
						})
					}
				}
			}
		default:
			panic(fmt.Errorf("unknown document type %T", t))
		}
	}
	return ctx.triples, nil
}

func (ctx *context) evaluateIRI(iri *IRI) (*ntriples.IRIReference, error) {
	if !iri.Prefixed {
		v := iri.Value

		// Validate IRI.
		if strings.Contains(v, "\\u") || strings.Contains(v, "\\U") {
			// Unescape unicode characters.
			if v_, err := strconv.Unquote(`"` + v + `"`); err == nil {
				v = v_
			}
		}
		p, err := parser.New([]rune(`<` + v + `>`))
		if err != nil {
			return nil, err
		}
		if _, err := p.Match(op.And{grammar.IRIReference, op.EOF{}}); err != nil {
			return nil, err
		}

		ref := ntriples.IRIReference(v)
		return &ref, nil
	}

	p := strings.SplitAfterN(iri.Value, ":", 2)
	if len(p) != 2 {
		return nil, fmt.Errorf("invalid prefixed IRI %q", iri.Value)
	}
	prefix, ok := ctx.prefixes[p[0]]
	if !ok {
		return nil, fmt.Errorf("prefix %q not defined", p[0])
	}
	ref := ntriples.IRIReference(prefix + p[1])
	return &ref, nil
}

func (ctx *context) evaluateObject(o Object) ([]ntriples.Object, error) {
	switch o := o.(type) {
	case *IRI:
		i, err := ctx.evaluateIRI(o)
		if err != nil {
			return nil, err
		}
		return []ntriples.Object{*i}, nil
	case *BlankNode:
		bn := ntriples.BlankNode(*o)
		return []ntriples.Object{&bn}, nil
	case BlankNodePropertyList:
		var objects []ntriples.Object
		for _, op := range o {
			bn := ctx.bn()
			p, os, err := ctx.evaluatePredicateObject(op)
			if err != nil {
				return nil, err
			}
			for _, o := range os {
				ctx.triples = append(ctx.triples, ntriples.Triple{
					Subject:   &bn,
					Predicate: p,
					Object:    o,
				})
			}
			objects = append(objects, &bn)
		}
		return objects, nil
	case Collection:
		var objects []ntriples.Object
		for _, o := range o {
			o, err := ctx.evaluateObject(o)
			if err != nil {
				return nil, err
			}
			objects = append(objects, o...)
		}
		if len(objects) == 0 {
			o := ntriples.IRIReference("http://www.w3.org/1999/02/22-rdf-syntax-ns#nil")
			return []ntriples.Object{&o}, nil
		}

		var first, el ntriples.BlankNode
		for i, o := range objects {
			e := ctx.el()
			if i == 0 {
				first = e
			}
			ctx.triples = append(ctx.triples, ntriples.Triple{
				Subject:   &e,
				Predicate: "http://www.w3.org/1999/02/22-rdf-syntax-ns#first",
				Object:    o,
			})
			if i+1 != len(objects) {
				el = ctx.el()
				ctx.triples = append(ctx.triples, ntriples.Triple{
					Subject:   &e,
					Predicate: "http://www.w3.org/1999/02/22-rdf-syntax-ns#rest",
					Object:    &el,
				})
			} else {
				el = e
			}
		}
		ctx.triples = append(ctx.triples, ntriples.Triple{
			Subject:   &el,
			Predicate: "http://www.w3.org/1999/02/22-rdf-syntax-ns#rest",
			Object:    ntriples.IRIReference("http://www.w3.org/1999/02/22-rdf-syntax-ns#nil"),
		})
		return []ntriples.Object{&first}, nil
	case *NumericLiteral:
		var ref ntriples.IRIReference
		switch o.Type {
		case Integer:
			ref = "http://www.w3.org/2001/XMLSchema#integer"
		case Decimal:
			ref = "http://www.w3.org/2001/XMLSchema#decimal"
		case Double:
			ref = "http://www.w3.org/2001/XMLSchema#double"
		default:
			panic(fmt.Errorf("unknown numeric literal type %q", o.Type))
		}
		return []ntriples.Object{&ntriples.Literal{
			Value:     o.Value,
			Reference: &ref,
		}}, nil
	case *StringLiteral:
		v := o.Value
		if o.Multiline {
			v = strings.ReplaceAll(v, "\n", "\\n")
		}
		if o.DatatypeIRI != "" {
			ref := ntriples.IRIReference(o.DatatypeIRI)
			return []ntriples.Object{&ntriples.Literal{
				Value:     v,
				Reference: &ref,
			}}, nil
		}
		if o.LanguageTag != "" {
			return []ntriples.Object{&ntriples.Literal{
				Value:    v,
				Language: o.LanguageTag,
			}}, nil
		}
		return []ntriples.Object{&ntriples.Literal{
			Value: v,
		}}, nil
	case *BooleanLiteral:
		ref := ntriples.IRIReference("http://www.w3.org/2001/XMLSchema#boolean")
		return []ntriples.Object{&ntriples.Literal{
			Value:     o.String(),
			Reference: &ref,
		}}, nil
	default:
		panic(fmt.Errorf("unknown objects type %T", o))
	}
}

func (ctx *context) evaluateObjectList(os ObjectList) ([]ntriples.Object, error) {
	var objects []ntriples.Object
	for _, o := range os {
		os, err := ctx.evaluateObject(o)
		if err != nil {
			return nil, err
		}
		objects = append(objects, os...)
	}
	return objects, nil
}

func (ctx *context) evaluatePredicateObject(po PredicateObject) (ntriples.IRIReference, []ntriples.Object, error) {
	var predicate ntriples.IRIReference
	switch v := po.Verb.(type) {
	case *IRI:
		p, err := ctx.evaluateIRI(v)
		if err != nil {
			return "", nil, err
		}
		predicate = *p
	}
	objects, err := ctx.evaluateObjectList(po.ObjectList)
	if err != nil {
		return "", nil, err
	}
	return predicate, objects, nil
}
