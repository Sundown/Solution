package main

import (
	"os"

	"github.com/sundown/solution/apotheosis"
	"github.com/sundown/solution/pilot"
	"github.com/sundown/solution/prism"
	"github.com/sundown/solution/subtle"
)

func main() {
	prism.Notify("Solution init...")

	if os.Args[1] == "pilot" {
		prism.Notify("Starting Pilot")
		pilot.Pilot()
		os.Exit(0)
	}

	env := prism.NewEnvironment()

	// Parse arguments
	prism.Init(env)

	// Open file and lex
	prism.Lex(env)

	// Parse lexed tokens to AST and resolve compiler directives
	subtle.Parse(env)

	// compile AST to LLVM
	apotheosis.Compile(env)

	// write LLVM IR to file or invoke Clang
	prism.Emit(env)
}
