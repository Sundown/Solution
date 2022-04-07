package subtle

import (
	"github.com/sundown/solution/palisade"
	"github.com/sundown/solution/prism"
)

func (env Environment) analyseMonadic(d *palisade.Monadic) prism.MonadicApplication {
	right := env.analyseExpression(d.Expression)

	fn := env.analyseApplicable(*d.Applicable, nil, right.Type())

	if _, ok := fn.(prism.MonadicFunction); !ok {
		prism.Panic("Expected monadic function, got " + fn.String())
	}

	function := fn.(prism.MonadicFunction)

	prism.DeferMonadicApplicationTypes(&function, &right)

	if function.Name.Package == "_" && function.Name.Name == "Return" {
		if !env.CurrentFunctionIR.Type().Equals(function.Returns) {
			if !env.CurrentFunctionIR.Type().IsAlgebraic() {
				prism.Panic("Return receives type which does not match determined-function's type")
			} else {
				prism.Panic("Not implemented, pain")
			}
		}
	}

	return prism.MonadicApplication{
		Operator: function,
		Operand:  right,
	}
}
