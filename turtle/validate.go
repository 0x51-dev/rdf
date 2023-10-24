package turtle

import (
	"fmt"
	"strings"
)

func (ctx *Context) ValidateCollection(c Collection) bool {
	for _, o := range c {
		if !ctx.validateObject(o) {
			return false
		}
	}
	return true
}

func (ctx *Context) ValidateIRI(i *IRI) bool {
	if !i.Prefixed {
		return true
	}
	p := strings.SplitAfterN(i.Value, ":", 2)
	if len(p) != 2 {
		return false
	}
	_, ok := ctx.Prefixes[p[0]]
	return ok
}

func (ctx *Context) ValidatePredicateObjectList(p PredicateObjectList) bool {
	for _, po := range p {
		if !ctx.validatePredicateObject(po) {
			return false
		}
	}
	return true
}

func (ctx *Context) ValidateTriple(t *Triple) bool {
	if t.Subject != nil {
		if !ctx.validateSubject(t.Subject) {
			return false
		}
		return ctx.ValidatePredicateObjectList(t.PredicateObjectList)
	}

	if !ctx.ValidatePredicateObjectList((PredicateObjectList)(t.BlankNodePropertyList)) {
		return false
	}
	return ctx.ValidatePredicateObjectList(t.PredicateObjectList)
}

func (ctx *Context) validateDocument(d Document) bool {
	for _, t := range d {
		switch t := t.(type) {
		case *Base:
			ctx.Base = string(*t)
		case *Prefix:
			ctx.Prefixes[t.Name] = t.IRI
		case *Triple:
			if !ctx.ValidateTriple(t) {
				return false
			}
		default:
			panic(fmt.Errorf("unknown document type %T", t))
		}
	}
	return true
}

func (ctx *Context) validateObject(o Object) bool {
	switch o := o.(type) {
	case *IRI:
		return ctx.ValidateIRI(o)
	case *BlankNode:
		return true
	case BlankNodePropertyList:
		return ctx.ValidatePredicateObjectList((PredicateObjectList)(o))
	case *NumericLiteral, *BooleanLiteral, *StringLiteral:
		return true
	case Collection:
		return ctx.ValidateCollection(o)
	default:
		panic(fmt.Errorf("unknown object type %T", o))
	}
}

func (ctx *Context) validatePredicateObject(po PredicateObject) bool {
	if !ctx.validateVerb(po.Verb) {
		return false
	}
	for _, o := range po.ObjectList {
		if !ctx.validateObject(o) {
			return false
		}
	}
	return true
}

func (ctx *Context) validateSubject(s Subject) bool {
	switch s := s.(type) {
	case *IRI:
		return ctx.ValidateIRI(s)
	case *BlankNode:
		return true
	case Collection:
		return ctx.ValidateCollection(s)
	default:
		panic(fmt.Errorf("unknown subject type %T", s))
	}
}

func (ctx *Context) validateVerb(v Verb) bool {
	switch v := v.(type) {
	case *IRI:
		return ctx.ValidateIRI(v)
	case *A:
		return true
	default:
		panic(fmt.Errorf("unknown verb type %T", v))
	}
}
