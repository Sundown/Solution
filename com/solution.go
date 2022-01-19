package main

import (
	"fmt"
	"sundown/solution/oversight"
	"sundown/solution/palisade"
	"sundown/solution/subtle"
)

func main() {
	oversight.Notify("Solution init...")

	r := &oversight.Runtime{}

	lexed := palisade.Begin(r.ParseArgs())
	env := subtle.Init(lexed)

	fmt.Println(env.String())
}
