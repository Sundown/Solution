package pilot

import (
	"github.com/sundown/solution/apotheosis"
	"github.com/sundown/solution/prism"
	"github.com/sundown/solution/subtle"
)

func Pilot() {
	for _, t := range Cases {
		newRun(t)
	}
}

type Test struct {
	Name   string
	Code   string
	Result string
}

func newRun(t Test) {
	env := prism.NewEnvironment()
	env.IsPilotRun = true

	env.File = `@Package pilot_test_output; @Entry Main; Main Int → Void{` + t.Code + `}`

	prism.Lex(env)
	subtle.Parse(env)
	apotheosis.Compile(env)
	res, p := prism.PilotEmit(env)

	if p && res == t.Result {
		prism.Notify(prism.Green("PASS"), t.Name)
	} else {
		prism.Error(prism.Red("FAIL"), t.Name, res)
	}

}
