package compiler

import (
	"sundown/solution/parse"

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
	case "Println":
		return state.CompileInlinePrintln(app)
	case "Print":
		return state.CompileInlinePrint(app)
	case "Len":
		if app.Argument.TypeOf.Vector == nil {
			panic("Can't take Len of non-vector")
		}

		return state.Block.NewLoad(types.I64,
			state.Block.NewGetElementPtr(
				app.Argument.TypeOf.AsLLType(),
				state.CompileExpression(app.Argument),
				I32(0), I32(0)))
	case "Cap":
		if app.Argument.TypeOf.Vector == nil {
			panic("Can't take Cap of non-vector")
		}

		return state.Block.NewLoad(types.I64,
			state.Block.NewGetElementPtr(
				app.Argument.TypeOf.AsLLType(),
				state.CompileExpression(app.Argument),
				I32(0), I32(1)))
	case "Map":
		return state.CompileInlineMap(app)
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
	leng := state.Block.NewLoad(types.I64, state.Block.NewGetElementPtr(
		src.Type().(*types.PointerType).ElemType,
		src,
		I32(0), I32(0)))

	state.LLVMPanic(bfalse, "Panic: index %d out of bounds [%d]\n", index, leng)

	bend := state.CurrentFunction.NewBlock("")
	btrue.NewBr(bend)
	bfalse.NewUnreachable()

	state.Block.NewCondBr(
		state.Block.NewICmp(
			enum.IPredSLE,
			leng,
			index),
		bfalse, btrue)

	state.Block = bend
}

// Supply the block in which to generate message and exit call, a printf formatter, and variadic params
func (state *State) LLVMPanic(block *ir.Block, format string, args ...value.Value) {
	var fmt value.Value = block.NewGetElementPtr(
		types.NewArray(uint64(len(format)+1), types.I8),
		state.Module.NewGlobalDef("", constant.NewCharArrayFromString(format+"\x00")), I32(0), I32(0))
	block.NewCall(state.GetPrintf(), append([]value.Value{fmt}, args...)...)
	block.NewCall(state.GetExit(), I32(1))
}
