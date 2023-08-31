package grammar

import (
	nt "github.com/0x51-dev/rdf/ntriples/grammar"
	"github.com/0x51-dev/upeg/parser/op"
)

var (
	Document = op.Capture{
		Name: "Document",
		Value: op.ZeroOrMore{Value: op.And{
			op.Optional{Value: op.Or{Statement, nt.Comment}},
			nt.OWhitespace, op.EndOfLine{},
		}},
	}
	Statement = op.Capture{
		Name: "Statement",
		Value: op.And{
			nt.OWhitespace,
			nt.Subject, nt.OWhitespace,
			nt.Predicate, nt.OWhitespace,
			nt.Object, nt.OWhitespace,
			op.Optional{Value: op.And{
				GraphLabel, nt.OWhitespace,
			}},
			'.', op.Optional{Value: nt.Comment},
		},
	}
	GraphLabel = op.Capture{
		Name:  "GraphLabel",
		Value: op.Or{nt.IRIReference, nt.BlankNodeLabel},
	}
)
