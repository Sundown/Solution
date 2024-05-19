package apotheosis

import (
	"github.com/sundown/solution/prism"

	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/value"
)

func (env *Environment) newInlineEqual(left prism.Value, right prism.Value) value.Value {
	switch left.Type.Kind() {
	case prism.RealType.ID:
		return env.Block.NewFCmp(enum.FPredOEQ, left.Value, right.Value)
	case prism.IntType.ID:
		return env.Block.NewICmp(enum.IPredEQ, left.Value, right.Value)
	case prism.CharType.ID:
		return env.Block.NewICmp(enum.IPredEQ, left.Value, right.Value)
	case prism.BoolType.ID:
		return env.Block.NewICmp(enum.IPredEQ, left.Value, right.Value)
	}

	prism.Panic("unreachable")
	panic("unlabelled error")
}
