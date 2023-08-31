package nquads

import (
	"github.com/0x51-dev/rdf/nquads/grammar"
	"github.com/0x51-dev/rdf/nquads/grammar/ir"
	"github.com/0x51-dev/upeg/parser"
	"github.com/0x51-dev/upeg/parser/op"
	"strings"
)

type Document = ir.Document

func ParseDocument(doc string) (Document, error) {
	if len(doc) == 0 {
		return nil, nil
	}
	if !strings.HasSuffix(doc, "\n") {
		doc += "\n"
	}
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
