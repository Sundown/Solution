package compiler

import (
	"sundown/solution/temporal"

	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

func (state *State) CompileExpression(expr *temporal.Expression) value.Value {
	if expr.Application != nil {
		return state.CompileApplication(expr.Application)
	} else if expr.Atom != nil {
		return state.CompileAtom(expr.Atom)
	} else {
		panic("unreachable")
	}
}

type Callable func(*temporal.Type, value.Value) value.Value

func (state *State) GetSpecialCallable(ident *temporal.Ident) Callable {
	switch *ident.Ident {
	case "Println":
		return state.CompileInlinePrintln
	default:
		panic("unreachable")
	}
}

func (state *State) CompileInlinePanic(_ *temporal.Type, val value.Value) value.Value {
	state.Block.NewCall(state.GetExit(), state.Block.NewTrunc(val, types.I32))
	state.Block.NewUnreachable()
	return nil
}

func (state *State) CompileInlineLen(typ *temporal.Type, val value.Value) value.Value {
	return state.Block.NewLoad(types.I64,
		state.Block.NewGetElementPtr(typ.AsLLType(), val, I32(0), vectorLenOffset))
}

func (state *State) CompileInlineCap(typ *temporal.Type, val value.Value) value.Value {
	return state.Block.NewLoad(types.I64,
		state.Block.NewGetElementPtr(typ.AsLLType(), val, I32(0), vectorCapOffset))
}

func (state *State) CompileApplication(app *temporal.Application) value.Value {
	switch *app.Function.Ident.Ident {
	case "Return":
		state.Block.NewRet(state.CompileExpression(app.Argument))
		return nil
	case "GEP":
		return state.CompileInlineIndex(app.Argument.TypeOf, state.CompileExpression(app.Argument))
	case "Println":
		return state.CompileInlinePrintln(app.Argument.TypeOf, state.CompileExpression(app.Argument))
	case "Panic":
		return state.CompileInlinePanic(nil, state.CompileExpression(app.Argument))
	case "Len":
		return state.CompileInlineLen(app.Argument.TypeOf, state.CompileExpression(app.Argument))
	case "Cap":
		return state.CompileInlineCap(app.Argument.TypeOf, state.CompileExpression(app.Argument))
	case "Map":
		return state.CompileInlineMap(app.Argument)
	case "Foldl":
		return state.CompileInlineFoldl(app)
	case "Sum":
		return state.CompileInlineSum(app)
	case "Product":
		return state.CompileInlineProduct(app)
	case "Append":
		return state.CompileInlineAppend(app.Argument)
	case "Equals":
		return state.CompileInlineEqual(app.Argument)
	case "First":
		return state.TupleGet(app.Argument.TypeOf, state.CompileExpression(app.Argument), 0)
	case "Second":
		return state.TupleGet(app.Argument.TypeOf, state.CompileExpression(app.Argument), 1)
	case "Third":
		return state.TupleGet(app.Argument.TypeOf, state.CompileExpression(app.Argument), 2)
	default:
		return state.Block.NewCall(
			state.Functions[app.Function.ToLLVMName()],
			state.CompileExpression(app.Argument))
	}
}
