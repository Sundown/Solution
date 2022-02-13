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
		var h prism.DyadicFunction
		if expr.Monadic.Expression.Monadic.Expression.Monadic != nil &&
			expr.Monadic.Expression.Monadic.Expression.Monadic.Expression != nil {
			h = env.AnalyseDyadicPartial(expr.Monadic.Expression.Monadic.Expression, left, right)
		} else {
			h = env.FetchDVerb(expr.Monadic.Expression.Monadic.Expression.Monadic.Verb)
		}

		dy = env.D3Train(env.FetchDVerb(expr.Monadic.Verb), g, h, left, right)
	} else {
		dy = env.D2Train(env.FetchMVerb(expr.Monadic.Verb), g, left, right)
	}

	env.DyadicFunctions[dy.Ident()] = &dy

	return dy
}

func Match(e *prism.Type, t *prism.Type) {
	if !prism.PureMatch(*e, *t) {
		if !prism.QueryCast(*e, *t) {
			_, err := prism.Delegate(t, e)
			if err != nil {
				panic(*err)
			}
		}
	}
}
