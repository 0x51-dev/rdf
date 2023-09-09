package turtle

import (
	"fmt"
	"strings"
)

func (ctx *context) validateCollection(c Collection) bool {
	for _, o := range c {
		if !ctx.validateObject(o) {
			return false
		}
	}
	return true
}

func (ctx *context) validateDocument(d Document) bool {
	for _, t := range d {
		switch t := t.(type) {
		case *Base:
			ctx.base = string(*t)
		case *Prefix:
			ctx.prefixes[t.Name] = t.IRI
		case *Triple:
			if t.Subject != nil {
				if !ctx.validateSubject(t.Subject) {
					return false
				}
				return ctx.validatePredicateObjectList(t.PredicateObjectList)
			} else {
				if !ctx.validatePredicateObjectList((PredicateObjectList)(t.BlankNodePropertyList)) {
					return false
				}
				return ctx.validatePredicateObjectList(t.PredicateObjectList)
			}
		default:
			panic(fmt.Errorf("unknown document type %T", t))
		}
	}
	return true
}

func (ctx *context) validateIRI(i *IRI) bool {
	if !i.Prefixed {
		return true
	}
	p := strings.SplitAfterN(i.Value, ":", 2)
	if len(p) != 2 {
		return false
	}
	_, ok := ctx.prefixes[p[0]]
	return ok
}

func (ctx *context) validateObject(o Object) bool {
	switch o := o.(type) {
	case *IRI:
		return ctx.validateIRI(o)
	case *BlankNode:
		return true
	case BlankNodePropertyList:
		return ctx.validatePredicateObjectList((PredicateObjectList)(o))
	case *NumericLiteral, *BooleanLiteral, *StringLiteral:
		return true
	case Collection:
		return ctx.validateCollection(o)
	default:
		panic(fmt.Errorf("unknown object type %T", o))
	}
}

func (ctx *context) validatePredicateObject(po PredicateObject) bool {
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

func (ctx *context) validatePredicateObjectList(p PredicateObjectList) bool {
	for _, po := range p {
		if !ctx.validatePredicateObject(po) {
			return false
		}
	}
	return true
}

func (ctx *context) validateSubject(s Subject) bool {
	switch s := s.(type) {
	case *IRI:
		return ctx.validateIRI(s)
	case *BlankNode:
		return true
	case Collection:
		return ctx.validateCollection(s)
	default:
		panic(fmt.Errorf("unknown subject type %T", s))
	}
}

func (ctx *context) validateVerb(v Verb) bool {
	switch v := v.(type) {
	case *IRI:
		return ctx.validateIRI(v)
	case *A:
		return true
	default:
		panic(fmt.Errorf("unknown verb type %T", v))
	}
}
