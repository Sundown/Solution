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

	fn := env.analyseApplicable(*d.Applicable, left.Type(), right.Type())

	if _, ok := fn.(prism.DyadicFunction); !ok {
		prism.Panic("Expected dyadic function, got " + fn.String())
	}

	function := fn.(prism.DyadicFunction)

	prism.DeferDyadicApplicationTypes(&function, &left, &right)

	if function.Ident().Package == "_" && function.Ident().Name == "←" {
		if !env.CurrentFunctionIR.Type().Equals(function.Type()) {
			if !env.CurrentFunctionIR.Type().IsAlgebraic() {
				prism.Panic("Return receives type which does not match determined-function's type")
			} else {
				prism.Panic("Not implemented, pain")
			}
		}
	}

	return prism.DyadicApplication{
		Operator: function,
		Left:     left,
		Right:    right,
	}
}
