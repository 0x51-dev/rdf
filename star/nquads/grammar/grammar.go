package grammar

import (
	nq "github.com/0x51-dev/rdf/nquads/grammar"
	nt "github.com/0x51-dev/rdf/ntriples/grammar"
	nts "github.com/0x51-dev/rdf/star/ntriples/grammar"
	"github.com/0x51-dev/upeg/parser"
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
			nts.Subject, nt.OWhitespace,
			nt.Predicate, nt.OWhitespace,
			nts.Object, nt.OWhitespace,
			op.Optional{Value: op.And{
				nq.GraphLabel, nt.OWhitespace,
			}},
			'.', op.Optional{Value: nt.Comment},
		},
	}
)

func NewParser(input []rune) (*parser.Parser, error) {
	return nts.NewParser(input)
}
