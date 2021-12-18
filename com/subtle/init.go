package subtle

import (
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
		exprs = append(exprs, env.AnalyseExpression(&e))
	}

	return &exprs
}

func (env Environment) AnalyseExpression(e *prism.Expression) prism.Expression
