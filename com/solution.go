package main

import (
	"sundown/solution/compiler"
	"sundown/solution/lexer"
	"sundown/solution/oversight"
	"sundown/solution/parse"
)

func main() {
	oversight.Notify("Solution init...")

	r := &oversight.Runtime{}
	l := &lexer.State{}
	p := &parse.State{}
	c := &compiler.State{Runtime: r}

	r.ParseArgs()

	r.HandleEmit(c.Compile(p.Parse(l.Lex(r))))
}
