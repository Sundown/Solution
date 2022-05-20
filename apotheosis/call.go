package apotheosis

import (
	"github.com/llir/llvm/ir/value"
	"github.com/sundown/solution/prism"
)

func (env *Environment) newMonadicApplication(app *prism.MonadicApplication) value.Value {
	if name := app.Operator.Ident().Name; name == "←" {
		env.Block.NewRet(env.newExpression(&app.Operand))
		return nil
	} else if fn := env.FetchMonadicCallable(name); fn != nil {
		return env.apply(fn, prism.Value{
			Value: env.newExpression(&app.Operand),
			Type:  app.Operand.Type()})
	}

	return env.Block.NewCall(
		env.LLMonadicFunctions[app.Operator.LLVMise()],
		env.newExpression(&app.Operand))

}

func (env *Environment) newDyadicApplication(app *prism.DyadicApplication) value.Value {
	if fn := env.FetchDyadicCallable(app.Operator.Ident().Name); fn != nil {
		return env.apply(fn, prism.Value{
			Value: env.newExpression(&app.Left),
			Type:  app.Operator.AlphaType},
			prism.Value{
				Value: env.newExpression(&app.Right),
				Type:  app.Operator.OmegaType})
	}

	call := env.Block.NewCall(
		env.LLDyadicFunctions[app.Operator.LLVMise()],
		env.newExpression(&app.Left),
		env.newExpression(&app.Right))

	return call

}
