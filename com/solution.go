package main

import (
	"sundown/solution/lexer"
	"sundown/solution/oversight"
	"sundown/solution/reform"
)

func main() {
	oversight.Notify("Solution init...")

	r := &oversight.Runtime{}
	l := &lexer.State{}
	//p := &temporal.State{}
	//c := &compiler.State{Runtime: r}
	m := reform.State{}

	lexed := l.Lex(r.ParseArgs())

	//repr.Println(lexed)

	m.Init(lexed)

	//r.HandleEmit(c.Compile(p.Parse(lexed)))
}
