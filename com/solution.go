package main

import (
	"sundown/solution/compiler"
	"sundown/solution/lexer"
	"sundown/solution/oversight"
	"sundown/solution/temporal"
)

func main() {
	oversight.Notify("Solution init...")

	r := &oversight.Runtime{}
	l := &lexer.State{}
	p := &temporal.State{}
	c := &compiler.State{Runtime: r}

	lexed := l.Lex(r.ParseArgs())

	//repr.Println(lexed)

	r.HandleEmit(c.Compile(p.Parse(lexed)))
}
