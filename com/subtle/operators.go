package subtle

import (
	"sundown/solution/palisade"
	"sundown/solution/prism"
)

func (env Environment) AnalyseDyadicOperator(d *palisade.Operator) prism.DyadicOperator {
	dop := prism.DyadicOperator{}

	lexpr := env.FetchVerb(d.Verb)

	rexpr := env.AnalyseExpression(d.Expression)

	switch *d.Operator {
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
