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
		return state.CompileInlineIndex(app)
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
		} else {
			format = state.Module.NewGlobalDef("", constant.NewCharArrayFromString("\n\x00"))
			return state.Block.NewCall(state.GetPrintf(),
				state.Block.NewGetElementPtr(types.NewArray(2, types.I8), format, I32(0), I32(0)))
		}
	case "Len":
		if app.Argument.TypeOf.Vector == nil {
			panic("Can't take Len of non-vector")
		}

		vec := state.CompileExpression(app.Argument)

		return state.Block.NewLoad(types.I64, state.Block.NewGetElementPtr(app.Argument.TypeOf.AsLLType(), vec, I32(0), I32(0)))
	case "Sum":
		return state.CompileInlineSum(app)
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
