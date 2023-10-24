package turtle

import (
	"fmt"
	nt "github.com/0x51-dev/rdf/ntriples"
)

type Context struct {
	Base     string
	Prefixes map[string]string

	BnIndex, ElIndex int
}

func NewContext() *Context {
	return &Context{
		Prefixes: make(map[string]string),
	}
}

func (ctx *Context) bn() nt.BlankNode {
	ctx.BnIndex++
	return nt.BlankNode(fmt.Sprintf("b%d", ctx.BnIndex))
}

func (ctx *Context) el() nt.BlankNode {
	ctx.ElIndex++
	return nt.BlankNode(fmt.Sprintf("el%d", ctx.ElIndex))
}
