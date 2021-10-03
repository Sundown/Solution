package main

import (
	"sundown/solution/compiler"
	"sundown/solution/lexer"
	"sundown/solution/parse"
	"sundown/solution/util"
)

func main() {
	util.Notify("Solution init...")

	r := &util.Runtime{}
	l := &lexer.State{}
	p := &parse.State{}
	c := &compiler.State{Runtime: r}

	r.ParseArgs()

	r.HandleEmit(c.Compile(p.Parse(l.Lex(r))))
}
