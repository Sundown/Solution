package subtle

import (
	"github.com/sundown/solution/palisade"
	"github.com/sundown/solution/prism"
)

func (env Environment) analyseDyadicOperator(d *palisade.Operator) prism.DyadicOperator {
	dop := prism.DyadicOperator{}

	rexpr := env.analyseExpression(d.Expression)

	switch *d.Operator {
	case "¨":
		var lexpr prism.Expression
		if d.Verb != nil {
			lexpr = env.FetchMVerb(d.Verb)
		}
		if _, ok := rexpr.Type().(prism.VectorType); !ok {
			prism.Panic("Right operand is not a vector")
		}

		if d.Subexpr != nil {
			prism.Panic("Not implemented")
			//lexpr = env.analyseMonadicPartial(d.Subexpr, rexpr.Type().(prism.VectorType).Type)
		}

		if _, ok := lexpr.(prism.Function); !ok {
			prism.Panic("Left operand is not a function")
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
		var lexpr prism.Expression
		if d.Verb != nil {
			lexpr = env.FetchDVerb(d.Verb)
		}
		if _, ok := rexpr.Type().(prism.VectorType); !ok {
			prism.Panic("Right operand is not a vector")
		}

		if d.Subexpr != nil {
			t := rexpr.Type().(prism.VectorType).Type
			lexpr = env.analyseDyadicPartial(d.Subexpr, t, t)
		}

		if _, ok := lexpr.(prism.DyadicFunction); !ok {
			prism.Panic("Left operand is not a function")
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
			Operator: prism.KindReduceOperator,
			Left:     fn,
			Right:    rexpr,
			Returns:  fn.Type(),
		}
	}

	return dop
}
