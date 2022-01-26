package apotheosis

import (
	"sundown/solution/prism"

	"github.com/llir/llvm/ir/value"
)

func (env *Environment) CompileInlineAdd(alpha, omega Value) value.Value {
	switch alpha.Type.Kind() {

	case prism.RealType.ID:
		return env.Block.NewFAdd(alpha.Value, omega.Value)
	case prism.IntType.ID:
		return env.Block.NewAdd(alpha.Value, omega.Value)
	case prism.CharType.ID:
		return env.Block.NewAdd(alpha.Value, omega.Value)
	}

	panic("unreachable")
}
