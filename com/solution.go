package main

import (
	"sundown/solution/oversight"
	"sundown/solution/palisade"
	"sundown/solution/subtle"

	"github.com/alecthomas/repr"
)

func main() {
	oversight.Notify("Solution init...")

	r := &oversight.Runtime{}

	lexed := palisade.Begin(r.ParseArgs())
	env := subtle.Init(lexed)

	repr.Println(env)
}
