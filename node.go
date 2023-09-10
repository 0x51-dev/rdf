package rdf

import (
	"fmt"
	"strings"
)

// BlankNode identifiers are local identifiers.
type BlankNode struct {
	Attribute string
}

func (b *BlankNode) Equal(other Node) bool {
	if other, ok := other.(*BlankNode); ok {
		return b.Attribute == other.Attribute
	}
	return false
}

func (b *BlankNode) GetValue() string {
	return b.Attribute
}

func (b *BlankNode) toObject(_ bool) (map[string]any, error) {
	return map[string]any{
		"@id": b.Attribute,
	}, nil
}

type IRIReference struct {
	Value string
}

func (r *IRIReference) Equal(other Node) bool {
	if other, ok := other.(*IRIReference); ok {
		return r.Value == other.Value
	}
	return false
}

func (r *IRIReference) GetValue() string {
	return r.Value
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

func (l *Literal) Equal(other Node) bool {
	if other, ok := other.(*Literal); ok {
		return l.Value == other.Value && l.Datatype == other.Datatype && l.Language == other.Language
	}
	return false
}

func (l *Literal) GetValue() string {
	return l.Value
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
	// GetValue returns the value of the node.
	GetValue() string
	// Equal returns true if the two nodes are equal.
	Equal(other Node) bool

	toObject(nativeTypes bool) (map[string]any, error)
}

func fromObject(obj any, acceptString bool) (Node, error) {
	var id string
	switch obj := obj.(type) {
	case map[string]any:
		if v, ok := obj["@value"]; ok {
			var literal = &Literal{
				Value:    fmt.Sprintf("%v", v),
				Datatype: XSDString,
			}
			if datatype, ok := obj["@type"]; ok {
				s, ok := datatype.(string)
				if !ok {
					return nil, fmt.Errorf("invalid @datatype attribute: %T", datatype)
				}
				literal.Datatype = DataType(s)
			}
			if !literal.Datatype.Validate(v, acceptString) {
				return nil, fmt.Errorf("invalid @value for datatype %s: %v", literal.Datatype, v)
			}

			if language, ok := obj["@language"]; ok {
				s, ok := language.(string)
				if !ok {
					return nil, fmt.Errorf("invalid @language attribute: %T", language)
				}
				literal.Language = s
				if literal.Datatype != "" && literal.Datatype != XSDNSString {
					return nil, fmt.Errorf("invalid datatype for language literal: %s", literal.Datatype)
				}
				if literal.Datatype == "" {
					literal.Datatype = XSDNSString
				}
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
