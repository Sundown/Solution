package subtle

import (
	"sundown/solution/palisade"
	"sundown/solution/prism"

	"github.com/alecthomas/repr"
)

// This returns a new DyadicFunction which is the proper composition of functions within
// the train, arranged according to the number of such functions
// https://aplwiki.com/wiki/Tacit_programming#Trains
func (env Environment) AnalyseDyadicPartial(expr *palisade.Expression, left, right prism.Expression) prism.DyadicFunction {
	g := env.FetchDVerb(expr.Monadic.Expression.Monadic.Verb)
	var dy prism.DyadicFunction
	if expr.Monadic.Expression.Monadic.Expression != nil {
		var h prism.DyadicFunction
		if expr.Monadic.Expression.Monadic.Expression.Monadic != nil && expr.Monadic.Expression.Monadic.Expression.Monadic.Expression != nil {
			//repr.Println(expr.Monadic.Expression.Monadic.Expression.Monadic)
			h = env.AnalyseDyadicPartial(expr.Monadic.Expression.Monadic.Expression, left, right)
		} else {
			h = env.FetchDVerb(expr.Monadic.Expression.Monadic.Expression.Monadic.Verb)
		}
		f := env.FetchDVerb(expr.Monadic.Verb)
		dy = env.D3Train(f, g, h, left, right)
	} else {
		f := env.FetchMVerb(expr.Monadic.Verb)
		dy = env.D2Train(f, g, left, right)
	}

	env.DyadicFunctions[dy.Ident()] = &dy

	repr.Println(dy)

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
