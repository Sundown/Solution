package subtle

import (
	"github.com/sundown/solution/palisade"
	"github.com/sundown/solution/prism"
)

func (env *Environment) analyseExpression(e *palisade.Expression) prism.Expression {
	if e.Monadic != nil {
		if e.Monadic.Expression == nil {
			return env.FetchMVerb(e.Monadic.Applicable.Verb)
		}
		return env.analyseMonadic(e.Monadic)
	} else if e.Dyadic != nil {
		return env.analyseDyadic(e.Dyadic)
	} else if e.Morphemes != nil {
		return env.analyseMorphemes(e.Morphemes)
	} else if e.Bool != nil {
		if len(*e.Bool) == 1 {
			return prism.Bool{Value: (*e.Bool)[0] == "true"}
		}

		vec := make([]prism.Expression, len(*e.Bool))
		for i, c := range *e.Bool {
			vec[i] = prism.Bool{Value: c == "true"}
		}

		return prism.Vector{
			ElementType: prism.VectorType{Type: prism.BoolType},
			Body:        &vec,
		}
	}

	prism.Panic("unreachable")
	panic("unlabelled error")
}
