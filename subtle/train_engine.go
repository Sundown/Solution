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
	var train prism.Function
	first := *expr.Monadic.Applicable
	second := *cdr(expr).Monadic.Applicable

	f := env.analyseApplicable(first, left, right)

	if l := trainLength(expr); l == 2 {
		if left == nil {
			h := env.analyseApplicable(second, nil, right).(prism.MonadicFunction)
			train = env.m2Train(env.analyseApplicable(first, nil, h.Type()).(prism.MonadicFunction),
				h,
				right)
		} else {
			h := env.analyseApplicable(second, left, right).(prism.DyadicFunction)
			train = env.d2Train(env.analyseApplicable(first, left, h.Type()).(prism.MonadicFunction),
				h,
				left, right)
		}
	} else if l == 3 {
		third := *cdr(cdr(expr)).Monadic.Applicable
		if left == nil {
			f := f.(prism.MonadicFunction)
			h := env.analyseApplicable(third, nil, right).(prism.MonadicFunction)
			train = env.m3Train(f,
				env.analyseApplicable(second, f.Type(), h.Type()).(prism.DyadicFunction),
				h,
				right)
		} else {
			f := f.(prism.DyadicFunction)
			h := env.analyseApplicable(third, left, right).(prism.DyadicFunction)

			train = env.d3Train(f,
				env.analyseApplicable(second, f.Type(), h.Type()).(prism.DyadicFunction),
				h,
				left, right)
		}
	} else if l%2 == 0 {
		if left == nil {
			h := env.boardTrain(cdr(expr), nil, right).(prism.MonadicFunction)
			train = env.m2Train(env.analyseApplicable(first, nil, h.Type()).(prism.MonadicFunction),
				h,
				right)
		} else {
			h := env.boardTrain(cdr(expr), left, right).(prism.DyadicFunction)
			train = env.d2Train(env.analyseApplicable(first, left, h.Type()).(prism.MonadicFunction),
				h,
				left, right)
		}
	} else {
		if left == nil {
			f := f.(prism.MonadicFunction)
			h := env.boardTrain(cdr(cdr(expr)), nil, right).(prism.MonadicFunction)
			train = env.m3Train(f,
				env.analyseApplicable(second, f.Type(), h.Type()).(prism.DyadicFunction),
				h, right)
		} else {
			f := f.(prism.DyadicFunction)
			h := env.boardTrain(cdr(cdr(expr)), left, right).(prism.DyadicFunction)
			train = env.d3Train(f,
				env.analyseApplicable(second, f.Type(), h.Type()).(prism.DyadicFunction),
				h, left, right)
		}
	}

	if m, ok := train.(prism.MonadicFunction); ok {
		env.MonadicFunctions[train.Ident()] = &m
	} else if d, ok := train.(prism.DyadicFunction); ok {
		env.DyadicFunctions[train.Ident()] = &d
	}

	return train
}

// Get tail of train expression
func cdr(expr *palisade.Expression) *palisade.Expression {
	if (*expr).Monadic != nil {
		return (*expr).Monadic.Expression
	}

	panic("Abominable expression structure")
}

func match(e *prism.Type, t *prism.Type) {
	if !(*e).Equals(*t) { // perhaps?
		if !prism.QueryCast(*e, *t) { // maybe?
			if _, err := prism.Delegate(t, e); err != nil { // possibly?
				panic(*err)
			}
		}
	}
}
