package rdf

import (
	"fmt"
	"strings"
)

// BlankNode identifiers are local identifiers.
type BlankNode struct {
	Attribute string
}

func (b *BlankNode) toObject(_ bool) (map[string]any, error) {
	return map[string]any{
		"@id": b.Attribute,
	}, nil
}

type IRIReference struct {
	Value string
}

func (r *IRIReference) toObject(_ bool) (map[string]any, error) {
	return map[string]any{
		"@id": r.Value,
	}, nil
}

// Literal is used to represent values such as strings, numbers, and dates.
type Literal struct {
	// Value is a Unicode string.
	Value string
	// Datatype is an IRI identifying a datatype that determines how the lexical form maps to a literal value.
	Datatype DataType
	// Language can only be present if Datatype is XSDNSString.
	Language string
}

func (l *Literal) toObject(nativeTypes bool) (map[string]any, error) {
	if l.Language != "" {
		if l.Datatype != "" && l.Datatype != XSDNSString {
			return nil, fmt.Errorf("invalid datatype for language literal: %s", l.Datatype)
		}
		return map[string]any{
			"@value":    l.Value,
			"@language": l.Language,
		}, nil
	}
	if l.Datatype != XSDString {
		if nativeTypes {
			v, _, err := l.Datatype.NativeType(l.Value)
			if err != nil {
				return nil, err
			}
			return map[string]any{
				"@value":    v,
				"@datatype": l.Datatype,
			}, nil
		}
		return map[string]any{
			"@value":    l.Value,
			"@datatype": l.Datatype,
		}, nil
	}
	return map[string]any{
		"@value": l.Value,
	}, nil
}

type Node interface {
	toObject(nativeTypes bool) (map[string]any, error)
}

func FromObject(obj any, acceptString bool) (Node, error) {
	var id string
	switch obj := obj.(type) {
	case map[string]any:
		if v, ok := obj["@value"]; ok {
			var literal = &Literal{
				Value:    fmt.Sprintf("%v", v),
				Datatype: XSDString,
			}
			if datatype, ok := obj["@datatype"]; ok {
				s, ok := datatype.(string)
				if !ok {
					return nil, fmt.Errorf("invalid @datatype attribute: %T", datatype)
				}
				literal.Datatype = DataType(s)
			}
			if !literal.Datatype.Validate(v, acceptString) {
				return nil, fmt.Errorf("invalid @value for datatype %s: %v", literal.Datatype, v)
			}
			return literal, nil
		}

		if v, ok := obj["@id"]; ok {
			id, ok = v.(string)
			if !ok {
				return nil, fmt.Errorf("invalid @id attribute: %T", v)
			}
		} else {
			return nil, fmt.Errorf("missing @id attribute")
		}
	case string:
		id = obj
	default:
		return nil, fmt.Errorf("unknown object type: %T", obj)
	}
	if strings.HasPrefix(id, "_:") {
		return &BlankNode{Attribute: id}, nil
	}
	return &IRIReference{Value: id}, nil
}
