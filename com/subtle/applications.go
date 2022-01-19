package subtle

import (
	"sundown/solution/palisade"
	"sundown/solution/prism"
)

func (env Environment) AnalyseMonadic(m *palisade.Monadic) prism.MApplication {
	op := env.FetchVerb(m.Verb)
	if _, ok := op.(prism.MFunction); !ok {
		panic("Verb is not a monadic function")
	}

	expr := env.AnalyseExpression(m.Expression)
	if !prism.EqType(expr.Type(), op.(prism.MFunction).Returns) {
		panic("Type mismatch")
	}

	return prism.MApplication{
		Operator: op.(prism.MFunction),
		Operand:  expr, // TODO check type
	}
}
func (env Environment) AnalyseDyadic(d *palisade.Dyadic) prism.DApplication {
	op := env.FetchVerb(d.Verb)
	if _, ok := op.(prism.DFunction); !ok {
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
	if !prism.EqType(left.Type(), op.(prism.DFunction).AlphaType) {
		panic("Alpha type mismatch")
	} else if !prism.EqType(right.Type(), op.(prism.DFunction).OmegaType) {
		panic("Omega type mismatch")
	}

	return prism.DApplication{
		Operator: op.(prism.DFunction),
		Left:     left,
		Right:    right,
	}
}
