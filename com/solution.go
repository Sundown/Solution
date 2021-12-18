package main

import (
	"fmt"
	"sundown/solution/oversight"
	"sundown/solution/palisade"
	"sundown/solution/prescience"
	"sundown/solution/prism"
	"sundown/solution/subtle"
	"sundown/solution/weave"
)

func main() {
	oversight.Notify("Solution init...")

	r := &oversight.Runtime{}
	env := prism.NewEnvironment()

	// Short pass: lexing/tokenisation
	lexed := palisade.Begin(r.ParseArgs())

	// Short pass: intern all function declarations
	prescience.Init(env, lexed)

	// Long pass: analyse all function declarations
	weave.Init(env, lexed)

	// Long pass: analyse all application
	subtle.Init(env)

	// Short pass: LLVM code generation
	// TODO!

	fmt.Println(env.String())

}
