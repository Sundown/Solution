package apotheosis

import (
	"sundown/solution/subtle"

	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/value"
)

func (state *State) CompileInlineEqual(
	typAlpha *subtle.Type, valAlpha value.Value,
	typOmega *subtle.Type, valOmega value.Value) value.Value {
	return state.Block.NewICmp(enum.IPredEQ, valAlpha, valOmega)
}
