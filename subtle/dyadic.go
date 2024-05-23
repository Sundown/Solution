package subtle

import (
	"github.com/sundown/solution/palisade"
	"github.com/sundown/solution/prism"
)

func (env *Environment) analyseDBody(f *prism.DyadicFunction) {
	if _, okr := f.OmegaType.(prism.Universal); okr || f.Attrs().Special || f.Attrs().SkipBuilder {
		return
	}

	if _, okl := f.AlphaType.(prism.Universal); okl {
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

func (env *Environment) analyseDyadic(d *palisade.Dyadic) prism.DyadicApplication {
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

	function, left, right = prism.DeferDyadicApplicationTypes(function, left, right)

	return prism.DyadicApplication{
		Operator: function,
		Left:     left,
		Right:    right,
	}
}
