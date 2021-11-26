package compiler

import (
	"sundown/solution/oversight"
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

func (state *State) TupleGet(temporal *parse.Expression, real value.Value, index int) value.Value {
	if len(temporal.TypeOf.Tuple) < index {
		oversight.Panic(oversight.CT_OOB, index, temporal.TypeOf.String(), len(temporal.TypeOf.Tuple))
	}

	if temporal.TypeOf.Tuple == nil {
		oversight.Panic(oversight.CT_Unexpected, oversight.Yellow("tuple"), oversight.Yellow(temporal.TypeOf.String()))
	}

	return state.Block.NewLoad(
		temporal.TypeOf.Tuple[index].AsLLType(),
		state.Block.NewGetElementPtr(
			temporal.TypeOf.AsLLType(), real,
			I32(0), I32(int64(index))))
}
