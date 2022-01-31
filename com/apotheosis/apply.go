package apotheosis

import (
	"sundown/solution/prism"

	"github.com/llir/llvm/ir/value"
)

func (env *Environment) Apply(e *prism.Expression, params ...Value) value.Value {
	switch (*e).(type) {
	case prism.DyadicFunction:
		if (*e).(prism.DyadicFunction).Special {
			id := (*e).(prism.DyadicFunction).Ident()
			return env.GetSpecialDCallable(&id)(params[0], params[1])
		}

		return env.Block.NewCall(env.CompileExpression(e), params[0].Value, params[1].Value)
	case prism.MonadicFunction:
		if (*e).(prism.MonadicFunction).Special {
			id := (*e).(prism.MonadicFunction).Ident()
			return env.GetSpecialMCallable(&id)(params[0])
		}

		return env.Block.NewCall(env.CompileExpression(e), params[0].Value)
	}
	panic("unreachable")
}
