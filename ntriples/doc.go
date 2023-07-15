package ntriples

import (
	"github.com/0x51-dev/rdf/ntriples/grammar"
	"github.com/0x51-dev/rdf/ntriples/grammar/ir"
	"github.com/0x51-dev/upeg/parser"
	"github.com/0x51-dev/upeg/parser/op"
)

type Document = []Triple

func ParseDocument(doc string) (Document, error) {
	p, err := parser.New([]rune(doc))
	if err != nil {
		return nil, err
	}
	n, err := p.Parse(op.And{grammar.Document, op.EOF{}})
	if err != nil {
		return nil, err
	}
	return ir.ParseDocument(n)
}
