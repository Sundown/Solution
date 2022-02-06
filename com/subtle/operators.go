package subtle

import (
	"fmt"
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
	case "Map":
		if _, ok := lexpr.(prism.Function); !ok {
			panic("Left operand is not a function")
		}
		if _, ok := rexpr.Type().(prism.VectorType); !ok {
			panic("Right operand is not a vector")
		}

		dop = prism.DyadicOperator{
			Operator: prism.KindMapOperator,
			Left:     lexpr.(prism.MonadicFunction),
			Right:    rexpr,
			Returns:  rexpr.Type(), // TODO incorrect, actually [left.type()] but needs algebraic handling
		}
		tmp := dop.Left.(prism.MonadicFunction).OmegaType
		tmp2 := dop.Right.Type().(prism.VectorType).Type
		_, err := prism.Delegate(&tmp, &tmp2)
		if err != nil {
			panic(*err)
		}

		fmt.Println(dop.Type())
	case "/":
		if _, ok := lexpr.(prism.DyadicFunction); !ok {
			panic("Left operand is not a function")
		}
		if _, ok := rexpr.Type().(prism.VectorType); !ok {
			panic("Right operand is not a vector")
		}

		elmtype := rexpr.Type().(prism.VectorType).Type
		fn := lexpr.(prism.DyadicFunction)

		if !prism.PureMatch(elmtype, fn.OmegaType) {
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

		if !prism.PureMatch(elmtype, fn.AlphaType) {
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

		if prism.PredicateGenericType(fn.Returns) {
			fn.Returns = prism.IntegrateGenericType(fn.AlphaType, fn.Returns)
		}

		dop = prism.DyadicOperator{
			Operator: prism.KindFoldlOperator,
			Left:     lexpr.(prism.DyadicFunction),
			Right:    rexpr,
			Returns:  rexpr.Type().(prism.VectorType).Type, // wrong source
		}
	}

	return dop
}
