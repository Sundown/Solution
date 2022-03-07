package subtle

import (
	"github.com/sundown/solution/palisade"
	"github.com/sundown/solution/prism"
)

func (env Environment) analyseDyadic(d *palisade.Dyadic) prism.DyadicApplication {
	var left prism.Expression
	if d.Monadic != nil {
		left = env.analyseMonadic(d.Monadic)
	} else if d.Morphemes != nil {
		left = env.analyseMorphemes(d.Morphemes)
	}

	right := env.analyseExpression(d.Expression)

	var fn prism.DyadicFunction

	// Has autoVectoring occured in this call?
	branchAutoVector := false

	if d.Verb == nil {
		fn = env.analyseDyadicPartial(d.Subexpr, left.Type(), right.Type())
	} else {
		fn = env.FetchDVerb(d.Verb)
		prism.DeferDyadicApplicationTypes(&fn, &left, &right)

		if fn.Name.Package == "_" && fn.Name.Name == "Return" {
			if !env.CurrentFunctionIR.Type().Equals(fn.Returns) {
				if !env.CurrentFunctionIR.Type().IsAlgebraic() {
					prism.Panic("Return receives type which does not match determined-function's type")
				} else {
					prism.Panic("Not implemented, pain")
				}
			}
		}
	}

	return prism.DyadicApplication{
		Operator:   fn,
		Left:       left,
		Right:      right,
		AutoVector: branchAutoVector,
	}
}
