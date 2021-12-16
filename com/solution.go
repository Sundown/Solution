package main

import (
	"sundown/solution/oversight"
	"sundown/solution/palisade"
	"sundown/solution/weave"
)

func main() {
	oversight.Notify("Solution init...")

	r := &oversight.Runtime{}

	//p := &subtle.State{}
	//c := &compiler.State{Runtime: r}
	m := weave.State{}

	lexed := palisade.Begin(r.ParseArgs())

	//repr.Println(lexed)

	m.Init(lexed)

	//r.HandleEmit(c.Compile(p.Parse(lexed)))
}
