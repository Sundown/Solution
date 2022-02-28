package subtle

import (
	"github.com/sundown/solution/palisade"
	"github.com/sundown/solution/prism"
)

func (env Environment) analyseMonadic(d *palisade.Monadic) prism.MonadicApplication {
	right := env.analyseExpression(d.Expression)

	var fn prism.MonadicFunction
	if d.Verb == nil {
		fn = env.analyseMonadicPartial(d.Subexpr, right.Type())
	} else {
		fn = env.FetchMVerb(d.Verb)
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

	return prism.MonadicApplication{
		Operator: fn,
		Operand:  right,
	}
}
