package apotheosis

import (
	"sundown/solution/oversight"
	"sundown/solution/prism"

	"github.com/llir/llvm/ir/value"
)

func (state *State) CompileTuple(tuple *prism.Morpheme) value.Value {
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

func (state *State) TupleGet(typ *prism.Type, real value.Value, index int) value.Value {
	if len(typ.Tuple) < index {
		oversight.Panic(oversight.CT_OOB, index, typ.String(), len(typ.Tuple))
	}

	if typ.Tuple == nil {
		oversight.Panic(oversight.CT_Unexpected, oversight.Yellow("tuple"), oversight.Yellow(typ.String()))
	}

	ptr := state.Block.NewGetElementPtr(
		typ.AsLLType(), real,
		I32(0), I32(int64(index)))

	if typ.Tuple[index].Atomic == nil {
		return ptr
	} else {
		return state.Block.NewLoad(typ.Tuple[index].AsLLType(), ptr)
	}
}
