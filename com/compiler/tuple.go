package compiler

import (
	"sundown/solution/parse"

	"github.com/llir/llvm/ir/value"
)

func (state *State) CompileTuple(tuple *parse.Atom) value.Value {
	tuple_type := tuple.TypeOf.AsLLType()
	ll_tuple := state.Block.NewAlloca(tuple_type)

	for index, expr := range tuple.Tuple {
		val := state.CompileExpression(expr)

		if expr.TypeOf.Atomic == nil {
			val = state.Block.NewLoad(expr.TypeOf.AsLLType(), val)
		}

		state.Block.NewStore(val, state.GEP(ll_tuple, I32(0), I32(int64(index))))
	}

	return ll_tuple
}
