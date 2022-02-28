package apotheosis

import (
	"github.com/sundown/solution/prism"

	"github.com/llir/llvm/ir/value"
)

func (env *Environment) Apply(c prism.Callable, params ...prism.Value) value.Value {
	switch fn := c.(type) {
	case prism.DyadicFunction:
		if fn.Special {
			return env.Apply(env.FetchDCallable(fn.Ident().Name), params...)
		}
		v, ok := c.(prism.Expression)
		if !ok {
			prism.Panic("Apply: not an expression")
		}
		return env.Block.NewCall(env.compileExpression(&v), params[0].Value, params[1].Value)
	case prism.MonadicFunction:
		if fn.Special {
			return env.Apply(env.FetchMCallable(fn.Ident().Name), params...)
		}

		v, ok := c.(prism.Expression)
		if !ok {
			prism.Panic("Apply: not an expression")
		}
		return env.Block.NewCall(env.compileExpression(&v), params[0].Value)
	case prism.DCallable:
		return fn(params[0], params[1])
	case prism.MCallable:
		return fn(params[0])
	}

	prism.Panic("unreachable")
	panic(nil)
}
