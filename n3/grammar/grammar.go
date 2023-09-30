package grammar

import (
	nt "github.com/0x51-dev/rdf/ntriples/grammar"
	ttl "github.com/0x51-dev/rdf/turtle/grammar"
	"github.com/0x51-dev/upeg/parser"
	"github.com/0x51-dev/upeg/parser/op"
)

var (
	Document = op.Capture{
		Name: "Document",
		Value: op.ZeroOrMore{Value: op.And{
			nt.OWhitespace,
			op.Or{
				op.And{Statement, ttl.WSPLNC, '.'},
				SparqlDirective,
				op.And{op.Optional{Value: nt.Comment}, op.EndOfLine{}},
			},
		}},
	}
	Statement = op.Capture{
		Name:  "Statement",
		Value: op.Or{Directive, Triples},
	}
	Directive = op.Or{
		PrefixID,
		Base,
	}
	PrefixID = op.Capture{
		Name: "Prefix",
		Value: op.And{
			"@prefix", nt.Whitespace,
			op.Capture{Name: "Prefix", Value: ttl.PNAME_NS}, nt.OWhitespace,
			nt.IRIReference,
		},
	}
	Base = op.Capture{
		Name:  "Base",
		Value: op.And{"@base", nt.Whitespace, nt.IRIReference},
	}
	SparqlDirective = op.Or{
		ttl.SparqlPrefix,
		ttl.SparqlBase,
	}
	Triples = op.Capture{
		Name: "Triples",
		Value: op.And{
			Subject,
			ttl.WSPLNC,
			op.Optional{Value: PredicateObjectList},
		},
	}
	PredicateObject = op.Capture{
		Name: "PredicateObject",
		Value: op.And{
			Verb,
			ttl.WSPLNC,
			ObjectList,
		},
	}
	PredicateObjectList = op.Capture{
		Name: "PredicateObjectList",
		Value: op.And{
			PredicateObject,
			op.ZeroOrMore{Value: op.And{
				ttl.WSPLNC,
				';',
				ttl.WSPLNC,
				op.Optional{Value: PredicateObject},
			}},
		},
	}
	ObjectList = op.Capture{
		Name: "ObjectList",
		Value: op.And{
			Object,
			op.ZeroOrMore{Value: op.And{
				nt.OWhitespace, ',',
				ttl.WSPLNC, Object,
			}},
		},
	}
	Verb = op.Capture{
		Name: "Verb",
		Value: op.Or{
			Predicate,
			op.Capture{Name: "a", Value: 'a'},
			op.Capture{Name: "has", Value: op.And{"has", Path}},
			op.Capture{Name: "isOf", Value: op.And{"is", Path, "of"}},
			op.Capture{Name: "eql", Value: "<="},
			op.Capture{Name: "eqg", Value: "=>"},
			op.Capture{Name: "eq", Value: '='},
		},
	}
	Subject = op.Capture{
		Name:  "Subject",
		Value: Path,
	}
	Predicate = op.Or{
		Path,
		op.Capture{Name: "inverse", Value: op.And{"<-", Path}},
	}
	Object = op.Capture{
		Name:  "Object",
		Value: Path,
	}
	Path = op.Capture{
		Name: "Path",
		Value: op.And{
			op.Reference{Name: "PathItem"},
			op.Optional{Value: op.Or{
				op.Capture{Name: "forward", Value: op.And{'!', op.Reference{Name: "Path"}}},
				op.Capture{Name: "reverse", Value: op.And{'^', op.Reference{Name: "Path"}}},
			}},
		},
	}
	PathItem = op.Capture{
		Name: "PathItem",
		Value: op.Or{
			ttl.IRI,
			ttl.BlankNode,
			QuickVar,
			Collection,
			BlankNodePropertyList,
			IRIPropertyList,
			ttl.Literal,
			Formula,
		},
	}
	QuickVar = op.Capture{
		Name: "QuickVar",
		Value: op.And{
			'?', op.Capture{Name: "Var", Value: ttl.PN_LOCAL},
		},
	}
	BlankNodePropertyList = op.Capture{
		Name: "BlankNodePropertyList",
		Value: op.And{
			'[', ttl.WSPLNC, PredicateObjectList, ttl.WSPLNC, ']',
		},
	}
	IRIPropertyList = op.And{
		'[', nt.OWhitespace, "id", nt.Whitespace, ttl.IRI, ttl.WSPLNC, PredicateObjectList, ttl.WSPLNC, ']',
	}
	Collection = op.Capture{
		Name: "Collection",
		Value: op.And{
			'(', ttl.WSPLNC, op.ZeroOrMore{Value: op.And{Object, ttl.WSPLNC}}, ')',
		},
	}
	Formula = op.Capture{
		Name: "Formula",
		Value: op.And{
			'{', ttl.WSPLNC, op.Optional{Value: op.And{FormulaContent, ttl.WSPLNC}}, '}',
		},
	}
	FormulaContent = op.Capture{
		Name: "FormulaContent",
		Value: op.Or{
			op.And{
				Statement,
				op.Optional{Value: op.And{
					nt.Whitespace, '.', ttl.WSPLNC,
					op.Optional{Value: op.Reference{Name: "FormulaContent"}}},
				},
			},
			op.And{
				SparqlDirective,
				op.Optional{Value: op.Reference{Name: "FormulaContent"}},
			},
		},
	}
)

func NewParser(input []rune) (*parser.Parser, error) {
	p, err := parser.New(input)
	if err != nil {
		return nil, err
	}
	p.Rules["Path"] = Path
	p.Rules["PathItem"] = PathItem
	p.Rules["FormulaContent"] = FormulaContent
	return p, nil
}
