package compiler

import (
	"sundown/sunday/parse"

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
			expr.Atom.Tuple[0].TypeOf.Vector == nil ||
			expr.Atom.Tuple[1].TypeOf.Atomic == nil ||
			*expr.Atom.Tuple[1].TypeOf.Atomic != "Int" {
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
