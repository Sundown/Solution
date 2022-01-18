package subtle

import (
	"sundown/solution/oversight"
	"sundown/solution/palisade"
	"sundown/solution/prism"
)

type Environment struct {
	*prism.Environment
}

func Init(result *palisade.PalisadeResult) prism.Environment {
	env := Environment{}

	for _, stmt := range result.Statements {
		if f := stmt.FnDef; f != nil {
			if f.Body != nil {
				panic("Body should be empty")
			}

			fn := env.InternFunction(*f)

			env.Functions[fn.Name] = fn

		}
	}

	return *env.Environment
}

func (env Environment) InternFunction(f palisade.FnDef) prism.Function {
	if f.C != nil {
		// Dyadic, has 2 types
		return prism.Function{
			Name:     f.Name,
	}
}

func (env Environment) AnalyseBody(f *prism.Function) *[]prism.Expression {
	exprs := []prism.Expression{}
	for _, e := range *f.PreBody {
		exprs = append(exprs, env.AnalyseExpression(e))
	}

	return &exprs
}

func (env Environment) AnalyseExpression(e prism.Expression) prism.Expression

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
