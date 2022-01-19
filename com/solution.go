package main

import (
	"fmt"
	"sundown/solution/apotheosis"
	"sundown/solution/oversight"
	"sundown/solution/palisade"
	"sundown/solution/subtle"
)

func main() {
	oversight.Notify("Solution init...")

	r := &oversight.Runtime{}
	s := &apotheosis.State{}
	lexed := palisade.Begin(r.ParseArgs())
	env := subtle.Init(lexed)

	ir := s.Compile(&env)
	fmt.Println(ir.String())
}
