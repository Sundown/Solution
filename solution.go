package main

import (
	"os"

	"github.com/sundown/solution/apotheosis"
	"github.com/sundown/solution/pilot"
	"github.com/sundown/solution/prism"
	"github.com/sundown/solution/subtle"
)

func main() {
	if len(os.Args) >= 2 && os.Args[1] == "test" {
		prism.Notify("Starting Pilot Tests")
		pilot.Pilot()
		os.Exit(0)
	}

	prism.Notify("Solution init...")

	// Make environment and intern buildin functions
	env := prism.NewEnvironment()

	// Parse all arguments
	prism.Init(env)

	// Open file, lex, and close
	prism.Lex(env)

	// Parse lexed tokens to AST and resolve compiler directives
	subtle.Parse(env)

	// Compile AST to LLVM
	apotheosis.Compile(env)

	// Write LLVM IR to file or invoke Clang on LLVM Bitcode
	prism.Emit(env)
}
