package subtle

import (
	"sundown/solution/palisade"
	"sundown/solution/prism"

	"github.com/alecthomas/repr"
)

func (env Environment) AnalyseDyadicOperator(d *palisade.Operator) prism.DyadicOperator {
	dop := prism.DyadicOperator{}

	rexpr := env.AnalyseExpression(d.Expression)

	switch *d.Operator {
	case "¨":
		var lexpr prism.Expression
		if d.Verb != nil {
			lexpr = env.FetchMVerb(d.Verb)
		}
		if _, ok := rexpr.Type().(prism.VectorType); !ok {
			panic("Right operand is not a vector")
		}

		if d.Subexpr != nil {
			panic("Not implemented")
			//lexpr = env.AnalyseMonadicPartial(d.Subexpr, rexpr.Type().(prism.VectorType).Type)
		}

		if _, ok := lexpr.(prism.Function); !ok {
			panic("Left operand is not a function")
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
			panic("Right operand is not a vector")
		}

		if d.Subexpr != nil {
			t := rexpr.Type().(prism.VectorType).Type
			lexpr = env.AnalyseDyadicPartial(d.Subexpr, t, t)
		}

		if _, ok := lexpr.(prism.DyadicFunction); !ok {
			repr.Println(lexpr)
			panic("Left operand is not a function")
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
