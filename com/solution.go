package main

import (
	"os"
	"sundown/solution/compiler"
	"sundown/solution/lex"
	"sundown/solution/parse"
)

func main() {
	l := &lex.State{}
	p := &parse.State{}
	c := &compiler.State{}

	c.Compile(p.Parse(l.Lex(os.Args)))
}
