package compiler

import (
	"sundown/solution/parse"

	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/value"
)

func (state *State) CompileInlineEqual(arg *parse.Expression) value.Value {
	p := state.CompileExpression(arg)
	return state.Block.NewICmp(enum.IPredEQ,
		state.TupleGet(arg, p, 0),
		state.TupleGet(arg, p, 1))
}
