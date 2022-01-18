package main

import (
	"sundown/solution/oversight"
	"sundown/solution/palisade"
	"sundown/solution/prism"
	"sundown/solution/subtle"

	"github.com/alecthomas/repr"
)

func main() {
	oversight.Notify("Solution init...")

	r := &oversight.Runtime{}
	env := prism.NewEnvironment()

	lexed := palisade.Begin(r.ParseArgs())
	subtle.Init(env)

	repr.Println(lexed)
}
