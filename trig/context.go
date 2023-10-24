package trig

import (
	"fmt"
	nt "github.com/0x51-dev/rdf/ntriples"
	ttl "github.com/0x51-dev/rdf/turtle"
)

type Context struct {
	*ttl.Context
}

func NewContext() *Context {
	return &Context{
		Context: ttl.NewContext(),
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
