package subtle

import (
	"sundown/solution/palisade"
	"sundown/solution/prism"
)

func (env Environment) AnalyseMonadic(m *palisade.Monadic) prism.MApplication {
	op := env.FetchVerb(m.Verb)
	if _, ok := op.(prism.MonadicFunction); !ok {
		panic("Verb is not a monadic function")
	}

	fn := op.(prism.MonadicFunction)
	expr := env.AnalyseExpression(m.Expression)

	if fn.OmegaType.SemiDetermined() {
		fn.OmegaType = expr.Type()

		if fn.Returns.SemiDetermined() {
			fn.Returns = expr.Type()
		}
	} else if !prism.EqType(expr.Type(), fn.OmegaType) {
		panic("Type mismatch")
	}

	return prism.MApplication{
		Operator: fn,
		Operand:  expr,
	}
}
func (env Environment) AnalyseDyadic(d *palisade.Dyadic) prism.DApplication {
	op := env.FetchVerb(d.Verb)
	if _, ok := op.(prism.DyadicFunction); !ok {
		panic("Verb is not a dyadic function")
	}
	var left prism.Expression
	if d.Monadic != nil {
		left = env.AnalyseMonadic(d.Monadic)
	} else if d.Morphemes != nil {
		left = env.AnalyseMorphemes(d.Morphemes)
	} else {
		panic("Dyadic expression has no left operand")
	}

	right := env.AnalyseExpression(d.Expression)
	if !prism.EqType(left.Type(), op.(prism.DyadicFunction).AlphaType) {
		panic("Alpha type mismatch")
	} else if !prism.EqType(right.Type(), op.(prism.DyadicFunction).OmegaType) {
		panic("Omega type mismatch")
	}

	return prism.DApplication{
		Operator: op.(prism.DyadicFunction),
		Left:     left,
		Right:    right,
	}
}
