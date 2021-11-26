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

func (state *State) CompileInlineIndex(app *temporal.Application) value.Value {
	if app.Argument.Atom == nil || app.Argument.Atom.Tuple == nil ||
		app.Argument.Atom.Tuple[0].TypeOf.Vector == nil {
		panic("Index requires tuple: ([T], Int | Nat)")
	}

	index := state.CompileExpression(app.Argument.Atom.Tuple[1])
	src := state.CompileExpression(app.Argument.Atom.Tuple[0])

	state.ValidateVectorIndex(src, index)

	head_typ := app.Argument.Atom.Tuple[0].TypeOf
	elem_typ := head_typ.Vector.AsLLType()

	element := state.Block.NewGetElementPtr(
		elem_typ, state.Block.NewLoad(
			types.NewPointer(elem_typ),
			state.Block.NewGetElementPtr(
				head_typ.AsLLType(),
				src,
				I32(0), I32(2))),
		index)

	if app.Argument.Atom.Tuple[0].TypeOf.Vector.Atomic != nil {
		return state.Block.NewLoad(elem_typ, element)
	}

	return element
}
