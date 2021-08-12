package compiler

import (
	"sundown/solution/parse"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

func (state *State) CompileInlinePrintln(app *parse.Application) value.Value {
	header := state.CompileExpression(app.Argument)

	var format *ir.Global

	if header.Type().Equal(types.NewPointer(parse.StringType.AsLLType())) {
		format = state.Module.NewGlobalDef("", constant.NewCharArrayFromString("%s\n\x00"))

		return state.Block.NewCall(state.GetPrintf(),
			state.Block.NewGetElementPtr(types.NewArray(4, types.I8), format, I32(0), I32(0)),
			state.Block.NewLoad(types.I8Ptr, state.Block.NewGetElementPtr(
				header.Type().(*types.PointerType).ElemType,
				header,
				I32(0), I32(2))))
	} else if header.Type().Equal(types.I64) {
		format = state.Module.NewGlobalDef("", constant.NewCharArrayFromString("%d\n\x00"))
		return state.Block.NewCall(state.GetPrintf(),
			state.Block.NewGetElementPtr(types.NewArray(4, types.I8), format, I32(0), I32(0)),
			header)
	} else if header.Type().Equal(types.Double) {
		format = state.Module.NewGlobalDef("", constant.NewCharArrayFromString("%f\n\x00"))
		return state.Block.NewCall(state.GetPrintf(),
			state.Block.NewGetElementPtr(types.NewArray(4, types.I8), format, I32(0), I32(0)),
			header)
	} else if header.Type().Equal(types.I8) {
		format = state.Module.NewGlobalDef("", constant.NewCharArrayFromString("%c\n\x00"))
		return state.Block.NewCall(state.GetPrintf(),
			state.Block.NewGetElementPtr(types.NewArray(4, types.I8), format, I32(0), I32(0)),
			header)
	} else {
		format = state.Module.NewGlobalDef("", constant.NewCharArrayFromString("\n\x00"))
		return state.Block.NewCall(state.GetPrintf(),
			state.Block.NewGetElementPtr(types.NewArray(2, types.I8), format, I32(0), I32(0)))
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

func (state *State) CompileInlineSum(app *parse.Application) value.Value {
	if app.Argument.TypeOf.Vector == nil {
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
			llvec.Type().(*types.PointerType).ElemType, llvec, I32(0), I32(0)))

		vec_body := state.Block.NewLoad(types.NewPointer(types.Double), state.Block.NewGetElementPtr(
			llvec.Type().(*types.PointerType).ElemType, llvec, I32(0), I32(2)))
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
	}

	return state.Block.NewLoad(types.I64, accum)
}

func (state *State) CompileInlineMap(app *parse.Application) value.Value {
	if app.Argument.TypeOf.Tuple == nil {
		panic("Map requires Tuple")
	}

	// The vector in AST
	vec := app.Argument.Atom.Tuple[1]
	// The vector in LLVM
	llvec := state.CompileExpression(vec)

	head_type := vec.TypeOf.AsLLType()
	elm_type := vec.Atom.Vector[0].TypeOf.AsLLType()

	should_store := true
	if app.Argument.Atom.Tuple[0].Atom.Function.Gives.Equals(parse.AtomicType("Void")) {
		should_store = false
	}

	// Map is 1:1 so leng and cap are just copied from input vector
	leng := state.Block.NewGetElementPtr(head_type, llvec, I32(0), I32(0))
	cap := state.Block.NewGetElementPtr(head_type, llvec, I32(0), I32(1))

	head := state.BuildVectorHeader(head_type)

	// Copy length
	state.Block.NewStore(
		state.Block.NewLoad(types.I64, leng),
		state.Block.NewGetElementPtr(head_type, llvec, I32(0), I32(0)))

	// Copy capacity
	state.Block.NewStore(
		state.Block.NewLoad(types.I64, cap),
		state.Block.NewGetElementPtr(head_type, llvec, I32(0), I32(1)))

	var body *ir.InstBitCast
	if should_store {
		// Allocate a body of capacity * element width, and cast to element type
		body = state.Block.NewBitCast(
			state.Block.NewCall(state.GetCalloc(),
				I32(vec.Atom.Vector[0].TypeOf.WidthInBytes()),                         // Byte size of elements
				state.Block.NewTrunc(state.Block.NewLoad(types.I64, cap), types.I32)), // How much memory to alloc
			types.NewPointer(elm_type)) // Cast alloc'd memory to typ
	}

	// -------------
	// # LOOP BODY #
	// -------------
	if app.Argument.Atom.Tuple[1].TypeOf.Vector != nil {
		vec_body := state.Block.NewLoad(
			types.NewPointer(elm_type),
			state.Block.NewGetElementPtr(head_type, llvec, I32(0), I32(2)))

		counter := state.Block.NewAlloca(types.I64)
		state.Block.NewStore(I64(0), counter)

		// Body
		// Get elem, add to accum, increment counter, conditional jump to body
		loopblock := state.CurrentFunction.NewBlock("")
		state.Block.NewBr(loopblock)
		// Add to accum
		cur_counter := loopblock.NewLoad(types.I64, counter)

		var cur_elm value.Value
		cur_elm = loopblock.NewGetElementPtr(elm_type, vec_body, cur_counter)

		if vec.Atom.Vector[0].TypeOf.Atomic != nil {
			cur_elm = loopblock.NewLoad(elm_type, cur_elm)
		}

		call := loopblock.NewCall(
			state.CompileExpression(app.Argument.Atom.Tuple[0]),
			cur_elm)

		if should_store {
			loopblock.NewStore(
				call,
				loopblock.NewGetElementPtr(elm_type, body, cur_counter))
		}
		// Increment counter
		loopblock.NewStore(loopblock.NewAdd(loopblock.NewLoad(types.I64, counter), I64(1)), counter)
		cond := loopblock.NewICmp(enum.IPredSLT, cur_counter, loopblock.NewLoad(types.I64, leng))
		exitblock := state.CurrentFunction.NewBlock("")
		loopblock.NewCondBr(cond, loopblock, exitblock)
		state.Block = exitblock

		if should_store {
			state.Block.NewStore(body, state.Block.NewGetElementPtr(vec.TypeOf.AsLLType(), head, I32(0), I32(2)))
			return head
		} else {
			return nil
		}

	} else {
		panic("Map needs (F, [T])")
	}
}
