package main

import (
	"github.com/sundown/solution/prism"
	"github.com/sundown/solution/subtle"
)

func main() {
	prism.Notify("Solution init...")

	env := prism.NewEnvironment()

	// Parse arguments
	prism.Init(env)

	// Open file and lex
	prism.Lex(env)

	// Parse lexed tokens to AST and resolve compiler directives
	subtle.Parse(env)

	// compile AST to LLVM
	apotheosis.Compile(env)

	// Write LLVM IR to file or invoke Clang
	prism.Emit(env)
}
