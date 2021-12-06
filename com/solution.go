package main

import (
	"sundown/solution/lexer"
	"sundown/solution/oversight"

	"github.com/alecthomas/repr"
)

func main() {
	oversight.Notify("Solution init...")

	r := &oversight.Runtime{}
	l := &lexer.State{}
	//p := &temporal.State{}
	//c := &compiler.State{Runtime: r}

	lexed := l.Lex(r.ParseArgs())

	repr.Println(lexed)

	//r.HandleEmit(c.Compile(p.Parse(lexed)))
}
