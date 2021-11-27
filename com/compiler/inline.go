package compiler

import (
	"sundown/solution/temporal"

	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

func (state *State) CompileInlinePrintln(typ *temporal.Type, val value.Value) value.Value {
	if typ.Equals(&temporal.StringType) {
		return state.Block.NewCall(
			state.GetPrintf(),
			state.GetFormatString(typ),
			state.Block.NewLoad(types.I8Ptr, state.Block.NewGetElementPtr(
				typ.AsLLType(),
				val,
				I32(0), vectorBodyOffset)))
	}

	return state.Block.NewCall(
		state.GetPrintf(),
		state.GetFormatString(typ),
		val)
}

func (state *State) CompileInlineIndex(typ *temporal.Type, val value.Value) value.Value {
	return state.ReadVectorElement(
		typ.Tuple[0],                // Vector type
		state.TupleGet(typ, val, 0), // Vector in LLVM
		state.TupleGet(typ, val, 1)) // Index in LLVM
}

func (state *State) CompileInlinePanic(_ *temporal.Type, val value.Value) value.Value {
	state.Block.NewCall(state.GetExit(), state.Block.NewTrunc(val, types.I32))
	state.Block.NewUnreachable()
	return nil
}
