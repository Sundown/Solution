package apotheosis

import (
	"github.com/sundown/solution/prism"

	"github.com/llir/llvm/ir/value"
)

func (env *Environment) compileExpression(expr *prism.Expression) value.Value {
	switch t := (*expr).(type) {
	case prism.MonadicApplication:
		return env.compileMonadicApplication(&t)
	case prism.DyadicApplication:
		return env.compileDyadicApplication(&t)
	case prism.Morpheme:
		return env.compileAtom(&t)
	case prism.DyadicOperator:
		return env.compileDyadicOperator(&t)
	case prism.Function:
		return env.compileFunction(&t)
	case prism.Alpha:
		return env.CurrentFunction.Params[0]
	case prism.Omega:
		if len(env.CurrentFunction.Params) == 1 {
			return env.CurrentFunction.Params[0]
		} else {
			return env.CurrentFunction.Params[1]
		}
	case prism.Cast:
		return env.compileCast(t)
	default:
		prism.Panic("unreachable")
	}
	panic(nil)
}

func (env *Environment) compileFunction(f *prism.Function) value.Value {
	if mfn, ok := env.LLMonadicFunctions[(*f).LLVMise()]; ok {
		return mfn
	} else if dfn, ok := env.LLDyadicFunctions[(*f).LLVMise()]; ok {
		return dfn
	}

	prism.Panic("Not found")
	panic(nil)
}

type MCallable func(val Value) value.Value
type DCallable func(left, right Value) value.Value

func (env *Environment) GetSpecialMCallable(ident *prism.Ident) MCallable {
	switch ident.Name {
	case "Println":
		return env.compileInlinePrintln
	case "Print":
		return env.compileInlinePrint
	case "Panic":
		return env.compileInlinePanic
	case "Len":
		return env.readVectorLength
	case "Cap":
		return env.readVectorCapacity
	case "Max":
		return env.compileInlineCeil
	case "Min":
		return env.compileInlineFloor
	default:
		prism.Panic("unreachable")
	}
	panic(nil)
}

func (env *Environment) GetSpecialDCallable(ident *prism.Ident) DCallable {
	switch ident.Name {
	case "GEP":
		return env.compileInlineIndex
	case "+":
		return env.compileInlineAdd
	case "-":
		return env.compileInlineSub
	case "*":
		return env.compileInlineMul
	case "÷":
		return env.compileInlineDiv
	case "=":
		return env.compileInlineEqual
	case "Max":
		return env.compileInlineMax
	case "Min":
		return env.compileInlineMin
	case "&":
		return env.compileInlineAnd
	case "|":
		return env.compileInlineAnd
	default:
		prism.Panic("unreachable")
	}
	panic(nil)
}

func (env *Environment) compileDyadicOperator(dop *prism.DyadicOperator) value.Value {
	switch dop.Operator {
	case prism.KindMapOperator:
		return env.compileInlineMap(
			dop.Left,
			Value{env.compileExpression(&dop.Right), dop.Right.Type()})

	case prism.KindReduceOperator:
		return env.compileInlineReduce(
			dop.Left,
			Value{env.compileExpression(&dop.Right), dop.Right.Type()})
	}
	prism.Panic("unreachable")
	panic(nil)
}

func (env *Environment) compileMonadicApplication(app *prism.MonadicApplication) value.Value {
	switch app.Operator.Ident().Name {
	case "Return":
		env.Block.NewRet(env.compileExpression(&app.Operand))
		return nil
	case "Println":
		return env.compileInlinePrintln(Value{env.compileExpression(&app.Operand), app.Operand.Type()})
	case "Print":
		return env.compileInlinePrint(Value{env.compileExpression(&app.Operand), app.Operand.Type()})
	case "Panic":
		return env.compileInlinePanic(Value{env.compileExpression(&app.Operand), app.Operand.Type()})
	case "Len":
		return env.readVectorLength(Value{env.compileExpression(&app.Operand), app.Operand.Type()})
	case "Cap":
		return env.readVectorCapacity(Value{env.compileExpression(&app.Operand), app.Operand.Type()})
	case "Max":
		return env.compileInlineCeil(Value{env.compileExpression(&app.Operand), app.Operand.Type()})
	case "Min":
		return env.compileInlineFloor(Value{env.compileExpression(&app.Operand), app.Operand.Type()})
	default:
		return env.Block.NewCall(
			env.LLMonadicFunctions[app.Operator.LLVMise()],
			env.compileExpression(&app.Operand))
	}
}

func (env *Environment) compileDyadicApplication(app *prism.DyadicApplication) value.Value {
	switch app.Operator.Ident().Name {
	case "+":
		return env.compileInlineAdd(
			Value{env.compileExpression(&app.Left), app.Left.Type()},
			Value{env.compileExpression(&app.Right), app.Right.Type()})
	case "-":
		return env.compileInlineSub(
			Value{env.compileExpression(&app.Left), app.Left.Type()},
			Value{env.compileExpression(&app.Right), app.Right.Type()})
	case "*":
		return env.compileInlineMul(
			Value{env.compileExpression(&app.Left), app.Left.Type()},
			Value{env.compileExpression(&app.Right), app.Right.Type()})
	case "÷":
		return env.compileInlineDiv(
			Value{env.compileExpression(&app.Left), app.Left.Type()},
			Value{env.compileExpression(&app.Right), app.Right.Type()})
	case "=":
		return env.compileInlineEqual(
			Value{env.compileExpression(&app.Left), app.Left.Type()},
			Value{env.compileExpression(&app.Right), app.Right.Type()})
	case "Max":
		return env.compileInlineMax(
			Value{env.compileExpression(&app.Left), app.Left.Type()},
			Value{env.compileExpression(&app.Right), app.Right.Type()})
	case "Min":
		return env.compileInlineMin(
			Value{env.compileExpression(&app.Left), app.Left.Type()},
			Value{env.compileExpression(&app.Right), app.Right.Type()})
	case "GEP":
		return env.compileInlineIndex(
			Value{env.compileExpression(&app.Left), app.Left.Type()},
			Value{env.compileExpression(&app.Right), app.Right.Type()})
	case "Append":
		return env.compileInlineAppend(
			Value{env.compileExpression(&app.Left), app.Left.Type()},
			Value{env.compileExpression(&app.Right), app.Right.Type()})
	case "Equals":
		return env.compileInlineEqual(
			Value{env.compileExpression(&app.Left), app.Left.Type()},
			Value{env.compileExpression(&app.Right), app.Right.Type()})
	case "&":
		return env.compileInlineAnd(
			Value{env.compileExpression(&app.Left), app.Left.Type()},
			Value{env.compileExpression(&app.Right), app.Right.Type()})
	case "|":
		return env.compileInlineAnd(
			Value{env.compileExpression(&app.Left), app.Left.Type()},
			Value{env.compileExpression(&app.Right), app.Right.Type()})
	default:
		call := env.Block.NewCall(
			env.LLDyadicFunctions[app.Operator.LLVMise()],
			env.compileExpression(&app.Left),
			env.compileExpression(&app.Right))

		return call
	}
}
