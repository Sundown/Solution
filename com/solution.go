package main

import (
	"fmt"
	"os"
	"sundown/solution/compiler"
	"sundown/solution/lex"
	"sundown/solution/parse"

	"github.com/alecthomas/kong"
)

var cli struct {
	EBNF  bool     `help"Dump EBNF."`
	Files []string `arg:"" optional:"" type:"existingfile" help:"GraphQL schema files to parse."`
}

func main() {
	/* 	util.Bene("Solution version " + util.GetSolutionVersion())
	   	util.VerifyClangVersion()

	   	rt := &util.Runtime{}
	   	if len(os.Args) == 1 {
	   		util.Error("No files input").Exit()
	   	}

	   	for i, s := range os.Args {
	   		if s[0] == '-' {
	   			switch s[1:] {
	   			case "emit":
	   				if len(os.Args) > i+1 {
	   					i++
	   					rt.Emit = util.Ref(os.Args[i])
	   				} else {
	   					util.Error("emit expected one of [bc, llvm, binary, asm]").Exit()
	   				}

	   			case "o":
	   				if len(os.Args) > i+1 {
	   					i++
	   					rt.Output = util.Ref(os.Args[i])
	   				} else {
	   					util.Error("output expected filename").Exit()
	   				}
	   			}
	   		}
	   	} */

	ctx := kong.Parse(&cli)
	if cli.EBNF {
		fmt.Println("Something")
		ctx.Exit(0)
	}
	file := cli.Files[0]
	l := &lex.State{}
	r, err := os.Open(file)
	ctx.FatalIfErrorf(err)
	err = lex.Parser.Parse(file, r, l)
	r.Close()

	ctx.FatalIfErrorf(err)

	p := &parse.State{}
	c := &compiler.State{}

	t := p.Parse(l)
	//fmt.Println(t.String())
	c.Compile(t)
}
