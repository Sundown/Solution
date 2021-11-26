package compiler

import (
	"sundown/solution/temporal"

	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/value"
)

func (state *State) CompileInlineEqual(arg *temporal.Expression) value.Value {
	p := state.CompileExpression(arg)
	return state.Block.NewICmp(enum.IPredEQ,
		state.TupleGet(arg.TypeOf, p, 0),
		state.TupleGet(arg.TypeOf, p, 1))
}
