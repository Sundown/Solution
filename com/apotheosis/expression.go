package apotheosis

import (
	"sundown/solution/prism"

	"github.com/llir/llvm/ir/value"
)

func (state *State) CompileExpression(expr *prism.Expression) value.Value {
	switch t := (*expr).(type) {
	case prism.MApplication:
		return state.CompileMApplication(&t)
	case prism.DApplication:
		return state.CompileDApplication(&t)
	case prism.Morpheme:
		return state.CompileAtom(&t)
	default:
		panic("unreachable")
	}
}

type MCallable func(val Value) value.Value
type DCallable func(left, right Value) value.Value

func (state *State) GetSpecialMCallable(ident *prism.Ident) MCallable {
	switch ident.Name {
	case "Println":
		return state.CompileInlinePrintln
	case "Print":
		return state.CompileInlinePrint
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

func (state *State) GetSpecialDCallable(ident *prism.Ident) DCallable {
	switch ident.Name {
	case "GEP":
		return state.CompileInlineIndex
	default:
		panic("unreachable")
	}
}

func (state *State) CompileApplication(app *prism.Application) value.Value {
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
