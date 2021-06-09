package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sundown/girl/codegen"
	"sundown/girl/parser"
)

var help = `Girl: the interactive compiler backend

Usage:
	Girl [SUBCOMMAND] [PATH]

Subcommands:
	build <file>	Compiles input file to LLVM IR
	grammar		Prints the Rib EBNF grammar
`

func main() {
	var filecontents []byte
	var err error

	if len(os.Args) < 2 {
		fmt.Println(help)
		os.Exit(0)
	}

	switch os.Args[1] {
	case "grammar":
		fmt.Println(parser.Parser)
		os.Exit(0)
	case "build":
		filecontents, err = ioutil.ReadFile(os.Args[2])
		if err != nil {
			panic(err)
		}

		prog := &parser.Program{}

		err = parser.Parser.ParseString(os.Args[2], string(filecontents), prog)
		if err != nil {
			panic(err)
		}

		codegen.StartCompiler("out.ll", prog)
	default:
		pretty_error("invalid subcommand" + os.Args[1])
		os.Exit(1)
	}
}
