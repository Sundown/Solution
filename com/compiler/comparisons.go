package compiler

import (
	"sundown/solution/temporal"

	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/value"
)

func (state *State) CompileInlineEqual(
	typAlpha *temporal.Type, valAlpha value.Value,
	typOmega *temporal.Type, valOmega value.Value) value.Value {
	return state.Block.NewICmp(enum.IPredEQ, valAlpha, valOmega)
}
