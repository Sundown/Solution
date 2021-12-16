package main

import (
	"sundown/solution/oversight"
	"sundown/solution/palisade"
	"sundown/solution/prescience"
	"sundown/solution/prism"
	"sundown/solution/weave"

	"github.com/alecthomas/repr"
)

func main() {
	oversight.Notify("Solution init...")

	r := &oversight.Runtime{}
	env := prism.NewEnvironment()
	//p := &subtle.State{}
	//c := &compiler.State{Runtime: r}
	m := weave.State{}

	lexed := palisade.Begin(r.ParseArgs())

	oracle := prescience.Init(env, lexed)
	repr.Println(oracle)

	m.Init(lexed)

	//r.HandleEmit(c.Compile(p.Parse(lexed)))
}
