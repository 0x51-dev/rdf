package rdf

import (
	"fmt"
	"math/big"
	"strings"
)

type DataType string

// INFO on datatypes: https://www.w3.org/TR/xmlschema11-2/type-hierarchy-201104.longdesc.html
const (
	XSD                 = "http://www.w3.org/2001/XMLSchema#"
	XSDAnyType DataType = XSD + "anyType"
	XSDBoolean DataType = XSD + "boolean"
	XSDInteger DataType = XSD + "integer"
	XSDDecimal DataType = XSD + "decimal"
	XSDDouble  DataType = XSD + "double"
	XSDString  DataType = XSD + "string"
	// TODO: add "other" build-in atomic types, e.g. dateTimeStamp etc.

	XSDNS       = "http://www.w3.org/1999/02/22-rdf-syntax-ns#"
	XSDNSString = XSDNS + "langString"
)

// NativeType returns the native type of the given value.
// Returns true if the value was converted, otherwise false.
func (d DataType) NativeType(value string) (any, bool, error) {
	switch d {
	case XSDAnyType, XSDString:
		return value, true, nil
	case XSDBoolean:
		switch value {
		case "true", "false":
			return value[0] == 't', true, nil
		}
		return nil, false, fmt.Errorf("invalid boolean value: %q", value)
	case XSDInteger:
		if i, ok := new(big.Int).SetString(value, 10); ok {
			return i, true, nil
		}
		return nil, false, fmt.Errorf("invalid integer value: %q", value)
	case XSDDecimal:
		if strings.ContainsAny(value, "eE") {
			return nil, false, fmt.Errorf("invalid decimal value: %q", value)
		}
		if f, ok := new(big.Float).SetString(value); ok {
			return f, true, nil
		}
		return nil, false, fmt.Errorf("invalid decimal value: %q", value)
	case XSDDouble:
		if value == "INF" || value == "-INF" || value == "NaN" {
			return value, true, nil
		}
		if f, ok := new(big.Float).SetString(value); ok {
			return f, true, nil
		}
		return nil, false, fmt.Errorf("invalid decimal value: %q", value)
	default:
		return value, false, nil
	}
}

// Validate returns true if the given value is valid for the datatype.
func (d DataType) Validate(v any, acceptString bool) bool {
	if _, ok := v.(string); acceptString && ok {
		// TODO: validate string against the datatype.
		return true
	}
	switch d {
	case XSDAnyType:
		return true
	case XSDString:
		_, ok := v.(string)
		return ok
	case XSDBoolean:
		_, ok := v.(bool)
		return ok
	case XSDInteger:
		_, ok := v.(*big.Int)
		return ok
	case XSDDecimal:
		_, ok := v.(*big.Float)
		return ok
	case XSDDouble:
		if _, ok := v.(*big.Float); ok {
			return true
		}
		_, ok := v.(string)
		return ok && (v == "INF" || v == "-INF" || v == "NaN")
	default:
		return false
	}
}
