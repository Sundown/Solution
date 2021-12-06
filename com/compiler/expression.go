package compiler

import (
	"sundown/solution/temporal"

	"github.com/llir/llvm/ir/value"
)

func (state *State) CompileExpression(expr *temporal.Expression) value.Value {
	if expr.Application != nil {
		return state.CompileApplication(expr.Application)
	} else if expr.Morpheme != nil {
		return state.CompileAtom(expr.Morpheme)
	} else {
		panic("unreachable")
	}
}

type Callable func(*temporal.Type, value.Value) value.Value

func (state *State) GetSpecialCallable(ident *temporal.Ident) Callable {
	switch *ident.Ident {
	case "Println":
		return state.CompileInlinePrintln
	case "Print":
		return state.CompileInlinePrint
	case "GEP":
		return state.CompileInlineIndex
	case "Panic":
		return state.CompileInlinePanic
	case "Len":
		return state.ReadVectorLength
	case "Cap":
		return state.ReadVectorCapacity
	case "Sum":
		return state.CompileInlineSum
	case "Product":
		return state.CompileInlineProduct
	default:
		panic("unreachable")
	}
}

func (state *State) CompileApplication(app *temporal.Application) value.Value {
	switch *app.Function.Ident.Ident {
	case "Return":
		state.Block.NewRet(state.CompileExpression(app.ArgumentAlpha))
		return nil
	case "GEP":
		return state.CompileInlineIndex(app.ArgumentAlpha.TypeOf, state.CompileExpression(app.ArgumentAlpha))
	case "Println":
		return state.CompileInlinePrintln(app.ArgumentAlpha.TypeOf, state.CompileExpression(app.ArgumentAlpha))
	case "Print":
		return state.CompileInlinePrint(app.ArgumentAlpha.TypeOf, state.CompileExpression(app.ArgumentAlpha))
	case "Panic":
		return state.CompileInlinePanic(nil, state.CompileExpression(app.ArgumentAlpha))
	case "Len":
		return state.ReadVectorLength(app.ArgumentAlpha.TypeOf, state.CompileExpression(app.ArgumentAlpha))
	case "Cap":
		return state.ReadVectorCapacity(app.ArgumentAlpha.TypeOf, state.CompileExpression(app.ArgumentAlpha))
	case "Map":
		return state.CompileInlineMap(app.ArgumentAlpha)
	case "Foldl":
		return state.CompileInlineFoldl(app)
	case "Sum":
		return state.CompileInlineSum(app.ArgumentAlpha.TypeOf, state.CompileExpression(app.ArgumentAlpha))
	case "Product":
		return state.CompileInlineProduct(app.ArgumentAlpha.TypeOf, state.CompileExpression(app.ArgumentAlpha))
	case "Append":
		return state.CompileInlineAppend(
			app.ArgumentAlpha.TypeOf, state.CompileExpression(app.ArgumentAlpha),
			app.ArgumentOmega.TypeOf, state.CompileExpression(app.ArgumentOmega))
	case "Equals":
		return state.CompileInlineEqual(
			app.ArgumentAlpha.TypeOf, state.CompileExpression(app.ArgumentAlpha),
			app.ArgumentOmega.TypeOf, state.CompileExpression(app.ArgumentOmega))
	case "First":
		return state.TupleGet(app.ArgumentAlpha.TypeOf, state.CompileExpression(app.ArgumentAlpha), 0)
	case "Second":
		return state.TupleGet(app.ArgumentAlpha.TypeOf, state.CompileExpression(app.ArgumentAlpha), 1)
	case "Third":
		return state.TupleGet(app.ArgumentAlpha.TypeOf, state.CompileExpression(app.ArgumentAlpha), 2)
	default:
		return state.Block.NewCall(
			state.Functions[app.Function.ToLLVMName()],
			state.CompileExpression(app.ArgumentAlpha))
	}
}
