package rdf

import "testing"

// https://www.w3.org/TR/xmlschema-2/#decimal
func TestDataType_NativeType_decimal(t *testing.T) {
	for _, test := range []string{
		"-1", "12678967", "+100000", "210",
		"-1.23", "12678967.543233", "+100000.00", "210",
	} {
		if _, ok, err := XSDDecimal.NativeType(test); !ok || err != nil {
			t.Errorf("XSDDecimal.NativeType(%q) = %v, %v; want %v, nil", test, ok, err, true)
		}
	}
	for _, test := range []string{
		"-1E4", "1267.43233E12", "12.78e-2", "INF",
	} {
		if _, ok, err := XSDDecimal.NativeType(test); ok || err == nil {
			t.Errorf("XSDDecimal.NativeType(%q) = %v, %v; want %v, nil", test, ok, err, false)
		}
	}
}

// https://www.w3.org/TR/xmlschema-2/#double
func TestDataType_NativeType_double(t *testing.T) {
	for _, test := range []string{
		"-1", "12678967", "+100000", "210",
		"-1.23", "12678967.543233", "+100000.00", "210",
		"-1E4", "1267.43233E12", "12.78e-2", "12", "-0", "0", "INF",
	} {
		if _, ok, err := XSDDouble.NativeType(test); !ok || err != nil {
			t.Errorf("XSDDouble.NativeType(%q) = %v, %v; want %v, nil", test, ok, err, true)
		}
	}
}

// https://www.w3.org/TR/xmlschema-2/#integer
func TestDataType_NativeType_integer(t *testing.T) {
	for _, test := range []string{
		"-1", "12678967", "+100000", "210",
	} {
		if _, ok, err := XSDInteger.NativeType(test); !ok || err != nil {
			t.Errorf("XSDInteger.NativeType(%q) = %v, %v; want %v, nil", test, ok, err, true)
		}
	}
	for _, test := range []string{
		"-1.23", "12678967.543233", "+100000.00",
		"-1E4", "1267.43233E12", "12.78e-2", "INF",
	} {
		if _, ok, err := XSDInteger.NativeType(test); ok || err == nil {
			t.Errorf("XSDInteger.NativeType(%q) = %v, %v; want %v, nil", test, ok, err, false)
		}
	}
}
