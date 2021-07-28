package main

import (
	"os"
	"sundown/sunday/compiler"
	"sundown/sunday/lex"
	"sundown/sunday/parse"
)

func main() {
	l := &lex.State{}
	p := &parse.State{}
	c := &compiler.State{}

	c.Compile(p.Parse(l.Lex(os.Args)))
}
