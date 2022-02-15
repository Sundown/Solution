package subtle

import (
	"sundown/solution/palisade"
	"sundown/solution/prism"
)

func (env Environment) AnalyseDyadicOperator(d *palisade.Monadic) prism.DyadicOperator {
	dop := prism.DyadicOperator{}
	var lexpr prism.Expression
	if d.Verb != nil {
		lexpr = env.FetchVerb(d.Verb)
	} else if d.Subexpr != nil {
		if d.Subexpr.Dyadic == nil {
			panic("unreachable")
		}

		lexpr = env.AnalysePartial(d.Subexpr.Dyadic)
	} else {
		panic("Dyadic expression has no left operand")
	}

	rexpr := env.AnalyseExpression(d.Expression.Monadic.Expression)

	switch *d.Expression.Monadic.Verb.Ident {
	case "¨":
		if _, ok := lexpr.(prism.Function); !ok {
			panic("Left operand is not a function")
		}
		if _, ok := rexpr.Type().(prism.VectorType); !ok {
			panic("Right operand is not a vector")
		}

		elmtype := rexpr.Type().(prism.VectorType).Type
		fn := lexpr.(prism.MonadicFunction)

		if !elmtype.Equals(fn.OmegaType) {
			if !prism.QueryCast(elmtype, fn.OmegaType) {
				tmp := elmtype
				_, err := prism.Delegate(&fn.OmegaType, &tmp)
				if err != nil {
					prism.Panic(*err)
				}
			} else {
				rexpr = prism.DelegateCast(rexpr, prism.VectorType{Type: fn.OmegaType})
			}
		}

		if fn.Returns.IsAlgebraic() {
			fn.Returns = fn.Returns.Resolve(fn.OmegaType)
		}

		dop = prism.DyadicOperator{
			Operator: prism.KindMapOperator,
			Left:     fn,
			Right:    rexpr,
			Returns:  fn.Type(),
		}
	case "/":
		if _, ok := lexpr.(prism.DyadicFunction); !ok {
			panic("Left operand is not a function")
		}
		if _, ok := rexpr.Type().(prism.VectorType); !ok {
			panic("Right operand is not a vector")
		}

		elmtype := rexpr.Type().(prism.VectorType).Type
		fn := lexpr.(prism.DyadicFunction)

		if !elmtype.Equals(fn.OmegaType) {
			if !prism.QueryCast(elmtype, fn.OmegaType) {
				tmp := elmtype
				_, err := prism.Delegate(&fn.OmegaType, &tmp)
				if err != nil {
					prism.Panic(*err)
				}
			} else {
				rexpr = prism.DelegateCast(rexpr, prism.VectorType{Type: fn.OmegaType})
			}
		}

		if !elmtype.Equals(fn.AlphaType) {
			if !prism.QueryCast(elmtype, fn.AlphaType) {
				tmp := elmtype
				_, err := prism.Delegate(&fn.AlphaType, &tmp)
				if err != nil {
					prism.Panic(*err)
				}
			} else {
				rexpr = prism.DelegateCast(rexpr, prism.VectorType{Type: fn.AlphaType})
			}
		}

		if _, err := prism.Delegate(&fn.AlphaType, &fn.OmegaType); err != nil {
			prism.Panic(*err)
		}

		if fn.Returns.IsAlgebraic() {
			fn.Returns = fn.Returns.Resolve(fn.AlphaType)
		}

		dop = prism.DyadicOperator{
			Operator: prism.KindFoldlOperator,
			Left:     fn,
			Right:    rexpr,
			Returns:  fn.Type(),
		}
	}

	return dop
}
