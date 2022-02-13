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
	case prism.DyadicOperator:
		return env.CompileDyadicOperator(&t)
	case prism.Function:
		return env.CompileFunction(&t)
	case prism.Alpha:
		return env.CurrentFunction.Params[0]
	case prism.Omega:
		if len(env.CurrentFunction.Params) == 1 {
			return env.CurrentFunction.Params[0]
		} else {
			return env.CurrentFunction.Params[1]
		}
	case prism.Cast:
		return env.CompileCast(t)
	default:
		panic("unreachable")
	}
}

func (env *Environment) CompileFunction(f *prism.Function) value.Value {
	if mfn, ok := env.LLMonadicFunctions[(*f).LLVMise()]; ok {
		return mfn
	} else if dfn, ok := env.LLDyadicFunctions[(*f).LLVMise()]; ok {
		return dfn
	}

	panic("Not found")
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
	case "Max":
		return env.CompileInlineCeil
	case "Min":
		return env.CompileInlineFloor
	default:
		panic("unreachable")
	}
}

func (env *Environment) GetSpecialDCallable(ident *prism.Ident) DCallable {
	switch ident.Name {
	case "GEP":
		return env.CompileInlineIndex
	case "+":
		return env.CompileInlineAdd
	case "-":
		return env.CompileInlineSub
	case "*":
		return env.CompileInlineMul
	case "÷":
		return env.CompileInlineDiv
	case "Max":
		return env.CompileInlineMax
	case "Min":
		return env.CompileInlineMin
	case "&":
		return env.CompileInlineAnd
	case "|":
		return env.CompileInlineAnd
	default:
		panic("unreachable")
	}
}

func (env *Environment) CompileDyadicOperator(dop *prism.DyadicOperator) value.Value {
	switch dop.Operator {
	case prism.KindMapOperator:
		return env.CompileInlineMap(
			dop.Left,
			Value{env.CompileExpression(&dop.Right), dop.Right.Type()})

	case prism.KindFoldlOperator:
		return env.CompileInlineFoldl(
			dop.Left,
			Value{env.CompileExpression(&dop.Right), dop.Right.Type()})
	}
	panic("unreachable")
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
	case "Max":
		return env.CompileInlineCeil(Value{env.CompileExpression(&app.Operand), app.Operand.Type()})
	case "Min":
		return env.CompileInlineFloor(Value{env.CompileExpression(&app.Operand), app.Operand.Type()})
	default:
		return env.Block.NewCall(
			env.LLMonadicFunctions[app.Operator.LLVMise()],
			env.CompileExpression(&app.Operand))
	}
}

func (env *Environment) CompileDApplication(app *prism.DApplication) value.Value {
	switch app.Operator.Ident().Name {
	case "+":
		return env.CompileInlineAdd(
			Value{env.CompileExpression(&app.Left), app.Left.Type()},
			Value{env.CompileExpression(&app.Right), app.Right.Type()})
	case "-":
		return env.CompileInlineSub(
			Value{env.CompileExpression(&app.Left), app.Left.Type()},
			Value{env.CompileExpression(&app.Right), app.Right.Type()})
	case "*":
		return env.CompileInlineMul(
			Value{env.CompileExpression(&app.Left), app.Left.Type()},
			Value{env.CompileExpression(&app.Right), app.Right.Type()})
	case "÷":
		return env.CompileInlineDiv(
			Value{env.CompileExpression(&app.Left), app.Left.Type()},
			Value{env.CompileExpression(&app.Right), app.Right.Type()})
	case "Max":
		return env.CompileInlineMax(
			Value{env.CompileExpression(&app.Left), app.Left.Type()},
			Value{env.CompileExpression(&app.Right), app.Right.Type()})
	case "Min":
		return env.CompileInlineMin(
			Value{env.CompileExpression(&app.Left), app.Left.Type()},
			Value{env.CompileExpression(&app.Right), app.Right.Type()})
	case "GEP":
		return env.CompileInlineIndex(
			Value{env.CompileExpression(&app.Left), app.Left.Type()},
			Value{env.CompileExpression(&app.Right), app.Right.Type()})
	case "Append":
		return env.CompileInlineAppend(
			Value{env.CompileExpression(&app.Left), app.Left.Type()},
			Value{env.CompileExpression(&app.Right), app.Right.Type()})
	case "Equals":
		return env.CompileInlineEqual(
			Value{env.CompileExpression(&app.Left), app.Left.Type()},
			Value{env.CompileExpression(&app.Right), app.Right.Type()})
	case "&":
		return env.CompileInlineAnd(
			Value{env.CompileExpression(&app.Left), app.Left.Type()},
			Value{env.CompileExpression(&app.Right), app.Right.Type()})
	case "|":
		return env.CompileInlineAnd(
			Value{env.CompileExpression(&app.Left), app.Left.Type()},
			Value{env.CompileExpression(&app.Right), app.Right.Type()})
	default:
		call := env.Block.NewCall(
			env.LLDyadicFunctions[app.Operator.LLVMise()],
			env.CompileExpression(&app.Left),
			env.CompileExpression(&app.Right))

		return call
	}
}
