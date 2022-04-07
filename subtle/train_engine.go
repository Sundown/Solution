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
	} else {
		return 1
	}
}

func (env Environment) boardTrain(
	expr *palisade.Expression, left, right prism.Type,
) prism.Function {
	var t prism.Function
	if l := trainLength(expr); l == 2 {
		if left == nil {
			t = env.m2Train(env.analyseApplicable(*expr.Monadic.Applicable, left, right).(prism.MonadicFunction),
				env.analyseApplicable(*expr.Monadic.Expression.Monadic.Applicable, left, right).(prism.MonadicFunction),
				right)
		} else {
			t = env.d2Train(env.analyseApplicable(*expr.Monadic.Applicable, left, right).(prism.MonadicFunction),
				env.analyseApplicable(*expr.Monadic.Expression.Monadic.Applicable, left, right).(prism.DyadicFunction),
				left, right)
		}
	} else if l == 3 {
		if left == nil {
			//first := innerExpression(expr)

			t = env.m3Train(env.analyseApplicable(*expr.Monadic.Applicable, left, right).(prism.MonadicFunction),
				env.analyseApplicable(*expr.Monadic.Expression.Monadic.Applicable, left, right).(prism.DyadicFunction),
				env.analyseApplicable(*expr.Monadic.Expression.Monadic.Expression.Monadic.Applicable, left, right).(prism.MonadicFunction),
				right)
		} else {

			t = env.d3Train(env.analyseApplicable(*expr.Monadic.Applicable, left, right).(prism.DyadicFunction),
				env.analyseApplicable(*expr.Monadic.Expression.Monadic.Applicable, left, right).(prism.DyadicFunction),
				env.analyseApplicable(*expr.Monadic.Expression.Monadic.Expression.Monadic.Applicable, left, right).(prism.DyadicFunction),
				left, right)
		}
	} else if l%2 == 0 {
		if left == nil {
			t = env.m2Train(env.analyseApplicable(*expr.Monadic.Applicable, left, right).(prism.MonadicFunction),
				env.boardTrain(expr.Monadic.Expression, left, right).(prism.MonadicFunction),
				right)
		} else {

			t = env.d2Train(env.analyseApplicable(*expr.Monadic.Applicable, left, right).(prism.MonadicFunction),
				env.boardTrain(expr.Monadic.Expression, left, right).(prism.DyadicFunction),
				left, right)
		}
	} else {
		if left == nil {
			t = env.m3Train(env.analyseApplicable(*expr.Monadic.Applicable, left, right).(prism.MonadicFunction),
				env.analyseApplicable(*expr.Monadic.Expression.Monadic.Applicable, left, right).(prism.DyadicFunction),
				env.boardTrain(expr.Monadic.Expression.Monadic.Expression,
					left, right).(prism.MonadicFunction), right)
		} else {
			t = env.d3Train(env.analyseApplicable(*expr.Monadic.Applicable, left, right).(prism.DyadicFunction),
				env.analyseApplicable(*expr.Monadic.Expression.Monadic.Applicable, left, right).(prism.DyadicFunction),
				env.boardTrain(expr.Monadic.Expression.Monadic.Expression,
					left, right).(prism.DyadicFunction), left, right)
		}
	}

	if m, ok := t.(prism.MonadicFunction); ok {
		env.MonadicFunctions[t.Ident()] = &m
	} else if d, ok := t.(prism.DyadicFunction); ok {
		env.DyadicFunctions[t.Ident()] = &d
	}

	return t
}

/* func innerExpression(expr *palisade.Expression) *palisade.Expression {
	if (*expr).Monadic != nil {
		return (*expr).Monadic.Expression
	} else if (*expr).Operator != nil {
		return (*expr).Operator.Expression
	}

	panic("Abominable expression structure")
} */

/* func (env *Environment) operatorOrMonadic(expr *palisade.Expression, right prism.Type) prism.MonadicFunction {
	if expr.Monadic != nil {
		return env.FetchMVerb(expr.Monadic.Verb)
	} else if expr.Operator != nil {
		return env.operatorToFunction(env.analyseMonadicOperator(expr.Operator, right))
	}

	panic("aaa")
} */

func match(e *prism.Type, t *prism.Type) {
	if !(*e).Equals(*t) { // perhaps?
		if !prism.QueryCast(*e, *t) { // maybe?
			if _, err := prism.Delegate(t, e); err != nil { // possibly?
				panic(*err)
			}
		}
	}
}
