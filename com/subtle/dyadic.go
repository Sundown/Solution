package subtle

import (
	"fmt"
	"sundown/solution/palisade"
	"sundown/solution/prism"
)

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

	fn := op.(prism.DyadicFunction)

	tmp := right.Type()
	resolved_right, err := prism.Delegate(&fn.OmegaType, &tmp)
	if err != nil {
		fmt.Println(tmp.String())
		prism.Panic(*err)
	}
	tmp = left.Type()
	resolved_left, err := prism.Delegate(&fn.AlphaType, &tmp)
	if err != nil {
		prism.Panic(*err)
	}

	if _, err := prism.Delegate(resolved_left, resolved_right); err != nil {
		prism.Panic(*err)
	}

	if prism.PredicateGenericType(fn.Returns) {
		fn.Returns = prism.IntegrateGenericType(*resolved_left, fn.Returns)
	}

	if fn.Name.Package == "_" && fn.Name.Name == "Return" {
		if !prism.PrimativeTypeEq(env.CurrentFunctionIR.Type(), fn.Returns) {
			if !prism.PredicateGenericType(env.CurrentFunctionIR.Type()) {
				panic("Return recieves type which does not match determined-function's type")
			} else {
				panic("Not implemented, pain")
			}
		}
	}

	return prism.DApplication{
		Operator: fn,
		Left:     left,
		Right:    right,
	}
}
