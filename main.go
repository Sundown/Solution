package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

var help = `Rib: the interactive compiler backend

Usage:
	rib [SUBCOMMAND] [PATH]

Subcommands:
	build <file>	Compiles input file to LLVM IR
	parse <file> 	Parses input file and prints syntax tree
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
		fmt.Println(parser)
		os.Exit(0)
	case "build":
		filecontents, err = ioutil.ReadFile(os.Args[2])
		if err != nil {
			panic(err)
		}

		prog := &Program{}

		err = parser.ParseString(os.Args[2], string(filecontents), prog)
		if err != nil {
			panic(err)
		}

		for _, expr := range prog.Expression {
			gen(expr)
		}

		print_module()
	default:
		pretty_error("invalid subcommand" + os.Args[1])
		os.Exit(1)
	}
}
