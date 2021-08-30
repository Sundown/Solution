package compiler

import (
	"sundown/solution/parse"

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
	case "Panic":
		if app.Argument.TypeOf.Equals(parse.AtomicType("Int")) {
			state.Block.NewCall(state.GetExit(), state.Block.NewTrunc(state.CompileExpression(app.Argument), types.I32))
			state.Block.NewUnreachable()
			return nil
		} else {
			panic("Cannot call Panic with non-int")
		}
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
