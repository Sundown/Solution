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

func (state *State) CompileMApplication(app *prism.MApplication) value.Value {
	switch app.Operator.Ident().Name {
	case "Return":
		state.Block.NewRet(state.CompileExpression(&app.Operand))
		return nil
	case "Println":
		return state.CompileInlinePrintln(Value{state.CompileExpression(&app.Operand), app.Operand.Type()})
	case "Print":
		return state.CompileInlinePrint(Value{state.CompileExpression(&app.Operand), app.Operand.Type()})
	case "Panic":
		return state.CompileInlinePanic(Value{state.CompileExpression(&app.Operand), app.Operand.Type()})
	case "Len":
		return state.ReadVectorLength(Value{state.CompileExpression(&app.Operand), app.Operand.Type()})
	case "Cap":
		return state.ReadVectorCapacity(Value{state.CompileExpression(&app.Operand), app.Operand.Type()})
	case "Sum":
		return state.CompileInlineSum(Value{state.CompileExpression(&app.Operand), app.Operand.Type()})
	case "Product":
		return state.CompileInlineProduct(Value{state.CompileExpression(&app.Operand), app.Operand.Type()})
	default:
		return state.Block.NewCall(
			state.MFunctions[app.Operator.LLVMise()],
			state.CompileExpression(&app.Operand))
	}
}

func (state *State) CompileDApplication(app *prism.DApplication) value.Value {
	switch app.Operator.Ident().Name {
	case "GEP":
		return state.CompileInlineIndex(
			Value{state.CompileExpression(&app.Left), app.Left.Type()},
			Value{state.CompileExpression(&app.Right), app.Right.Type()})
	case "Map":
		return state.CompileInlineMap(app.Left, app.Right)
	case "Append":
		return state.CompileInlineAppend(
			Value{state.CompileExpression(&app.Left), app.Left.Type()},
			Value{state.CompileExpression(&app.Right), app.Right.Type()})
	case "Equals":
		return state.CompileInlineEqual(
			Value{state.CompileExpression(&app.Left), app.Left.Type()},
			Value{state.CompileExpression(&app.Right), app.Right.Type()})
	default:
		call := state.Block.NewCall(
			state.DFunctions[app.Operator.LLVMise()],
			state.CompileExpression(&app.Left),
			state.CompileExpression(&app.Right))

		return call
	}
}
