package compiler

import (
	"fmt"
	"sundown/solution/parse"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

func (state *State) CompileInlinePrintln(app *parse.Application) value.Value {
	header := state.CompileExpression(app.Argument)
	typ := app.Argument.TypeOf
	if typ.Equals(&parse.StringType) {
		return state.Block.NewCall(state.GetPrintf(), state.Block.NewGetElementPtr(types.NewArray(4, types.I8), state.Module.NewGlobalDef("", constant.NewCharArrayFromString("%s\n\x00")), I32(0), I32(0)), state.Block.NewLoad(types.I8Ptr, state.Block.NewGetElementPtr(header.Type().(*types.PointerType).ElemType, header, I32(0), I32(2))))
	} else if typ.Equals(&parse.IntType) {
		return state.Block.NewCall(state.GetPrintf(),
			state.Block.NewGetElementPtr(
				types.NewArray(4, types.I8),
				state.Module.NewGlobalDef("", constant.NewCharArrayFromString("%d\n\x00")),
				I32(0), I32(0)),
			header)
	} else if typ.Equals(&parse.RealType) {
		return state.Block.NewCall(state.GetPrintf(),
			state.Block.NewGetElementPtr(
				types.NewArray(4, types.I8),
				state.Module.NewGlobalDef("", constant.NewCharArrayFromString("%f\n\x00")),
				I32(0), I32(0)),
			header)
	} else if typ.Equals(&parse.CharType) {
		return state.Block.NewCall(
			state.GetPrintf(),
			state.Block.NewGetElementPtr(
				types.NewArray(4, types.I8),
				state.Module.NewGlobalDef("", constant.NewCharArrayFromString("%c\n\x00")),
				I32(0), I32(0)),
			header)
	} else {
		return state.Block.NewCall(state.GetPrintf(),
			state.Block.NewGetElementPtr(
				types.NewArray(2, types.I8),
				state.Module.NewGlobalDef("", constant.NewCharArrayFromString("Can't print that type yet!\n\x00")),
				I32(0), I32(0)))
	}
}

func (state *State) CompileInlinePrint(app *parse.Application) value.Value {
	header := state.CompileExpression(app.Argument)

	var format *ir.Global

	if header.Type().Equal(types.NewPointer(parse.StringType.AsLLType())) {
		format = state.Module.NewGlobalDef("", constant.NewCharArrayFromString("%s\x00"))

		return state.Block.NewCall(state.GetPrintf(),
			state.Block.NewGetElementPtr(types.NewArray(3, types.I8), format, I32(0), I32(0)),
			state.Block.NewLoad(types.I8Ptr, state.Block.NewGetElementPtr(
				header.Type().(*types.PointerType).ElemType,
				header,
				I32(0), I32(2))))
	} else if header.Type().Equal(types.I64) {
		format = state.Module.NewGlobalDef("", constant.NewCharArrayFromString("%d\x00"))
		return state.Block.NewCall(state.GetPrintf(),
			state.Block.NewGetElementPtr(types.NewArray(3, types.I8), format, I32(0), I32(0)),
			header)
	} else if header.Type().Equal(types.Double) {
		format = state.Module.NewGlobalDef("", constant.NewCharArrayFromString("%f\x00"))
		return state.Block.NewCall(state.GetPrintf(),
			state.Block.NewGetElementPtr(types.NewArray(3, types.I8), format, I32(0), I32(0)),
			header)
	} else if header.Type().Equal(types.I8) {
		format = state.Module.NewGlobalDef("", constant.NewCharArrayFromString("%c\x00"))
		return state.Block.NewCall(state.GetPrintf(),
			state.Block.NewGetElementPtr(types.NewArray(3, types.I8), format, I32(0), I32(0)),
			header)
	} else {
		format = state.Module.NewGlobalDef("", constant.NewCharArrayFromString("\n\x00"))
		return state.Block.NewCall(state.GetPrintf(),
			state.Block.NewGetElementPtr(types.NewArray(2, types.I8), format, I32(0), I32(0)))
	}
}

func (state *State) CompileInlineIndex(app *parse.Application) value.Value {
	if app.Argument.Atom == nil || app.Argument.Atom.Tuple == nil ||
		app.Argument.Atom.Tuple[0].TypeOf.Vector == nil {
		panic("Index requires tuple: ([T], Int | Nat)")
	}

	index := state.CompileExpression(app.Argument.Atom.Tuple[1])
	src := state.CompileExpression(app.Argument.Atom.Tuple[0])

	//state.ValidateVectorIndex(src, index)

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

func (state *State) CompileInlineSum(app *parse.Application) value.Value {
	if app.Argument.TypeOf.Vector == nil {
		fmt.Println(app.Argument.TypeOf.String())
		panic("Sum requires Vector")
	}

	accum := state.Block.NewAlloca(types.I64)
	state.Block.NewStore(I64(0), accum)
	if *app.Argument.TypeOf.Vector.Atomic == "Int" {
		vec := app.Argument
		llvec := state.CompileExpression(vec)
		vec_len := state.Block.NewLoad(types.I64, state.Block.NewGetElementPtr(
			llvec.Type().(*types.PointerType).ElemType, llvec, I32(0), I32(0)))

		vec_body := state.Block.NewLoad(types.I64Ptr, state.Block.NewGetElementPtr(
			llvec.Type().(*types.PointerType).ElemType, llvec, I32(0), I32(2)))
		counter := state.Block.NewAlloca(types.I64)
		state.Block.NewStore(I64(0), counter)

		accum := state.Block.NewAlloca(types.I64)
		state.Block.NewStore(I64(0), accum)

		// Body
		// Get elem, add to accum, increment counter, conditional jump to body
		loopblock := state.CurrentFunction.NewBlock("")
		state.Block.NewBr(loopblock)
		// Add to accum
		cur_counter := loopblock.NewLoad(types.I64, counter)

		cur_elm := loopblock.NewLoad(types.I64, loopblock.NewGetElementPtr(types.I64, vec_body, cur_counter))

		loopblock.NewStore(loopblock.NewAdd(loopblock.NewLoad(types.I64, accum), cur_elm), accum)
		// Increment counter
		loopblock.NewStore(loopblock.NewAdd(loopblock.NewLoad(types.I64, counter), I64(1)), counter)
		cond := loopblock.NewICmp(enum.IPredSLT, cur_counter, vec_len)
		exitblock := state.CurrentFunction.NewBlock("")
		loopblock.NewCondBr(cond, loopblock, exitblock)
		state.Block = exitblock
		return state.Block.NewLoad(types.I64, accum)
	} else if *app.Argument.TypeOf.Vector.Atomic == "Real" {
		vec := app.Argument
		llvec := state.CompileExpression(vec)
		vec_len := state.Block.NewLoad(types.I64, state.Block.NewGetElementPtr(
			llvec.Type().(*types.PointerType).ElemType, llvec, I32(0), vectorLenOffset))

		vec_body := state.Block.NewLoad(types.NewPointer(types.Double), state.Block.NewGetElementPtr(
			llvec.Type().(*types.PointerType).ElemType, llvec, I32(0), vectorBodyOffset))
		counter := state.Block.NewAlloca(types.I64)
		state.Block.NewStore(I64(0), counter)

		accum := state.Block.NewAlloca(types.Double)
		state.Block.NewStore(constant.NewFloat(types.Double, 0), accum)

		// Body
		// Get elem, add to accum, increment counter, conditional jump to body
		loopblock := state.CurrentFunction.NewBlock("")
		state.Block.NewBr(loopblock)
		// Add to accum
		cur_counter := loopblock.NewLoad(types.I64, counter)

		cur_elm := loopblock.NewLoad(types.Double, loopblock.NewGetElementPtr(types.Double, vec_body, cur_counter))

		loopblock.NewStore(loopblock.NewFAdd(loopblock.NewLoad(types.Double, accum), cur_elm), accum)
		// Increment counter
		loopblock.NewStore(loopblock.NewAdd(loopblock.NewLoad(types.I64, counter), I64(1)), counter)
		cond := loopblock.NewICmp(enum.IPredSLT, cur_counter, vec_len)
		exitblock := state.CurrentFunction.NewBlock("")
		loopblock.NewCondBr(cond, loopblock, exitblock)
		state.Block = exitblock
		return state.Block.NewLoad(types.Double, accum)
	} else if *app.Argument.TypeOf.Vector.Atomic == "Char" {
		vec := app.Argument
		llvec := state.CompileExpression(vec)
		vec_len := state.Block.NewLoad(types.I64, state.Block.NewGetElementPtr(
			llvec.Type().(*types.PointerType).ElemType, llvec, I32(0), vectorLenOffset))

		vec_body := state.Block.NewLoad(types.NewPointer(types.I8), state.Block.NewGetElementPtr(
			llvec.Type().(*types.PointerType).ElemType, llvec, I32(0), vectorBodyOffset))

		counter := state.Block.NewAlloca(types.I64)
		state.Block.NewStore(I64(0), counter)

		accum := state.Block.NewAlloca(types.I8)
		state.Block.NewStore(constant.NewInt(types.I8, 0), accum)

		// Body
		// Get elem, add to accum, increment counter, conditional jump to body
		loopblock := state.CurrentFunction.NewBlock("")
		state.Block.NewBr(loopblock)
		// Add to accum
		cur_counter := loopblock.NewLoad(types.I64, counter)

		cur_elm := loopblock.NewLoad(types.I8, loopblock.NewGetElementPtr(types.I8, vec_body, cur_counter))

		loopblock.NewStore(loopblock.NewAdd(loopblock.NewLoad(types.I8, accum), cur_elm), accum)
		// Increment counter
		loopblock.NewStore(loopblock.NewAdd(loopblock.NewLoad(types.I64, counter), I64(1)), counter)
		cond := loopblock.NewICmp(enum.IPredSLT, cur_counter, vec_len)
		exitblock := state.CurrentFunction.NewBlock("")
		loopblock.NewCondBr(cond, loopblock, exitblock)
		state.Block = exitblock
		return state.Block.NewLoad(types.I8, accum)
	}

	return state.Block.NewLoad(types.I64, accum)
}
