package rdf

import "fmt"

type DataType string

// INFO on datatypes: https://www.w3.org/TR/xmlschema11-2/type-hierarchy-201104.longdesc.html
const (
	XSD                      = "http://www.w3.org/2001/XMLSchema#"
	XSDAnyType      DataType = XSD + "anyType"
	XSDAnyURI       DataType = XSD + "anyURI"
	XSDBase64Binary DataType = XSD + "base64Binary"
	XSDBoolean      DataType = XSD + "boolean"
	XSDDate         DataType = XSD + "date"
	XDSDecimal      DataType = XSD + "decimal"
	XSDDouble       DataType = XSD + "double"
	XSDDuration     DataType = XSD + "duration"
	XSDFloat        DataType = XSD + "float"
	XSDString       DataType = XSD + "string"
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
	default:
		return false
	}
}
