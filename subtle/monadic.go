package subtle

import (
	"fmt"

	"github.com/sundown/solution/palisade"
	"github.com/sundown/solution/prism"
)

func (env *Environment) analyseMBody(f *prism.MonadicFunction) {
	if _, ok := f.OmegaType.(prism.Universal); ok || f.Attrs().Special || f.Attrs().SkipBuilder {
		return
	}

	t := env.CurrentFunctionIR
	env.CurrentFunctionIR = *f

	if len(f.Body) == 0 {
		for _, expr := range *f.PreBody {
			f.Body = append(f.Body, env.analyseExpression(&expr))
		}
	}

	env.CurrentFunctionIR = t
}

func (env *Environment) analyseMonadic(d *palisade.Monadic) prism.MonadicApplication {
	right := env.analyseExpression(d.Expression)
	fn := env.analyseApplicable(*d.Applicable, nil, right.Type())

	if _, ok := fn.(prism.MonadicFunction); !ok {
		prism.Panic("Expected monadic function, got " + fn.String())
	}

	function := fn.(prism.MonadicFunction)

	prism.DeferMonadicApplicationTypes(&function, &right)

	if isReturn(function) {
		if !env.CurrentFunctionIR.Type().Equals(function.Returns) {
			if !env.CurrentFunctionIR.Type().IsAlgebraic() {
				prism.Panic(fmt.Sprintf("Return receives type (%s) which does not match determined-function's type (%s)", function.Returns, env.CurrentFunctionIR.Type()))
			}
		}
	}

	return prism.MonadicApplication{
		Operator: function,
		Operand:  right,
	}
}
