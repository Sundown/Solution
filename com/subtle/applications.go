package subtle

import (
	"sundown/solution/palisade"
	"sundown/solution/prism"
)

func (env Environment) AnalyseMonadic(m *palisade.Monadic) (app prism.MApplication) {
	op := env.FetchVerb(m.Verb)
	if _, ok := op.(prism.MonadicFunction); !ok {
		panic("Verb is not a monadic function")
	}

	fn := op.(prism.MonadicFunction)
	expr := env.AnalyseExpression(m.Expression)

	if !prism.PrimativeTypeEq(expr.Type(), fn.OmegaType) {
		if derived := prism.DeriveSemiDeterminedType(fn.OmegaType, expr.Type()); derived != nil {
			integrated_omega := prism.IntegrateSemiDeterminedType(derived, fn.OmegaType)

			fn.OmegaType = integrated_omega

			if prism.PredicateSemiDeterminedType(fn.Returns) {
				integrated_return := prism.IntegrateSemiDeterminedType(derived, fn.Returns)

				fn.Returns = integrated_return
			}
		} else {
			panic("Omega type mismatch between" + fn.OmegaType.String() + " and " + expr.Type().String())
		}
	}

	if fn.Name.Package == "_" && fn.Name.Name == "Return" {
		if !prism.PrimativeTypeEq(env.CurrentFunctionIR.Type(), fn.Returns) {
			if !prism.PredicateSemiDeterminedType(env.CurrentFunctionIR.Type()) {
				panic("Return recieves type which does not match determined-function's type")
			} else {
				panic("Not implemented, pain")
			}
		}
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

	fn := op.(prism.DyadicFunction)
	if !prism.PrimativeTypeEq(left.Type(), fn.AlphaType) {
		if derived := prism.DeriveSemiDeterminedType(fn.AlphaType, left.Type()); derived != nil {
			fn.AlphaType = prism.IntegrateSemiDeterminedType(derived, fn.AlphaType)
			if prism.PredicateSemiDeterminedType(fn.Returns) {
				fn.Returns = prism.IntegrateSemiDeterminedType(derived, fn.Returns)
			}
		} else {
			panic("Alpha type mismatch between " + fn.AlphaType.String() + " and " + right.Type().String())
		}
	} else if !prism.PrimativeTypeEq(right.Type(), fn.OmegaType) {
		if derived := prism.DeriveSemiDeterminedType(fn.OmegaType, right.Type()); derived != nil {
			fn.OmegaType = prism.IntegrateSemiDeterminedType(derived, fn.OmegaType)
			if prism.PredicateSemiDeterminedType(fn.Returns) {
				fn.Returns = prism.IntegrateSemiDeterminedType(derived, fn.Returns)
			}
		} else {
			panic("Omega type mismatch between " + fn.OmegaType.String() + " and " + right.Type().String())
		}
	}

	if fn.Name.Package == "_" && fn.Name.Name == "Return" {
		if !prism.PrimativeTypeEq(env.CurrentFunctionIR.Type(), fn.Returns) {
			if !prism.PredicateSemiDeterminedType(env.CurrentFunctionIR.Type()) {
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
