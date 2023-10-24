package grammar

import (
	nt "github.com/0x51-dev/rdf/ntriples/grammar"
	"github.com/0x51-dev/upeg/parser"
	"github.com/0x51-dev/upeg/parser/op"
)

var (
	Document = op.Capture{
		Name: "Document",
		Value: op.ZeroOrMore{Value: op.And{
			op.Optional{Value: op.Or{Triple, nt.Comment}},
			nt.OWhitespace, op.EndOfLine{},
		}},
	}
	Triple = op.Capture{
		Name: "Triple",
		Value: op.And{
			nt.OWhitespace,
			Subject, nt.OWhitespace,
			nt.Predicate, nt.OWhitespace,
			Object, nt.OWhitespace,
			'.', op.Optional{Value: nt.Comment},
		},
	}
	Subject = op.Capture{
		Name:  "Subject",
		Value: op.Or{nt.IRIReference, nt.BlankNodeLabel, QuotedTriple},
	}
	Object = op.Capture{
		Name:  "Object",
		Value: op.Or{nt.IRIReference, nt.BlankNodeLabel, nt.Literal, QuotedTriple},
	}
	QuotedTriple = op.Capture{
		Name: "QuotedTriple",
		Value: op.And{
			"<<",
			nt.OWhitespace,
			op.Reference{Name: "Subject"}, nt.OWhitespace,
			nt.Predicate, nt.OWhitespace,
			op.Reference{Name: "Object"}, nt.OWhitespace,
			">>",
		},
	}
)

func NewParser(input []rune) (*parser.Parser, error) {
	p, err := parser.New(input)
	if err != nil {
		return nil, err
	}
	p.Rules["Subject"] = Subject
	p.Rules["Object"] = Object
	return p, nil
}
