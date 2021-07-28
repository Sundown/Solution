package compiler

import (
	"sundown/sunday/parse"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

func (state *State) CompileExpression(expr *parse.Expression) value.Value {
	if expr.Atom != nil {
		return state.CompileAtom(expr.Atom)
	} else if expr.Application != nil {
		return state.CompileApplication(expr.Application)
	} else {
		panic("unreachable")
	}
}

func (state *State) CompileApplication(app *parse.Application) value.Value {
	switch *app.Function.Ident.Ident {
	case "Return":
		state.Block.NewRet(state.CompileExpression(app.Argument))
		return nil
	case "GEP":
		expr := app.Argument
		if expr.Atom == nil || expr.Atom.Tuple == nil ||
			expr.Atom.Tuple[0].TypeOf.Vector == nil {
			panic("Index requires tuple: ([T], Int | Nat)")
		}

		index := state.CompileExpression(expr.Atom.Tuple[1])
		src := state.CompileExpression(expr.Atom.Tuple[0])

		state.ValidateVectorIndex(src, index)

		elem_typ := expr.Atom.Tuple[0].TypeOf.Vector.AsLLType()

		element := state.Block.NewGetElementPtr(
			elem_typ, state.Block.NewLoad(
				types.NewPointer(elem_typ),
				state.Block.NewGetElementPtr(
					src.Type().(*types.PointerType).ElemType,
					src,
					I32(0), I32(2))),
			index)

		if expr.Atom.Tuple[0].TypeOf.Vector.Atomic != nil {
			return state.Block.NewLoad(element.Type().(*types.PointerType).ElemType, element)
		}

		return element
	case "Print":
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
		} else {
			format = state.Module.NewGlobalDef("", constant.NewCharArrayFromString("%d\x00"))
			return state.Block.NewCall(state.GetPrintf(),
				state.Block.NewGetElementPtr(types.NewArray(3, types.I8), format, I32(0), I32(0)),
				header)
		}
	case "Len":
		if app.Argument.TypeOf.Vector == nil {
			panic("Can't take Len of non-vector")
		}

		vec := state.CompileExpression(app.Argument)

		return state.Block.NewLoad(types.I64, state.Block.NewGetElementPtr(app.Argument.TypeOf.AsLLType(), vec, I32(0), I32(0)))
	case "Sum":
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
			vec_len.SetName("len")
			vec_body := state.Block.NewLoad(types.I64Ptr, state.Block.NewGetElementPtr(
				llvec.Type().(*types.PointerType).ElemType, llvec, I32(0), I32(2)))
			counter := state.Block.NewAlloca(types.I64)
			state.Block.NewStore(I64(0), counter)
			counter.SetName("counter_ptr")
			accum := state.Block.NewAlloca(types.I64)
			state.Block.NewStore(I64(0), accum)
			accum.SetName("accum")
			// Body
			// Get elem, add to accum, increment counter, conditional jump to body
			loopblock := state.CurrentFunction.NewBlock("loop_body")
			state.Block.NewBr(loopblock)
			// Add to accum
			cur_counter := loopblock.NewLoad(types.I64, counter)
			cur_counter.SetName("counter")
			cur_elm := loopblock.NewLoad(types.I64, loopblock.NewGetElementPtr(types.I64, vec_body, cur_counter))
			cur_elm.SetName("elm")
			loopblock.NewStore(loopblock.NewAdd(loopblock.NewLoad(types.I64, accum), cur_elm), accum)
			// Increment counter
			loopblock.NewStore(loopblock.NewAdd(loopblock.NewLoad(types.I64, counter), I64(1)), counter)
			cond := loopblock.NewICmp(enum.IPredSLT, cur_counter, vec_len)
			exitblock := state.CurrentFunction.NewBlock("exit_loop")
			loopblock.NewCondBr(cond, loopblock, exitblock)
			state.Block = exitblock
			return state.Block.NewLoad(types.I64, accum)
		} else if *app.Argument.TypeOf.Vector.Atomic == "Real" {

		}

		return state.Block.NewLoad(types.I64, accum)
	default:
		return state.Block.NewCall(
			state.Functions[app.Function.ToLLVMName()],
			state.CompileExpression(app.Argument))
	}
}

func (state *State) ValidateVectorIndex(src value.Value, index value.Value) {
	btrue := state.CurrentFunction.NewBlock("")
	bfalse := state.CurrentFunction.NewBlock("")
	bfalse.NewCall(state.GetExit(), I32(10))
	bend := state.CurrentFunction.NewBlock("")
	btrue.NewBr(bend)
	bfalse.NewUnreachable()

	state.Block.NewCondBr(
		state.Block.NewICmp(
			enum.IPredSLE,
			state.Block.NewLoad(types.I64, state.Block.NewGetElementPtr(
				src.Type().(*types.PointerType).ElemType,
				src,
				I32(0), I32(0))),
			index),
		bfalse, btrue)

	state.Block = bend
}
