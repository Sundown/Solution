package apotheosis

import (
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/value"
)

func (state *State) CompileInlineEqual(left Value, right Value) value.Value {
	return state.Block.NewICmp(enum.IPredEQ, left.Value, right.Value)
}
