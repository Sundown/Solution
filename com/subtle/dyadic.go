package subtle

import (
	"sundown/solution/palisade"
	"sundown/solution/prism"
)

func (env Environment) AnalyseDyadic(d *palisade.Dyadic) prism.DApplication {
	var left prism.Expression
	if d.Monadic != nil {
		left = env.AnalyseMonadic(d.Monadic)
	} else if d.Morphemes != nil {
		left = env.AnalyseMorphemes(d.Morphemes)
	}

	right := env.AnalyseExpression(d.Expression)

	var fn prism.DyadicFunction
	if d.Verb == nil {
		fn = env.AnalyseDyadicPartial(d.Subexpr, left, right)
	} else {
		fn = env.FetchDVerb(d.Verb)

		if !prism.LoTypeEq(right.Type(), fn.OmegaType) {
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

		if !prism.LoTypeEq(left.Type(), fn.AlphaType) {
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

		if prism.PredicateGenericType(fn.Returns) {
			fn.Returns = prism.IntegrateGenericType(fn.AlphaType, fn.Returns)
		}

		if fn.Name.Package == "_" && fn.Name.Name == "Return" {
			if !prism.LoTypeEq(env.CurrentFunctionIR.Type(), fn.Returns) {
				if !prism.PredicateGenericType(env.CurrentFunctionIR.Type()) {
					panic("Return recieves type which does not match determined-function's type")
				} else {
					panic("Not implemented, pain")
				}
			}
		}
	}

	return prism.DApplication{
		Operator: fn,
		Left:     left,
		Right:    right,
	}
}