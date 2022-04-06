package apotheosis

import (
	"github.com/alecthomas/repr"
	"github.com/sundown/solution/prism"

	"github.com/llir/llvm/ir/value"
)

func (env *Environment) compileExpression(expr *prism.Expression) value.Value {
	switch t := (*expr).(type) {
	case prism.MonadicApplication:
		return env.compileMonadicApplication(&t)
	case prism.DyadicApplication:
		return env.compileDyadicApplication(&t)
	case prism.Morpheme:
		return env.compileAtom(&t)
	case prism.MonadicOperator:
		panic("Can't do solo operators yet")
	case prism.OperatorApplication:
		return env.compileMonadicOperator(&t)
	case prism.Function:
		return env.compileFunction(&t)
	case prism.Alpha:
		return env.CurrentFunction.Params[0]
	case prism.Omega:
		if len(env.CurrentFunction.Params) == 1 {
			return env.CurrentFunction.Params[0]
		}

		return env.CurrentFunction.Params[1]

	// TODO find out what's causing this
	case prism.Cast:
		return env.compileCast(t)
	case *prism.Cast:
		return env.compileCast(*t)
	}

	repr.Println(*expr)
	panic(expr)
}

func (env *Environment) compileMonadicOperator(dop *prism.OperatorApplication) value.Value {
	switch dop.Op.Operator {
	case prism.KindMapOperator:
		return env.compileInlineMap(
			dop.Op.Fn.(prism.MonadicFunction),
			prism.Value{Value: env.compileExpression(&dop.Expr), Type: dop.Expr.Type()})

	case prism.KindReduceOperator:
		return env.compileInlineReduce(
			dop.Op.Fn.(prism.DyadicFunction),
			prism.Value{Value: env.compileExpression(&dop.Expr), Type: dop.Expr.Type()})
	}

	prism.Panic("unreachable")
	panic(nil)
}
