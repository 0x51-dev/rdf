package grammar

import (
	nt "github.com/0x51-dev/rdf/ntriples/grammar"
	"github.com/0x51-dev/upeg/parser"
	"github.com/0x51-dev/upeg/parser/op"
	"unicode"
)

var (
	Document = op.Capture{
		Name: "Document",
		Value: op.ZeroOrMore{Value: op.And{
			nt.OWhitespace,
			op.Or{
				op.Or{Directive, Triples},
				op.And{op.Optional{Value: nt.Comment}, op.EndOfLine{}},
			},
		}},
	}
	Directive = op.Capture{
		Name:  "Directive",
		Value: op.Or{PrefixID, Base, SparqlPrefix, SparqlBase},
	}
	PrefixID = op.Capture{
		Name: "Prefix",
		Value: op.And{
			"@prefix", nt.Whitespace,
			op.Capture{Name: "Prefix", Value: PNAME_NS}, nt.Whitespace,
			nt.IRIReference, nt.OWhitespace,
			'.',
		},
	}
	Base = op.Capture{
		Name: "Base",
		Value: op.And{
			"@base", nt.Whitespace,
			nt.IRIReference, nt.OWhitespace,
			'.',
		},
	}
	SparqlPrefix = op.Capture{
		Name: "Prefix",
		Value: op.And{
			CaseInsensitiveString("PREFIX"), nt.Whitespace,
			op.Capture{Name: "Prefix", Value: PNAME_NS}, nt.Whitespace,
			nt.IRIReference,
		},
	}
	SparqlBase = op.Capture{
		Name: "Base",
		Value: op.And{
			CaseInsensitiveString("BASE"), nt.Whitespace,
			nt.IRIReference,
		},
	}
	Triples = op.Capture{
		Name: "Triples",
		Value: op.And{
			op.Or{
				op.Capture{
					Name: "TripleSubject",
					Value: op.And{
						Subject,
						WSPLNC,
						PredicateObjectList,
					},
				},
				op.Capture{
					Name: "TripleBlankNodePropertyList",
					Value: op.And{
						BlankNodePropertyList,
						WSPLNC,
						op.Optional{Value: PredicateObjectList},
					},
				},
			},
			WSPLNC,
			'.',
		},
	}
	PredicateObject = op.Capture{
		Name: "PredicateObject",
		Value: op.And{
			Verb,
			WSPLNC,
			ObjectList,
		},
	}
	PredicateObjectList = op.Capture{
		Name: "PredicateObjectList",
		Value: op.And{
			PredicateObject,
			op.ZeroOrMore{Value: op.And{
				WSPLNC,
				';',
				WSPLNC,
				op.Optional{Value: PredicateObject},
			}},
		},
	}
	ObjectList = op.Capture{
		Name: "ObjectList",
		Value: op.And{
			Object, op.ZeroOrMore{Value: op.And{
				nt.OWhitespace, ',',
				WSPLNC, Object,
			}},
		},
	}
	Verb = op.Capture{
		Name:  "Verb",
		Value: op.Or{IRI, op.Capture{Name: "a", Value: 'a'}},
	}
	Subject = op.Capture{
		Name:  "Subject",
		Value: op.Or{IRI, BlankNode, Collection},
	}
	Object = op.Capture{
		Name: "Object",
		Value: op.Or{
			Literal,
			IRI,
			BlankNode,
			op.Reference{Name: "Collection"},
			op.Reference{Name: "BlankNodePropertyList"},
		},
	}
	Literal = op.Capture{
		Name:  "Literal",
		Value: op.Or{RDFLiteral, NumericLiteral, BooleanLiteral},
	}
	BlankNodePropertyList = op.Capture{
		Name: "BlankNodePropertyList",
		Value: op.And{
			'[',
			WSPLNC,
			PredicateObjectList,
			WSPLNC,
			']',
		},
	}
	Collection = op.Capture{
		Name: "Collection",
		Value: op.And{
			'(',
			WSPLNC,
			op.ZeroOrMore{Value: op.And{
				Object,
				WSPLNC,
			}},
			WSPLNC,
			nt.OWhitespace, ')'},
	}
	NumericLiteral = op.Capture{
		Name:  "NumericLiteral",
		Value: op.Or{Double, Decimal, Integer},
	}
	RDFLiteral = op.Capture{
		Name:  "RDFLiteral",
		Value: op.And{String, op.Optional{Value: op.Or{nt.LanguageTag, op.And{"^^", IRI}}}},
	}
	BooleanLiteral = op.Capture{
		Name:  "BooleanLiteral",
		Value: op.Or{"true", "false"},
	}
	String = op.Or{
		STRING_LITERAL_LONG_SINGLE_QUOTE,
		STRING_LITERAL_LONG_QUOTE,
		STRING_LITERAL_QUOTE,
		STRING_LITERAL_SINGLE_QUOTE,
	}
	IRI = op.Capture{
		Name: "IRI",
		Value: op.Or{
			nt.IRIReference,
			PrefixedName,
		},
	}
	PrefixedName = op.Capture{
		Name: "PrefixedName",
		Value: op.And{
			op.Not{Value: "_:"},
			op.And{PNAME_NS, op.Optional{Value: PN_LOCAL}},
		},
	}
	BlankNode = op.Capture{
		Name:  "BlankNode",
		Value: op.Or{nt.BlankNodeLabel, ANON},
	}
	PNAME_NS = op.OneOrMore{Value: op.And{
		op.Optional{Value: op.And{
			nt.PN_CHARS_BASE,
			op.ZeroOrMore{Value: op.And{
				op.Not{Value: ':'},
				op.Or{
					nt.PN_CHARS,
					op.And{
						op.OneOrMore{Value: '.'},
						op.Peek{Value: op.And{
							op.Not{Value: ':'},
							nt.PN_CHARS,
						}},
					},
				},
			}},
		}},
		':',
	}}
	Integer = op.Capture{
		Name: "Integer",
		Value: op.And{
			op.Optional{Value: op.Or{"+", "-"}},
			op.OneOrMore{Value: op.RuneRange{Min: '0', Max: '9'}},
		},
	}
	Decimal = op.Capture{
		Name: "Decimal",
		Value: op.And{
			op.Optional{Value: op.Or{"+", "-"}},
			op.ZeroOrMore{Value: op.RuneRange{Min: '0', Max: '9'}},
			'.',
			op.OneOrMore{Value: op.RuneRange{Min: '0', Max: '9'}},
		},
	}
	Double = op.Capture{
		Name: "Double",
		Value: op.And{
			op.Optional{Value: op.Or{"+", "-"}},
			op.Or{
				op.And{
					op.OneOrMore{Value: op.RuneRange{Min: '0', Max: '9'}},
					'.',
					op.ZeroOrMore{Value: op.RuneRange{Min: '0', Max: '9'}},
					Exponent,
				},
				op.And{
					op.Optional{Value: '.'},
					op.OneOrMore{Value: op.RuneRange{Min: '0', Max: '9'}},
					Exponent,
				},
			},
		},
	}
	Exponent = op.And{
		op.Or{'e', 'E'},
		op.Optional{Value: op.Or{'+', '-'}},
		op.OneOrMore{Value: op.RuneRange{Min: '0', Max: '9'}},
	}
	STRING_LITERAL_QUOTE = op.And{
		'"',
		op.Capture{
			Name: "StringLiteralQ",
			Value: op.ZeroOrMore{Value: op.Or{
				op.AnyBut{Value: op.Or{rune(0x22), rune(0x5C), rune(0x0A), rune(0x0D)}},
				ECHAR,
				UCHAR,
			}},
		},
		'"',
	}
	STRING_LITERAL_SINGLE_QUOTE = op.And{
		'\'',
		op.Capture{
			Name: "StringLiteralSQ",
			Value: op.ZeroOrMore{Value: op.Or{
				op.AnyBut{Value: op.Or{rune(0x27), rune(0x5C), rune(0x0A), rune(0x0D)}},
				ECHAR,
				UCHAR,
			}},
		},
		'\'',
	}
	STRING_LITERAL_LONG_QUOTE = op.And{
		"\"\"\"",
		op.Capture{
			Name: "StringLiteralLQ",
			Value: op.ZeroOrMore{Value: op.And{
				op.Optional{Value: op.Or{"\"\"", '"'}},
				op.Or{
					op.AnyBut{Value: op.Or{'"', '\\'}},
					ECHAR,
					UCHAR,
				},
			}},
		},
		"\"\"\"",
	}
	STRING_LITERAL_LONG_SINGLE_QUOTE = op.And{
		"'''",
		op.Capture{
			Name: "StringLiteralLSQ",
			Value: op.ZeroOrMore{Value: op.And{
				op.Optional{Value: op.Or{"''", '\''}},
				op.Or{
					op.AnyBut{Value: op.Or{'\'', '\\'}},
					ECHAR,
					UCHAR,
				},
			}},
		},
		"'''",
	}
	UCHAR = op.Or{
		op.And{"\\u", nt.Hexadecimal, nt.Hexadecimal, nt.Hexadecimal, nt.Hexadecimal},
		op.And{"\\U", nt.Hexadecimal, nt.Hexadecimal, nt.Hexadecimal, nt.Hexadecimal, nt.Hexadecimal, nt.Hexadecimal, nt.Hexadecimal, nt.Hexadecimal},
	}
	ECHAR = op.And{'\\', op.Or{'t', 'b', 'n', 'r', 'f', '"', '\'', '\\'}}
	WS    = op.Or{rune(0x20), rune(0x09), rune(0x0D), rune(0x0A)}
	ANON  = op.Capture{
		Name:  "Anon",
		Value: op.And{'[', op.ZeroOrMore{Value: WS}, ']'},
	}
	PN_LOCAL = op.And{
		op.Or{nt.PN_CHARS_U, ':', op.RuneRange{Min: '0', Max: '9'}, PLX},
		op.Optional{Value: op.And{
			op.ZeroOrMore{Value: op.Or{
				op.And{
					op.OneOrMore{Value: '.'},
					op.Peek{Value: op.Or{nt.PN_CHARS, ':', PLX}},
				},
				op.Or{nt.PN_CHARS, ':', PLX},
			}},
		}},
	}
	PLX          = op.Or{PERCENT, PN_LOCAL_ESC}
	PERCENT      = op.And{'%', nt.Hexadecimal, nt.Hexadecimal}
	PN_LOCAL_ESC = op.And{
		'\\',
		op.Or{'_', '~', '.', '-', '!', '$', '&', '\'', '(', ')', '*', '+', ',', ';', '=', '/', '?', '#', '@', '%'},
	}
	WSPLNC = op.ZeroOrMore{Value: op.Or{
		nt.Whitespace,
		op.EndOfLine{},
		nt.Comment,
	}}
)

func NewParser(input []rune) (*parser.Parser, error) {
	p, err := parser.New(input)
	if err != nil {
		return nil, err
	}
	p.Rules["BlankNodePropertyList"] = BlankNodePropertyList
	p.Rules["Collection"] = Collection
	return p, nil
}

// CaseInsensitiveString checks if a string matches a given string, ignoring case.
type CaseInsensitiveString string

func (s CaseInsensitiveString) Match(start parser.Cursor, p *parser.Parser) (parser.Cursor, error) {
	end := start
	for _, r := range s {
		if unicode.ToLower(end.Character()) != unicode.ToLower(r) {
			p.Reader.Jump(start)
			return start, p.NewNoMatchError(r, start, end)
		}
		end = p.Reader.Next().Cursor()
	}
	return end, nil
}

func (s CaseInsensitiveString) String() string {
	return string(s)
}
