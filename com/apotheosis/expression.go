package apotheosis

import (
	"sundown/solution/prism"

	"github.com/llir/llvm/ir/value"
)

func (env *Environment) CompileExpression(expr *prism.Expression) value.Value {
	switch t := (*expr).(type) {
	case prism.MApplication:
		return env.CompileMApplication(&t)
	case prism.DApplication:
		return env.CompileDApplication(&t)
	case prism.Morpheme:
		return env.CompileAtom(&t)
	default:
		panic("unreachable")
	}
}

type MCallable func(val Value) value.Value
type DCallable func(left, right Value) value.Value

func (env *Environment) GetSpecialMCallable(ident *prism.Ident) MCallable {
	switch ident.Name {
	case "Println":
		return env.CompileInlinePrintln
	case "Print":
		return env.CompileInlinePrint
	case "Panic":
		return env.CompileInlinePanic
	case "Len":
		return env.ReadVectorLength
	case "Cap":
		return env.ReadVectorCapacity
	case "Sum":
		return env.CompileInlineSum
	case "Product":
		return env.CompileInlineProduct
	default:
		panic("unreachable")
	}
}

func (env *Environment) GetSpecialDCallable(ident *prism.Ident) DCallable {
	switch ident.Name {
	case "GEP":
		return env.CompileInlineIndex
	default:
		panic("unreachable")
	}
}

func (env *Environment) CompileMApplication(app *prism.MApplication) value.Value {
	switch app.Operator.Ident().Name {
	case "Return":
		env.Block.NewRet(env.CompileExpression(&app.Operand))
		return nil
	case "Println":
		return env.CompileInlinePrintln(Value{env.CompileExpression(&app.Operand), app.Operand.Type()})
	case "Print":
		return env.CompileInlinePrint(Value{env.CompileExpression(&app.Operand), app.Operand.Type()})
	case "Panic":
		return env.CompileInlinePanic(Value{env.CompileExpression(&app.Operand), app.Operand.Type()})
	case "Len":
		return env.ReadVectorLength(Value{env.CompileExpression(&app.Operand), app.Operand.Type()})
	case "Cap":
		return env.ReadVectorCapacity(Value{env.CompileExpression(&app.Operand), app.Operand.Type()})
	case "Sum":
		return env.CompileInlineSum(Value{env.CompileExpression(&app.Operand), app.Operand.Type()})
	case "Product":
		return env.CompileInlineProduct(Value{env.CompileExpression(&app.Operand), app.Operand.Type()})
	default:
		return env.Block.NewCall(
			env.LLMFunctions[app.Operator.LLVMise()],
			env.CompileExpression(&app.Operand))
	}
}

func (env *Environment) CompileDApplication(app *prism.DApplication) value.Value {
	switch app.Operator.Ident().Name {
	case "GEP":
		return env.CompileInlineIndex(
			Value{env.CompileExpression(&app.Left), app.Left.Type()},
			Value{env.CompileExpression(&app.Right), app.Right.Type()})
	case "Map":
		return env.CompileInlineMap(app.Left, app.Right)
	case "Append":
		return env.CompileInlineAppend(
			Value{env.CompileExpression(&app.Left), app.Left.Type()},
			Value{env.CompileExpression(&app.Right), app.Right.Type()})
	case "Equals":
		return env.CompileInlineEqual(
			Value{env.CompileExpression(&app.Left), app.Left.Type()},
			Value{env.CompileExpression(&app.Right), app.Right.Type()})
	default:
		call := env.Block.NewCall(
			env.LLDFunctions[app.Operator.LLVMise()],
			env.CompileExpression(&app.Left),
			env.CompileExpression(&app.Right))

		return call
	}
}
