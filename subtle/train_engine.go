package subtle

import (
	"github.com/sundown/solution/palisade"
	"github.com/sundown/solution/prism"
)

func (env Environment) analyseDyadicPartial(expr *palisade.Expression, left, right prism.Type) prism.DyadicFunction {
	return env.boardTrain(expr, left, right).(prism.DyadicFunction)
}

func (env Environment) analyseMonadicPartial(expr *palisade.Expression, right prism.Type) prism.MonadicFunction {
	return env.boardTrain(expr, nil, right).(prism.MonadicFunction)
}

func trainLength(expr *palisade.Expression) int {
	if expr.Monadic != nil && expr.Monadic.Expression != nil {
		return 1 + trainLength(expr.Monadic.Expression)
	} else if expr.Operator != nil && expr.Operator.Expression != nil {
		return 1 + trainLength(expr.Operator.Expression)
	} else {
		return 1
	}
}

func (env Environment) boardTrain(
	expr *palisade.Expression, left, right prism.Type,
) prism.Function {
	if l := trainLength(expr); l == 2 {
		if left == nil {
			t := env.m2Train(env.FetchMVerb(expr.Monadic.Verb),
				env.FetchMVerb(expr.Monadic.Expression.Monadic.Verb),
				right)

			env.MonadicFunctions[t.Ident()] = &t

			return t
		}

		t := env.d2Train(env.FetchMVerb(expr.Monadic.Verb),
			env.FetchDVerb(expr.Monadic.Expression.Monadic.Verb),
			left, right)

		env.DyadicFunctions[t.Ident()] = &t

		return t
	} else if l == 3 {
		if left == nil {
			first := innerExpression(expr)

			t := env.m3Train(env.operatorOrMonadic(expr, right),
				env.FetchDVerb(first.Monadic.Verb),
				env.operatorOrMonadic(innerExpression(first), right),
				right)

			env.MonadicFunctions[t.Ident()] = &t

			return t
		}

		t := env.d3Train(env.FetchDVerb(expr.Monadic.Verb),
			env.FetchDVerb(expr.Monadic.Expression.Monadic.Verb),
			env.FetchDVerb(expr.Monadic.Expression.Monadic.Expression.Monadic.Verb),
			left, right)

		env.DyadicFunctions[t.Ident()] = &t

		return t
	} else if l%2 == 0 {
		if left == nil {
			t := env.m2Train(env.FetchMVerb(expr.Monadic.Verb),
				env.boardTrain(expr.Monadic.Expression, left, right).(prism.MonadicFunction),
				right)
			env.MonadicFunctions[t.Ident()] = &t

			return t
		}

		t := env.d2Train(env.FetchMVerb(expr.Monadic.Verb),
			env.boardTrain(expr.Monadic.Expression, left, right).(prism.DyadicFunction),
			left, right)
		env.DyadicFunctions[t.Ident()] = &t

		return t
	} else {
		if left == nil {
			t := env.m3Train(env.FetchMVerb(expr.Monadic.Verb),
				env.FetchDVerb(expr.Monadic.Expression.Monadic.Verb),
				env.boardTrain(expr.Monadic.Expression.Monadic.Expression,
					left, right).(prism.MonadicFunction), right)
			env.MonadicFunctions[t.Ident()] = &t

			return t
		}

		t := env.d3Train(env.FetchDVerb(expr.Monadic.Verb),
			env.FetchDVerb(expr.Monadic.Expression.Monadic.Verb),
			env.boardTrain(expr.Monadic.Expression.Monadic.Expression,
				left, right).(prism.DyadicFunction), left, right)
		env.DyadicFunctions[t.Ident()] = &t

		return t
	}
}

func innerExpression(expr *palisade.Expression) *palisade.Expression {
	if (*expr).Monadic != nil {
		return (*expr).Monadic.Expression
	} else if (*expr).Operator != nil {
		return (*expr).Operator.Expression
	}

	panic("aaa")
}

func (env *Environment) operatorOrMonadic(expr *palisade.Expression, right prism.Type) prism.MonadicFunction {
	if expr.Monadic != nil {
		return env.FetchMVerb(expr.Monadic.Verb)
	} else if expr.Operator != nil {
		return env.operatorToFunction(env.analyseDyadicOperator(expr.Operator, right))
	}

	panic("aaa")
}

func match(e *prism.Type, t *prism.Type) {
	if !(*e).Equals(*t) {
		if !prism.QueryCast(*e, *t) {
			_, err := prism.Delegate(t, e)
			if err != nil {
				panic(*err)
			}
		}
	}
}
