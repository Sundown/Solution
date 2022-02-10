package subtle

import (
	"sundown/solution/palisade"
	"sundown/solution/prism"

	"github.com/alecthomas/repr"
)

func (env Environment) AnalyseDyadicPartial(expr *palisade.Expression) prism.Expression {
	//---
	f := env.FetchDVerb(expr.Monadic.Verb)
	g := env.FetchDVerb(expr.Monadic.Expression.Monadic.Verb)
	h := env.FetchDVerb(expr.Monadic.Expression.Monadic.Expression.Monadic.Verb)

	A := &f.AlphaType
	B := &f.OmegaType
	X := &f.Returns

	C := &h.AlphaType
	D := &h.OmegaType
	Y := &h.Returns

	E := &g.AlphaType
	F := &g.OmegaType
	//Z := &g.Returns

	Match(A, C)
	Match(B, D)
	Match(X, E)
	Match(Y, F)

	repr.Println(expr)
	return nil
}

func Match(e *prism.Type, t *prism.Type) {
	if !prism.PureMatch((*e), *t) {
		if !prism.QueryCast((*e), *t) {
			tmp := (*e)
			_, err := prism.Delegate(t, &tmp)
			if err != nil {
				prism.Panic(*err)
			}
		}
	}
}
