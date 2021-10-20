package compiler

import (
	"sundown/solution/parse"

	"github.com/llir/llvm/ir/value"
)

func (state *State) CompileTuple(tuple *parse.Atom) value.Value {
	ll_tuple := state.Block.NewAlloca(tuple.TypeOf.AsLLType())

	for index, expr := range tuple.Tuple {
		val := state.CompileExpression(expr)

		if expr.TypeOf.Atomic == nil {
			val = state.Block.NewLoad(expr.TypeOf.AsLLType(), val)
		}

		state.Block.NewStore(val, state.GEP(ll_tuple, I32(0), I32(int64(index))))
	}

	return ll_tuple
}

func (state *State) InlineFirst(arg *parse.Expression) value.Value {
	val := state.Block.NewGetElementPtr(arg.TypeOf.AsLLType(),
		state.CompileExpression(arg), I32(0), I32(0))

	if arg.TypeOf.Tuple[0].Atomic != nil {
		return state.Block.NewLoad(arg.TypeOf.Tuple[0].AsLLType(), val)
	}

	return val
}

func (state *State) InlineSecond(arg *parse.Expression) value.Value {
	val := state.Block.NewGetElementPtr(arg.TypeOf.AsLLType(),
		state.CompileExpression(arg), I32(0), I32(1))

	if arg.TypeOf.Tuple[1].Atomic != nil {
		return state.Block.NewLoad(arg.TypeOf.Tuple[1].AsLLType(), val)
	}

	return val
}

func (state *State) InlineThird(arg *parse.Expression) value.Value {
	val := state.Block.NewGetElementPtr(arg.TypeOf.AsLLType(),
		state.CompileExpression(arg), I32(0), I32(2))

	if arg.TypeOf.Tuple[2].Atomic != nil {
		return state.Block.NewLoad(arg.TypeOf.Tuple[2].AsLLType(), val)
	}

	return val
}
