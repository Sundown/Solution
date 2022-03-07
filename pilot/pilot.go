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
	Expr   bool
}

func newRun(t Test) {
	env := prism.NewEnvironment()
	env.IsPilotRun = true
	if t.Expr {
		env.File = `@Package pilot_test_output; @Entry Main; Δ Main Int → Void: ` + t.Code + `∇`
	} else {
		env.File = t.Code
	}
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
