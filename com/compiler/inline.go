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
	src := state.TupleGet(typ, val, 0)
	index := state.TupleGet(typ, val, 1)

	state.ValidateVectorIndex(typ.Tuple[0], src, index)

	head_typ := typ.Tuple[0]
	elem_typ := typ.Tuple[0].Vector.AsLLType()

	element := state.Block.NewGetElementPtr(
		elem_typ, state.Block.NewLoad(
			types.NewPointer(elem_typ),
			state.Block.NewGetElementPtr(
				head_typ.AsLLType(),
				src,
				I32(0), I32(2))),
		index)

	if typ.Tuple[0].Vector.Atomic != nil {
		return state.Block.NewLoad(elem_typ, element)
	}

	return element
}
