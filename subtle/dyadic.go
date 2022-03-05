package subtle

import (
	"fmt"

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
	didAutoVector := false
	if d.Verb == nil {
		fn = env.analyseDyadicPartial(d.Subexpr, left.Type(), right.Type())
	} else {
		fn = env.FetchDVerb(d.Verb)
		avr := prism.QueryAutoVector(fn.OmegaType, right.Type())
		avl := prism.QueryAutoVector(fn.AlphaType, left.Type())
		fmt.Println(avr, avl)
		if !right.Type().Equals(fn.OmegaType) {
			if !fn.NoAutoVector() && avl && avr {
				fn.AlphaType = prism.VectorType{Type: fn.AlphaType}
				fn.OmegaType = prism.VectorType{Type: fn.OmegaType}
				fn.Returns = prism.VectorType{Type: fn.Returns}
				didAutoVector = true
			}

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

		if !left.Type().Equals(fn.AlphaType) && !didAutoVector {
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

		if _, err := prism.Delegate(&fn.AlphaType, &fn.OmegaType); !didAutoVector && err != nil {
			prism.Panic(*err)
		}

		if fn.Returns.IsAlgebraic() {
			fn.Returns = fn.Returns.Resolve(fn.OmegaType)
		}

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
		AutoVector: didAutoVector,
	}
}
