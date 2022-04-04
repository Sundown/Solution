package apotheosis

import (
	"github.com/llir/llvm/ir/value"
	"github.com/sundown/solution/prism"
)

func (env *Environment) compileMonadicApplication(app *prism.MonadicApplication) value.Value {
	if name := app.Operator.Ident().Name; name == "Return" {
		env.Block.NewRet(env.compileExpression(&app.Operand))
		return nil
	} else if fn := env.FetchMonadicCallable(name); fn != nil {
		return env.apply(fn, prism.Value{
			Value: env.compileExpression(&app.Operand),
			Type:  app.Operand.Type()})
	}

	return env.Block.NewCall(
		env.LLMonadicFunctions[app.Operator.LLVMise()],
		env.compileExpression(&app.Operand))

}

func (env *Environment) compileDyadicApplication(app *prism.DyadicApplication) value.Value {
	if fn := env.FetchDyadicCallable(app.Operator.Ident().Name); fn != nil {
		return env.apply(fn, prism.Value{
			Value: env.compileExpression(&app.Left),
			Type:  app.Operator.AlphaType},
			prism.Value{
				Value: env.compileExpression(&app.Right),
				Type:  app.Operator.OmegaType})
	}

	call := env.Block.NewCall(
		env.LLDyadicFunctions[app.Operator.LLVMise()],
		env.compileExpression(&app.Left),
		env.compileExpression(&app.Right))

	return call

}
