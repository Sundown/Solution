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

		if !header.Type().Equal(types.NewPointer(parse.StringType.AsLLType())) {
			panic("Print requires String")
		}

		format := state.Module.NewGlobalDef("", constant.NewCharArrayFromString("%s\x00"))

		return state.Block.NewCall(state.GetPrintf(),
			state.Block.NewGetElementPtr(types.NewArray(3, types.I8), format, I32(0), I32(0)),
			state.Block.NewLoad(types.I8Ptr, state.Block.NewGetElementPtr(
				header.Type().(*types.PointerType).ElemType,
				header,
				I32(0), I32(2))))
	case "Sum":
		if app.Argument.TypeOf.Vector == nil {
			panic("Sum requires Vector")
		}

		accum := state.Block.NewAlloca(types.I64)
		state.Block.NewStore(I64(0), accum)
		if *app.Argument.TypeOf.Vector.Atomic == "Int" {
			vec := app.Argument
			llvec := state.CompileExpression(vec)
			pbody := state.Block.NewGetElementPtr(
				llvec.Type().(*types.PointerType).ElemType,
				llvec, I32(0), I32(2))

			//body := state.Block.NewLoad(types.I64Ptr, pbody)

			leng := state.Block.NewLoad(
				types.I64,
				state.Block.NewGetElementPtr(
					llvec.Type().(*types.PointerType).ElemType,
					llvec,
					I32(0), I32(0)))

			loopblock := state.CurrentFunction.NewBlock("")
			state.Block.NewBr(loopblock)
			first := loopblock.NewPhi(ir.NewIncoming(I64(0), state.Block))
			first.Incs = append(first.Incs, ir.NewIncoming(loopblock.NewAdd(first, I64(1)), loopblock))
			leaveblock := state.CurrentFunction.NewBlock("")

			a := loopblock.NewLoad(types.I64, loopblock.NewLoad(types.I64Ptr, loopblock.NewGetElementPtr(types.I64Ptr, pbody, first)))

			loopblock.NewStore(loopblock.NewAdd(a, loopblock.NewLoad(types.I64, accum)), accum)

			cond := loopblock.NewICmp(enum.IPredEQ, leng, first)

			loopblock.NewCondBr(cond, leaveblock, loopblock)

			state.Block = leaveblock
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
