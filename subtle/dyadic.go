package subtle

import (
	"github.com/sundown.solution/palisade"
	"github.com/sundown.solution/prism"
)

func (env Environment) AnalyseDyadic(d *palisade.Dyadic) prism.DyadicApplication {
	var left prism.Expression
	if d.Monadic != nil {
		left = env.AnalyseMonadic(d.Monadic)
	} else if d.Morphemes != nil {
		left = env.AnalyseMorphemes(d.Morphemes)
	}

	right := env.AnalyseExpression(d.Expression)

	var fn prism.DyadicFunction
	if d.Verb == nil {
		fn = env.AnalyseDyadicPartial(d.Subexpr, left.Type(), right.Type())
	} else {
		fn = env.FetchDVerb(d.Verb)
		if !right.Type().Equals(fn.OmegaType) {
			if !prism.QueryCast(right.Type(), fn.OmegaType) {
				tmp := right.Type()
				_, err := prism.Delegate(&fn.OmegaType, &tmp)
				if err != nil {
					prism.Panic(*err)
				}
			} else {
				right = prism.DelegateCast(right, fn.OmegaType)
			}
		}

		if !left.Type().Equals(fn.AlphaType) {
			if !prism.QueryCast(left.Type(), fn.AlphaType) {
				tmp := left.Type()
				_, err := prism.Delegate(&fn.AlphaType, &tmp)
				if err != nil {
					prism.Panic(*err)
				}
			} else {
				left = prism.DelegateCast(left, fn.AlphaType)
			}
		}

		if _, err := prism.Delegate(&fn.AlphaType, &fn.OmegaType); err != nil {
			prism.Panic(*err)
		}

		if fn.Returns.IsAlgebraic() {
			fn.Returns = fn.Returns.Resolve(fn.OmegaType)
		}

		if fn.Name.Package == "_" && fn.Name.Name == "Return" {
			if !env.CurrentFunctionIR.Type().Equals(fn.Returns) {
				if !env.CurrentFunctionIR.Type().IsAlgebraic() {
					panic("Return recieves type which does not match determined-function's type")
				} else {
					panic("Not implemented, pain")
				}
			}
		}
	}

	return prism.DyadicApplication{
		Operator: fn,
		Left:     left,
		Right:    right,
	}
}
