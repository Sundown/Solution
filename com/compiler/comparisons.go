package compiler

import (
	"sundown/solution/temporal"

	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/value"
)

func (state *State) CompileInlineEqual(typ *temporal.Type, val value.Value) value.Value {
	return state.Block.NewICmp(enum.IPredEQ,
		state.TupleGet(typ, val, 0),
		state.TupleGet(typ, val, 1))
}
