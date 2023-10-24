package turtle

import (
	"fmt"
	nt "github.com/0x51-dev/rdf/ntriples"
	"github.com/0x51-dev/rdf/ntriples/grammar"
	"github.com/0x51-dev/upeg/parser"
	"github.com/0x51-dev/upeg/parser/op"
	"strconv"
	"strings"
)

func (ctx *Context) EvaluateBlankNodePropertyList(pl BlankNodePropertyList) ([]nt.Object, []nt.Triple, error) {
	var objects []nt.Object
	var triples []nt.Triple
	for _, n := range pl {
		bn := ctx.bn()
		p, os, ts, err := ctx.EvaluatePredicateObject(n)
		if err != nil {
			return nil, nil, err
		}
		triples = append(triples, ts...)
		for _, o := range os {
			triples = append(triples, nt.Triple{
				Subject:   &bn,
				Predicate: p,
				Object:    o,
			})
		}
		objects = append(objects, &bn)
	}
	return objects, triples, nil
}

func (ctx *Context) EvaluateBooleanLiteral(o *BooleanLiteral) (*nt.Literal, error) {
	ref := nt.IRIReference("http://www.w3.org/2001/XMLSchema#boolean")
	return &nt.Literal{
		Value:     o.String(),
		Reference: &ref,
	}, nil
}

func (ctx *Context) EvaluateCollection(c Collection) (nt.Object, []nt.Triple, error) {
	var objects []nt.Object
	var triples []nt.Triple
	for _, o := range c {
		o, ts, err := ctx.EvaluateObject(o)
		if err != nil {
			return nil, nil, err
		}
		objects = append(objects, o...)
		triples = append(triples, ts...)
	}
	if len(objects) == 0 {
		o := nt.IRIReference("http://www.w3.org/1999/02/22-rdf-syntax-ns#nil")
		return &o, triples, nil
	}

	var first, el nt.BlankNode
	for i, o := range objects {
		e := ctx.el()
		if i == 0 {
			first = e
		}
		triples = append(triples, nt.Triple{
			Subject:   &e,
			Predicate: "http://www.w3.org/1999/02/22-rdf-syntax-ns#first",
			Object:    o,
		})
		if i+1 != len(objects) {
			el = ctx.el()
			triples = append(triples, nt.Triple{
				Subject:   &e,
				Predicate: "http://www.w3.org/1999/02/22-rdf-syntax-ns#rest",
				Object:    &el,
			})
		} else {
			el = e
		}
	}
	triples = append(triples, nt.Triple{
		Subject:   &el,
		Predicate: "http://www.w3.org/1999/02/22-rdf-syntax-ns#rest",
		Object:    nt.IRIReference("http://www.w3.org/1999/02/22-rdf-syntax-ns#nil"),
	})
	return &first, triples, nil
}

func (ctx *Context) EvaluateIRI(iri *IRI) (*nt.IRIReference, error) {
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

		ref := nt.IRIReference(strings.ReplaceAll(v, "\\", ""))
		return &ref, nil
	}

	p := strings.SplitAfterN(iri.Value, ":", 2)
	if len(p) != 2 {
		return nil, fmt.Errorf("invalid prefixed IRI %q", iri.Value)
	}
	prefix, ok := ctx.Prefixes[p[0]]
	if !ok {
		return nil, fmt.Errorf("prefix %q not defined", p[0])
	}
	ref := nt.IRIReference(prefix + strings.ReplaceAll(p[1], "\\", ""))
	return &ref, nil
}

func (ctx *Context) EvaluateNumericLiteral(o *NumericLiteral) (*nt.Literal, error) {
	var ref nt.IRIReference
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
	return &nt.Literal{
		Value:     o.Value,
		Reference: &ref,
	}, nil
}

func (ctx *Context) EvaluateObject(o Object) ([]nt.Object, []nt.Triple, error) {
	switch o := o.(type) {
	case *IRI:
		i, err := ctx.EvaluateIRI(o)
		if err != nil {
			return nil, nil, err
		}
		return []nt.Object{*i}, nil, nil
	case *BlankNode:
		if *o == "[]" {
			bn := ctx.bn()
			return []nt.Object{&bn}, nil, nil
		}
		bn := nt.BlankNode(*o)
		return []nt.Object{&bn}, nil, nil
	case BlankNodePropertyList:
		return ctx.EvaluateBlankNodePropertyList(o)
	case Collection:
		obj, ts, err := ctx.EvaluateCollection(o)
		if err != nil {
			return nil, nil, err
		}
		return []nt.Object{obj}, ts, nil
	case *NumericLiteral:
		obj, err := ctx.EvaluateNumericLiteral(o)
		if err != nil {
			return nil, nil, err
		}
		return []nt.Object{obj}, nil, nil
	case *StringLiteral:
		obj, err := ctx.EvaluateStringLiteral(o)
		if err != nil {
			return nil, nil, err
		}
		return []nt.Object{obj}, nil, nil
	case *BooleanLiteral:
		obj, err := ctx.EvaluateBooleanLiteral(o)
		if err != nil {
			return nil, nil, err
		}
		return []nt.Object{obj}, nil, nil
	default:
		panic(fmt.Errorf("unknown objects type %T", o))
	}
}

// EvaluateObjectList evaluates a list of objects.
func (ctx *Context) EvaluateObjectList(os ObjectList) ([]nt.Object, []nt.Triple, error) {
	var objects []nt.Object
	var triples []nt.Triple
	for _, o := range os {
		os, ts, err := ctx.EvaluateObject(o)
		if err != nil {
			return nil, nil, err
		}
		objects = append(objects, os...)
		triples = append(triples, ts...)
	}
	return objects, triples, nil
}

func (ctx *Context) EvaluatePredicateObject(po PredicateObject) (nt.IRIReference, []nt.Object, []nt.Triple, error) {
	var predicate nt.IRIReference
	switch v := po.Verb.(type) {
	case *IRI:
		p, err := ctx.EvaluateIRI(v)
		if err != nil {
			return "", nil, nil, err
		}
		predicate = *p
	case *A:
		predicate = "http://www.w3.org/1999/02/22-rdf-syntax-ns#type"
	default:
		panic(fmt.Errorf("unknown predicate type %T", v))
	}
	objects, triples, err := ctx.EvaluateObjectList(po.ObjectList)
	if err != nil {
		return "", nil, nil, err
	}
	return predicate, objects, triples, nil
}

func (ctx *Context) EvaluateStringLiteral(o *StringLiteral) (*nt.Literal, error) {
	v := o.Value
	if o.Multiline {
		v = strings.ReplaceAll(v, "\n", "\\n")
		v = strings.ReplaceAll(v, "\r", "\\r")
		v = strings.ReplaceAll(v, "\\\"", "\"")
	}
	v = strings.ReplaceAll(v, "\"", "\\\"")
	if o.LanguageTag != "" {
		return &nt.Literal{
			Value:    v,
			Language: o.LanguageTag,
		}, nil
	}
	return &nt.Literal{
		Value: v,
	}, nil
}

func (ctx *Context) EvaluateTriple(t *Triple) ([]nt.Triple, error) {
	var triples []nt.Triple
	var subject nt.Subject
	if t.Subject != nil {
		switch t := t.Subject.(type) {
		case *IRI:
			s, err := ctx.EvaluateIRI(t)
			if err != nil {
				return nil, err
			}
			subject = s
		case *BlankNode:
			if *t == "[]" {
				bn := ctx.bn()
				subject = &bn
			} else {
				bn := nt.BlankNode(*t)
				subject = &bn
			}
		case Collection:
			var objects []nt.Object
			for _, o := range t {
				os, ts, err := ctx.EvaluateObject(o)
				if err != nil {
					return nil, err
				}
				objects = append(objects, os...)
				triples = append(triples, ts...)
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
				triples = append(triples, nt.Triple{
					Subject:   &e,
					Predicate: "http://www.w3.org/1999/02/22-rdf-syntax-ns#first",
					Object:    o,
				})
				if i+1 != len(objects) {
					el = ctx.el()
					triples = append(triples, nt.Triple{
						Subject:   &e,
						Predicate: "http://www.w3.org/1999/02/22-rdf-syntax-ns#rest",
						Object:    &el,
					})
				} else {
					el = e
				}
			}
			triples = append(triples, nt.Triple{
				Subject:   &el,
				Predicate: "http://www.w3.org/1999/02/22-rdf-syntax-ns#rest",
				Object:    nt.IRIReference("http://www.w3.org/1999/02/22-rdf-syntax-ns#nil"),
			})
		default:
			panic(fmt.Errorf("unknown subject type %T", t))
		}

		for _, po := range t.PredicateObjectList {
			p, os, ts, err := ctx.EvaluatePredicateObject(po)
			if err != nil {
				return nil, err
			}
			triples = append(triples, ts...)
			for _, o := range os {
				triples = append(triples, nt.Triple{
					Subject:   subject,
					Predicate: p,
					Object:    o,
				})
			}
		}
	} else {
		bn := ctx.bn()
		for _, po := range t.BlankNodePropertyList {
			p, os, ts, err := ctx.EvaluatePredicateObject(po)
			if err != nil {
				return nil, err
			}
			triples = append(triples, ts...)
			for _, o := range os {
				triples = append(triples, nt.Triple{
					Subject:   &bn,
					Predicate: p,
					Object:    o,
				})
			}
		}
		for _, po := range t.PredicateObjectList {
			p, os, ts, err := ctx.EvaluatePredicateObject(po)
			if err != nil {
				return nil, err
			}
			triples = append(triples, ts...)
			for _, o := range os {
				triples = append(triples, nt.Triple{
					Subject:   &bn,
					Predicate: p,
					Object:    o,
				})
			}
		}
	}
	return triples, nil
}

func (ctx *Context) evaluateDocument(d Document) (nt.Document, error) {
	var triples []nt.Triple
	for _, t := range d {
		switch t := t.(type) {
		case *Base:
			ctx.Base = string(*t)
		case *Prefix:
			ctx.Prefixes[t.Name] = t.IRI
		case *Triple:
			ts, err := ctx.EvaluateTriple(t)
			if err != nil {
				return nil, err
			}
			triples = append(triples, ts...)
		default:
			panic(fmt.Errorf("unknown document type %T", t))
		}
	}
	return triples, nil
}
