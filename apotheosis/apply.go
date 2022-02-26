package apotheosis

import (
	"sundown/solution/prism"

	"github.com/llir/llvm/ir/value"
)

func (env *Environment) Apply(c interface{}, params ...Value) value.Value {
	switch c.(type) {
	case prism.DyadicFunction:
		if c.(prism.DyadicFunction).Special {
			id := c.(prism.DyadicFunction).Ident()
			return env.GetSpecialDCallable(&id)(params[0], params[1])
		}
		v, ok := c.(prism.Expression)
		if !ok {
			panic("Apply: not an expression")
		}
		return env.Block.NewCall(env.CompileExpression(&v), params[0].Value, params[1].Value)
	case prism.MonadicFunction:
		if c.(prism.MonadicFunction).Special {
			id := c.(prism.MonadicFunction).Ident()
			return env.GetSpecialMCallable(&id)(params[0])
		}

		v, ok := c.(prism.Expression)
		if !ok {
			panic("Apply: not an expression")
		}
		return env.Block.NewCall(env.CompileExpression(&v), params[0].Value)
	case DCallable:
		return c.(DCallable)(params[0], params[1])
	}

	panic("unreachable")
}
