package main

import (
	"fmt"
	"os"
	"sundown/sunday/compiler"
	"sundown/sunday/lex"
	"sundown/sunday/parse"
)

func main() {
	l := &lex.State{}
	p := &parse.State{}
	c := &compiler.State{}

	p.Parse(l.Lex(os.Args))
	fmt.Println(p.String())

	c.Compile(p)
}
