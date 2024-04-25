package apotheosis

import (
	"github.com/sundown/solution/prism"

	"github.com/llir/llvm/ir/value"
)

func (env *Environment) newExpression(expr *prism.Expression) value.Value {
	switch t := (*expr).(type) {
	case prism.MonadicApplication:
		return env.newMonadicApplication(&t)
	case prism.DyadicApplication:
		return env.newDyadicApplication(&t)
	case prism.Morpheme:
		return env.newAtom(&t)
	case prism.MonadicOperator:
		panic("Can't do solo operators yet")
	case prism.OperatorApplication:
		return env.newMonadicOperator(&t)
	case prism.Function:
		return env.newFunction(&t)
	case prism.Alpha:
		return env.CurrentFunction.Params[0]
	case prism.Omega:
		if len(env.CurrentFunction.Params) == 1 {
			return env.CurrentFunction.Params[0]
		}

		return env.CurrentFunction.Params[1]

	case prism.Cast, *prism.Cast:
		panic("Casts don't exist yet")
	}

	panic(expr)
}

func (env *Environment) newMonadicOperator(dop *prism.OperatorApplication) value.Value {
	switch dop.Op.Operator {
	case prism.KindMapOperator:
		return env.newInlineMap(
			dop.Op.Fn.(prism.MonadicFunction),
			prism.Value{Value: env.newExpression(&dop.Expr), Type: dop.Expr.Type()})

	case prism.KindReduceOperator:
		return env.newInlineReduce(
			dop.Op.Fn.(prism.DyadicFunction),
			prism.Value{Value: env.newExpression(&dop.Expr), Type: dop.Expr.Type()})
	}

	prism.Panic("unreachable")
	panic(nil)
}
