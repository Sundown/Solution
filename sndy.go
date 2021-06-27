package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sundown/sunday/codegen"
	"sundown/sunday/parser"
	"sundown/sunday/util"
	"time"
)

var help = `Sunday

Usage:
	sunday [SUBCOMMAND] [PATH]

Subcommands:
	build <file>	Compiles input file to LLVM IR
	grammar		Prints the Rib EBNF grammar
`

func main() {
	start_time := time.Now()
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

		path := "out.ll"

		prog := &parser.Program{}

		err = parser.Parser.ParseString(os.Args[2], string(filecontents), prog)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Parsed %s in %s\n", path, time.Since(start_time).Round(1000))
		codegen.StartCompiler(path, prog)
	default:
		util.Error("invalid subcommand" + os.Args[1])
		os.Exit(1)
	}
}
