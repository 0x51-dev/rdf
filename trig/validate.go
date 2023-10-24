package trig

import (
	"fmt"
	ttl "github.com/0x51-dev/rdf/turtle"
)

func (ctx *Context) EvaluateWrappedGraph(wg WrappedGraph) bool {
	for _, t := range wg {
		if !ctx.ValidateTriple(&t) {
			return false
		}
	}
	return true
}

func (ctx *Context) validateDocument(d Document) bool {
	for _, t := range d {
		switch t := t.(type) {
		case *Base:
			ctx.Base = string(*t)
		case *Prefix:
			ctx.Prefixes[t.Name] = t.IRI
		case *TriplesOrGraph:
			switch t := t.LabelOrSubject.(type) {
			case *IRI:
				if !ctx.ValidateIRI((*ttl.IRI)(t)) {
					return false
				}
			case *BlankNode:
			default:
				panic(fmt.Errorf("unknown document type %T", t))
			}
			if len(t.WrappedGraph) != 0 {
				if !ctx.EvaluateWrappedGraph(t.WrappedGraph) {
					return false
				}
			} else {
				if !ctx.ValidatePredicateObjectList(t.PredicateObjectList) {
					return false
				}
			}
		case WrappedGraph:
			if !ctx.EvaluateWrappedGraph(t) {
				return false
			}
		case Triple2:
			if len(t.BlankNodePropertyList) != 0 {
				if !ctx.ValidatePredicateObjectList((ttl.PredicateObjectList)(t.BlankNodePropertyList)) {
					return false
				}
				if len(t.PredicateObjectList) != 0 {
					if !ctx.ValidatePredicateObjectList(t.PredicateObjectList) {
						return false
					}
				}
			} else {
				if !ctx.ValidateCollection(t.Collection) {
					return false
				}
				if !ctx.ValidatePredicateObjectList(t.PredicateObjectList) {
					return false
				}
			}
		default:
			panic(fmt.Errorf("unknown document type %T", t))
		}
	}
	return true
}
