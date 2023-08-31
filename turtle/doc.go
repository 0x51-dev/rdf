package turtle

import (
	"github.com/0x51-dev/rdf/turtle/grammar"
	"github.com/0x51-dev/rdf/turtle/grammar/ir"
	"github.com/0x51-dev/upeg/parser/op"
	"strings"
)

type Document = ir.Document

func ParseDocument(doc string) (*Document, error) {
	if len(doc) == 0 {
		return &Document{}, nil
	}
	if !strings.HasSuffix(doc, "\n") {
		doc += "\n"
	}
	p, err := grammar.NewParser([]rune(doc))
	if err != nil {
		return nil, err
	}
	n, err := p.Parse(op.And{grammar.Document, op.EOF{}})
	if err != nil {
		return nil, err
	}
	return ir.ParseDocument(n)
}
