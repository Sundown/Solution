package apotheosis

import (
	"sundown/solution/prism"

	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

func (state *State) CompileInlinePrintln(val Value) value.Value {
	if prism.EqType(val.Type, &prism.StringType) {
		return state.Block.NewCall(
			state.GetPrintf(),
			state.GetFormatStringln(&val.Type),
			state.Block.NewLoad(types.I8Ptr, state.Block.NewGetElementPtr(
				val.Type.Realise(),
				val.Value,
				I32(0), vectorBodyOffset)))
	}

	return state.Block.NewCall(
		state.GetPrintf(),
		state.GetFormatStringln(&val.Type),
		val.Value)
}

func (state *State) CompileInlinePrint(val Value) value.Value {
	if prism.EqType(val.Type, &prism.StringType) {
		return state.Block.NewCall(
			state.GetPrintf(),
			state.GetFormatString(&val.Type),
			state.Block.NewLoad(types.I8Ptr, state.Block.NewGetElementPtr(
				val.Type.Realise(),
				val.Value,
				I32(0), vectorBodyOffset)))
	}

	return state.Block.NewCall(
		state.GetPrintf(),
		state.GetFormatStringln(&val.Type),
		val.Value)
}

func (state *State) CompileInlineIndex(left, right Value) value.Value {
	return state.ReadVectorElement(right, left.Value)
}

func (state *State) CompileInlinePanic(val Value) value.Value {
	state.Block.NewCall(state.GetExit(), state.Block.NewTrunc(val.Value, types.I32))
	state.Block.NewUnreachable()
	return nil
}
