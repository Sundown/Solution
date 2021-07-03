package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"sundown/sunday/codegen"
	"sundown/sunday/parser"
	"sundown/sunday/util"
	"time"
)

var build = 1
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

		version = strconv.FormatInt(num, 36)

	} else {
		version = "unknown"
	}

	return version
}

func main() {
	GetVersion()
	start_time := time.Now()
	fmt.Printf(`#===================================#
#                                   #
#  Sunday compiler - build: %s  #
#  https://sundow.nl/sunday         #
#                                   #
#===================================#`+"\n", GetVersion()[0:6])

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
		fmt.Println("Parsing")
		prog := &parser.Program{}

		err = parser.Parser.ParseString(os.Args[2], string(filecontents), prog)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Parsed: %s in %s\n", path, time.Since(start_time).Round(1000))
		fmt.Println("Compiling")
		codegen.StartCompiler(path, prog)
	default:
		util.Error("invalid subcommand" + os.Args[1])
		os.Exit(1)
	}
}
