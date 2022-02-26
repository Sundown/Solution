package subtle

import (
	"github.com/sundown/solution/palisade"
	"github.com/sundown/solution/prism"
)

func (env Environment) AnalyseMonadic(m *palisade.Monadic) (app prism.MonadicApplication) {
	if m.Verb != nil {
		return env.AnalyseStandardMonadic(m)
	} else if m.Subexpr != nil {
		return env.AnalysePartialMonadic(m)
	} else {
		panic("unreachable")
	}
}

func (env Environment) AnalysePartialMonadic(m *palisade.Monadic) (app prism.MonadicApplication) {
	app = prism.MonadicApplication{
		Operator: env.AnalyseExpression(m.Subexpr).(prism.MonadicFunction),
		Operand:  env.AnalyseExpression(m.Expression),
	}

	tmp := app.Operand.Type()
	resolved_right, err := prism.Delegate(&app.Operator.OmegaType, &tmp)
	if err != nil {
		prism.Panic(*err)
	}

	if app.Operator.Returns.IsAlgebraic() {
		app.Operator.Returns = app.Operator.Returns.Resolve(resolved_right)
	}

	return app
}

func (env Environment) AnalyseStandardMonadic(m *palisade.Monadic) (app prism.MonadicApplication) {
	fn := env.FetchMVerb(m.Verb)

	expr := env.AnalyseExpression(m.Expression)

	tmp := expr.Type()
	resolved_right, err := prism.Delegate(&fn.OmegaType, &tmp)

	if !tmp.Equals(fn.OmegaType) {
		if !prism.QueryCast(tmp, fn.OmegaType) {
			tmp := tmp
			t, err := prism.Delegate(&fn.OmegaType, &tmp)
			_ = t
			if err != nil {
				prism.Panic(*err)
			}
		} else {
			expr = prism.DelegateCast(expr, fn.OmegaType)
		}
	}

	if err != nil {
		prism.Panic(*err)
	}

	if fn.Returns.IsAlgebraic() {
		fn.Returns = fn.Returns.Resolve(resolved_right)
	}

	if fn.Name.Package == "_" && fn.Name.Name == "Return" {
		if !env.CurrentFunctionIR.Type().Equals(fn.Returns) {
			if !env.CurrentFunctionIR.Type().IsAlgebraic() {
				panic("Return recieves " + fn.Returns.String() + " which does not match determined-function's type " + env.CurrentFunctionIR.Type().String())
			} else {
				panic("Not implemented, pain")
			}
		}
	}

	return prism.MonadicApplication{
		Operator: fn,
		Operand:  expr,
	}
}
