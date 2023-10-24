package grammar

import (
	"github.com/0x51-dev/upeg/parser/op"
)

var (
	Document = op.Capture{
		Name: "Document",
		Value: op.ZeroOrMore{Value: op.And{
			op.Optional{Value: op.Or{Triple, Comment}},
			OWhitespace, op.EndOfLine{},
		}},
	}
	Triple = op.Capture{
		Name: "Triple",
		Value: op.And{
			OWhitespace,
			Subject, OWhitespace,
			Predicate, OWhitespace,
			Object, OWhitespace,
			'.', op.Optional{Value: Comment},
		},
	}
	Subject = op.Capture{
		Name:  "Subject",
		Value: op.Or{IRIReference, BlankNodeLabel},
	}
	Predicate = op.Capture{
		Name:  "Predicate",
		Value: IRIReference,
	}
	Object = op.Capture{
		Name:  "Object",
		Value: op.Or{IRIReference, BlankNodeLabel, Literal},
	}
	Literal = op.Capture{
		Name: "Literal",
		Value: op.And{StringLiteral, op.Optional{Value: op.Or{
			op.And{"^^", IRIReference},
			LanguageTag,
		}}},
	}

	OWhitespace = op.ZeroOrMore{Value: op.Or{rune(0x20), rune(0x09)}}
	Whitespace  = op.OneOrMore{Value: op.Or{rune(0x20), rune(0x09)}}
	Comment     = op.And{
		OWhitespace,
		'#',
		op.ZeroOrMore{Value: op.AnyBut{Value: op.Or{rune(0x0A), rune(0x0D)}}},
	}
	LanguageTag = op.And{
		'@',
		op.Capture{
			Name: "LanguageTag",
			Value: op.And{
				op.OneOrMore{Value: op.Or{op.RuneRange{Min: 'a', Max: 'z'}, op.RuneRange{Min: 'A', Max: 'Z'}}},
				op.ZeroOrMore{Value: op.And{
					'-',
					op.OneOrMore{Value: op.Or{op.RuneRange{Min: 'a', Max: 'z'}, op.RuneRange{Min: 'A', Max: 'Z'}, op.RuneRange{Min: '0', Max: '9'}}}, // [a-zA-Z0-9]+
				}},
			},
		},
	}
	IRIReference = op.And{
		'<',
		op.Capture{
			Name: "IRIReference",
			Value: op.And{
				op.ZeroOrMore{Value: op.Or{
					op.AnyBut{Value: op.Or{
						op.RuneRange{Min: 0x00, Max: 0x20},
						'<', '>', '"', '{', '}', '|', '^', '`', '\\',
					}},
					UnicodeCharacter,
				}},
			},
		},
		'>',
	}
	StringLiteral = op.And{
		'"',
		op.Capture{
			Name: "StringLiteral",
			Value: op.And{
				op.ZeroOrMore{Value: op.Or{
					op.AnyBut{Value: op.Or{
						rune(0x22), rune(0x5C), rune(0x0A), rune(0x0D),
					}},
					EscapedCharacter, UnicodeCharacter,
				}},
			},
		},
		'"',
	}
	BlankNodeLabel = op.And{
		"_:",
		op.Capture{
			Name: "BlankNodeLabel",
			Value: op.And{
				op.Or{op.RuneRange{Min: '0', Max: '9'}, PN_CHARS_U},
				op.Optional{Value: op.And{
					op.ZeroOrMore{Value: op.And{
						op.Or{PN_CHARS, '.'},
						// To make sure that the last character is not a dot.
						op.Peek{Value: op.Or{
							PN_CHARS,              // Either end or another character.
							op.And{'.', PN_CHARS}, // '.' always followed by a '.' or character.
							op.And{'.', '.'},
						}},
					}},
					PN_CHARS,
				}},
			},
		},
	}
	UnicodeCharacter = op.Or{
		op.And{"\\u", op.Repeat{Min: 4, Max: 4, Value: Hexadecimal}},
		op.And{"\\U", op.Repeat{Min: 8, Max: 8, Value: Hexadecimal}},
	}
	EscapedCharacter = op.And{
		'\\',
		op.Or{
			't', 'b', 'n', 'r', 'f', '"', '\'', '\\',
		},
	}
	PN_CHARS_BASE = op.Or{
		op.RuneRange{Min: 'A', Max: 'Z'},
		op.RuneRange{Min: 'a', Max: 'z'},
		op.RuneRange{Min: 0x00C0, Max: 0x00D6},
		op.RuneRange{Min: 0x00D8, Max: 0x00F6},
		op.RuneRange{Min: 0x00F8, Max: 0x02FF},
		op.RuneRange{Min: 0x0370, Max: 0x037D},
		op.RuneRange{Min: 0x037F, Max: 0x1FFF},
		op.RuneRange{Min: 0x200C, Max: 0x200D},
		op.RuneRange{Min: 0x2070, Max: 0x218F},
		op.RuneRange{Min: 0x2C00, Max: 0x2FEF},
		op.RuneRange{Min: 0x3001, Max: 0xD7FF},
		op.RuneRange{Min: 0xF900, Max: 0xFDCF},
		op.RuneRange{Min: 0xFDF0, Max: 0xFFFD},
		op.RuneRange{Min: 0x10000, Max: 0xEFFFF},
	}
	PN_CHARS_U = op.Or{'_', PN_CHARS_BASE}
	PN_CHARS   = op.Or{
		'-',
		op.RuneRange{Min: '0', Max: '9'},
		PN_CHARS_U,
		rune(0x00B7),
		op.RuneRange{Min: 0x0300, Max: 0x036F},
		op.RuneRange{Min: 0x203F, Max: 0x2040},
	}
	Hexadecimal = op.Or{
		op.RuneRange{Min: '0', Max: '9'},
		op.RuneRange{Min: 'A', Max: 'F'},
		op.RuneRange{Min: 'a', Max: 'f'},
	}
)
