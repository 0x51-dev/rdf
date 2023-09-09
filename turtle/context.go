package turtle

import (
	"fmt"
	"github.com/0x51-dev/rdf/ntriples"
)

type context struct {
	base     string
	prefixes map[string]string

	bnIndex int
	elIndex int

	// list of generated triples
	triples []ntriples.Triple
}

func (ctx *context) bn() ntriples.BlankNode {
	ctx.bnIndex++
	return ntriples.BlankNode(fmt.Sprintf("_:b%d", ctx.bnIndex))
}

func (ctx *context) el() ntriples.BlankNode {
	ctx.elIndex++
	return ntriples.BlankNode(fmt.Sprintf("_:el%d", ctx.elIndex))
}
