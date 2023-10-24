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
				ttl.Directive,
				Block,
				op.And{op.Optional{Value: nt.Comment}, op.EndOfLine{}},
			},
		}},
	}
	Block = op.Capture{
		Name: "Block",
		Value: op.Or{
			TriplesOrGraph,
			WrappedGraph,
			Triples,
			op.Capture{
				Name:  "Graph",
				Value: op.And{"GRAPH", nt.Whitespace, LabelOrSubject, ttl.WSPLNC, WrappedGraph},
			},
		},
	}
	TriplesOrGraph = op.Capture{
		Name: "TriplesOrGraph",
		Value: op.And{
			LabelOrSubject, ttl.WSPLNC,
			op.Or{
				WrappedGraph,
				op.And{ttl.PredicateObjectList, ttl.WSPLNC, '.'},
			},
		},
	}
	Triples = op.Capture{
		Name: "Triples2",
		Value: op.Or{
			op.Capture{
				Name: "Triples2BlankNodePropertyList",
				Value: op.And{
					ttl.BlankNodePropertyList,
					op.Optional{Value: op.And{ttl.WSPLNC, ttl.PredicateObjectList}},
					ttl.WSPLNC, '.',
				},
			},
			op.Capture{
				Name: "Triples2Collection",
				Value: op.And{
					ttl.Collection,
					ttl.WSPLNC,
					ttl.PredicateObjectList,
					ttl.WSPLNC, '.',
				},
			},
		},
	}
	WrappedGraph = op.Capture{
		Name: "WrappedGraph",
		Value: op.And{
			'{', ttl.WSPLNC, op.Optional{Value: op.And{TriplesBlock, ttl.WSPLNC}}, '}',
		},
	}
	TriplesBlock = op.Capture{
		Name: "TriplesBlock",
		Value: op.And{
			ttl.Triples,
			op.Optional{Value: op.And{
				ttl.WSPLNC, '.',
				op.Optional{Value: op.And{ttl.WSPLNC, op.Reference{Name: "TriplesBlock"}}},
			}},
		},
	}
	LabelOrSubject = op.Capture{
		Name: "LabelOrSubject",
		Value: op.Or{
			ttl.IRI, ttl.BlankNode,
		},
	}
)

func NewParser(input []rune) (*parser.Parser, error) {
	p, err := ttl.NewParser(input)
	if err != nil {
		return nil, err
	}
	p.Rules["TriplesBlock"] = TriplesBlock
	return p, nil
}
