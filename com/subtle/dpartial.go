package subtle

import (
	"sundown/solution/palisade"
	"sundown/solution/prism"
)

// This returns a new DyadicFunction which is the proper composition of functions within
// the train, arranged according to the number of such functions
// https://aplwiki.com/wiki/Tacit_programming#Trains
func (env Environment) AnalyseDyadicPartial(expr *palisade.Expression, left, right prism.Expression) prism.DyadicFunction {
	g := env.FetchDVerb(expr.Monadic.Expression.Monadic.Verb)

	var dy prism.DyadicFunction
	if expr.Monadic.Expression.Monadic.Expression != nil {
		f := env.FetchDVerb(expr.Monadic.Verb)
		h := env.FetchDVerb(expr.Monadic.Expression.Monadic.Expression.Monadic.Verb)

		dy = env.D3Train(f, g, h, left, right)
	} else {
		f := env.FetchMVerb(expr.Monadic.Verb)

		X := &f.Returns

		F := &g.OmegaType
		//Z := &g.Returns

		UnifyTypes(X, F)

		dy = env.D2Train(f, g)
	}

	env.DyadicFunctions[dy.Ident()] = &dy

	return dy
}

func UnifyTypes(e *prism.Type, t *prism.Type) {
	if !prism.PureMatch((*e), *t) {
		//if !prism.QueryCast((*e), *t) {
		if a, ok := (*e).(prism.SumType); ok {
			if b, ok := (*t).(prism.SumType); ok {
				in := prism.TypeIntersection(a, b)
				if len(in.Types) > 0 {
					unified := prism.Type(in)
					e = &unified
					t = &unified
					return
				}
			}
		}

		tmp := (*e)
		_, err := prism.Delegate(t, &tmp)
		if err != nil {
			panic(*err)
		}
		//}
	}
}

func Match(e *prism.Type, t *prism.Type) {
	if !prism.PureMatch((*e), *t) {
		if !prism.QueryCast((*e), *t) {
			tmp := (*e)
			_, err := prism.Delegate(t, &tmp)
			if err != nil {
				panic(*err)
			}
		}
	}
}
