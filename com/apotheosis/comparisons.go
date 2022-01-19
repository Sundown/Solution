package apotheosis

import (
	"sundown/solution/prism"

	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/value"
)

func (state *State) CompileInlineEqual(
	typAlpha *prism.Type, valAlpha value.Value,
	typOmega *prism.Type, valOmega value.Value) value.Value {
	return state.Block.NewICmp(enum.IPredEQ, valAlpha, valOmega)
}
