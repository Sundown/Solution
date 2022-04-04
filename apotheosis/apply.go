package apotheosis

import (
	"github.com/sundown/solution/prism"

	"github.com/llir/llvm/ir/value"
)

func (env *Environment) apply(c prism.Callable, params ...prism.Value) value.Value {
	switch fn := c.(type) {
	case prism.DyadicFunction:
		if fn.Special {
			return env.apply(env.FetchDyadicCallable(fn.Ident().Name), params...)
		}
		v, ok := c.(prism.Expression)
		if !ok {
			prism.Panic("apply: not an expression")
		}

		if prism.IsVector(params[0].Type) && prism.IsVector(params[1].Type) && !fn.NoAutoVector() {
			return env.combineOf(fn, params[0], params[1])
		}

		return env.Block.NewCall(env.compileExpression(&v), params[0].Value, params[1].Value)
	case prism.MonadicFunction:
		if fn.Special {
			return env.apply(env.FetchMonadicCallable(fn.Ident().Name), params...)
		}

		v, ok := c.(prism.Expression)
		if !ok {
			prism.Panic("apply: not an expression")
		}

		if prism.IsVector(params[0].Type) && !fn.NoAutoVector() {
			return env.compileInlineMap(fn, params[0])
		}

		return env.Block.NewCall(env.compileExpression(&v), params[0].Value)
	case prism.DyadicCallable:
		if prism.IsVector(params[0].Type) && prism.IsVector(params[1].Type) && !c.NoAutoVector() {
			return env.combineOf(fn, params[0], params[1])
		}
		return fn.DCallable(params[0], params[1])
	case prism.MonadicCallable:
		return fn.MCallable(params[0])
	}

	prism.Panic("unreachable")
	panic(nil)
}
