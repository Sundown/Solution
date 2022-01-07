package subtle

import (
	"sundown/solution/oversight"
	"sundown/solution/prism"
)

type Environment struct {
	*prism.Environment
}

func Init(env *prism.Environment) *prism.Environment {
	if env == nil {
		panic("why")
	}

	nenv := &Environment{env}

	for _, f := range env.Functions {
		if f.Body != nil {
			panic("Body should be empty")
		}

		f.Body = nenv.AnalyseBody(f)
	}

	return env
}

func (env Environment) AnalyseBody(f *prism.Function) *[]prism.Expression {
	exprs := []prism.Expression{}
	for _, e := range *f.PreBody {
		exprs = append(exprs, env.AnalyseExpression(e))
	}

	return &exprs
}

func (env Environment) AnalyseExpression(e prism.Expression) prism.Expression {
	switch e.(type) {
	case prism.Dangle:
		return env.AnalyseDangle(e)
	}
	return e
}

func (env Environment) AnalyseDangle(e prism.Expression) prism.Expression {
	if _, fn := e.(prism.Dangle).Outer.(prism.Function); fn {
		panic("cant chain yet")
	}

	if inner := e.(prism.Dangle).Inner; inner != nil {
		// test dyadic status
		if app, ok := inner.(prism.Application); ok {
			return env.AnalyseDyadicApplication(app, inner)
		}

	}
	// this should possible be Expression not Atom
	return env.AnalyseAtom(e)

}

func (env Environment) AnalyseDyadicApplication(inner prism.Expression, outer prism.Expression) prism.Expression {
	fn := inner.(prism.Function)

	// TODO
	// Handle the case of the operand being empty, non-applied function
	// pain

	if fn.AlphaType == nil {
		oversight.Panic("Trying to call monadic function with 2 arguments")
	}

	if !prism.EqType(fn.AlphaType, outer.Type()) {
		oversight.Panic("Trying to call function with wrong left type")
	}

	if !prism.EqType(fn.OmegaType, inner.(prism.Application).Operand.Type()) {
		oversight.Panic("Trying to call function with wrong right type")
	}

	return prism.Dyadic{
		Operator: fn,
		Left:     outer,
		Right:    inner.(prism.Application).Operand,
	}

}

func (env Environment) AnalyseAtom(e prism.Expression) prism.Expression
