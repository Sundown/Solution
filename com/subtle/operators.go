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
	case "Map":
		if _, ok := lexpr.(prism.Function); !ok {
			panic("Left operand is not a function")
		}
		if _, ok := rexpr.(prism.Vector); !ok {
			panic("Right operand is not a vector")
		}

		dop = prism.DyadicOperator{
			Operator: prism.KindMapOperator,
			Left:     lexpr.(prism.MonadicFunction),
			Right:    rexpr.(prism.Vector),
			Returns:  rexpr.Type(), // TODO incorrect, actually [left.type()] but needs algebraic handling
		}
		tmp := dop.Left.(prism.MonadicFunction).OmegaType
		tmp2 := dop.Right.Type().(prism.VectorType).Type
		_, err := prism.Delegate(&tmp, &tmp2)
		if err != nil {
			panic(*err)
		}
	case "/":
		if _, ok := lexpr.(prism.Function); !ok {
			panic("Left operand is not a function")
		}
		if _, ok := rexpr.Type().(prism.VectorType); !ok {
			panic("Right operand is not a vector")
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
