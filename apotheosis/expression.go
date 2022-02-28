package apotheosis

import (
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
	case prism.DyadicOperator:
		return env.compileDyadicOperator(&t)
	case prism.Function:
		return env.compileFunction(&t)
	case prism.Alpha:
		return env.CurrentFunction.Params[0]
	case prism.Omega:
		if len(env.CurrentFunction.Params) == 1 {
			return env.CurrentFunction.Params[0]
		} else {
			return env.CurrentFunction.Params[1]
		}
	case prism.Cast:
		return env.compileCast(t)
	default:
		prism.Panic("unreachable")
	}
	panic(nil)
}

func (env *Environment) compileDyadicOperator(dop *prism.DyadicOperator) value.Value {
	switch dop.Operator {
	case prism.KindMapOperator:
		return env.compileInlineMap(
			dop.Left.(prism.MonadicFunction),
			prism.Value{Value: env.compileExpression(&dop.Right), Type: dop.Right.Type()})

	case prism.KindReduceOperator:
		return env.compileInlineReduce(
			dop.Left.(prism.DyadicFunction),
			prism.Value{Value: env.compileExpression(&dop.Right), Type: dop.Right.Type()})
	}
	prism.Panic("unreachable")
	panic(nil)
}
