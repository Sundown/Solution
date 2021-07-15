package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"

	"sundown/sunday/ir"
	"sundown/sunday/parser"
	"sundown/sunday/util"

	"github.com/alecthomas/repr"
)

var help = `Sunday

Usage:
	sunday [SUBCOMMAND] [PATH]

Subcommands:
	build <file>	Compiles input file to LLVM IR
	grammar		Prints the Rib EBNF grammar
`

func GetVersion() string {
	var version string
	out, _ := exec.Command("git", "rev-parse", "--short", "HEAD").Output()
	if out != nil {
		num, err := strconv.ParseInt(string(out)[:len(out)-1], 16, 64)
		if err != nil {
			panic(err)
		}

		version = strconv.FormatInt(num, 36)[0:6]

	} else {
		version = "unknown"
	}

	return version
}

func main() {
	var filecontents []byte
	var err error

	if len(os.Args) < 2 {
		fmt.Println(help)
		os.Exit(0)
	}

	filecontents, err = ioutil.ReadFile(os.Args[2])
	if err != nil {
		panic(err)
	}

	prog := &parser.Program{}

	err = parser.Parser.ParseString(os.Args[2], string(filecontents), prog)
	ioutil.WriteFile("tree.yml", []byte(repr.String(prog, repr.Indent("	"))), 0644)
	if err != nil {
		panic(err)
	}

	switch os.Args[1] {
	case "analyse":
		s := &ir.State{}

		s.Analyse(prog)
		//repr.Println(res)
		fmt.Println(s.String())
	default:
		util.Error("invalid subcommand" + os.Args[1])
		os.Exit(1)
	}
}
