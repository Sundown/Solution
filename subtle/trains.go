package subtle

import (
	"fmt"

	"github.com/sundown/solution/palisade"
	"github.com/sundown/solution/prism"
)

func (env Environment) analyseDyadicPartial(expr *palisade.Expression, left, right prism.Type) prism.DyadicFunction {
	return env.BoardTrain(expr, left, right).(prism.DyadicFunction)
}

func (env Environment) analyseMonadicPartial(expr *palisade.Expression, right prism.Type) prism.MonadicFunction {
	return env.BoardTrain(expr, nil, right).(prism.MonadicFunction)
}

func trainLength(expr *palisade.Expression) int {
	if expr.Monadic.Expression == nil {
		return 1
	}

	return 1 + trainLength(expr.Monadic.Expression)
}

func (env Environment) BoardTrain(
	expr *palisade.Expression, left, right prism.Type,
) prism.Function {
	if l := trainLength(expr); l == 2 {
		if left == nil {
			t := env.M2Train(env.FetchMVerb(expr.Monadic.Verb),
				env.FetchMVerb(expr.Monadic.Expression.Monadic.Verb),
				right)

			env.MonadicFunctions[t.Ident()] = &t

			return t
		} else {
			t := env.D2Train(env.FetchMVerb(expr.Monadic.Verb),
				env.FetchDVerb(expr.Monadic.Expression.Monadic.Verb),
				left, right)

			env.DyadicFunctions[t.Ident()] = &t

			return t
		}
	} else if l == 3 {
		if left == nil {
			t := env.M3Train(env.FetchMVerb(expr.Monadic.Verb),
				env.FetchDVerb(expr.Monadic.Expression.Monadic.Verb),
				env.FetchMVerb(expr.Monadic.Expression.Monadic.Expression.Monadic.Verb),
				right)

			env.MonadicFunctions[t.Ident()] = &t

			return t
		} else {
			t := env.D3Train(env.FetchDVerb(expr.Monadic.Verb),
				env.FetchDVerb(expr.Monadic.Expression.Monadic.Verb),
				env.FetchDVerb(expr.Monadic.Expression.Monadic.Expression.Monadic.Verb),
				left, right)

			env.DyadicFunctions[t.Ident()] = &t

			return t
		}
	} else if l%2 == 0 {
		if left == nil {
			t := env.M2Train(env.FetchMVerb(expr.Monadic.Verb),
				env.BoardTrain(expr.Monadic.Expression, left, right).(prism.MonadicFunction),
				right)
			env.MonadicFunctions[t.Ident()] = &t

			return t
		} else {
			t := env.D2Train(env.FetchMVerb(expr.Monadic.Verb),
				env.BoardTrain(expr.Monadic.Expression, left, right).(prism.DyadicFunction),
				left, right)
			env.DyadicFunctions[t.Ident()] = &t

			return t
		}
	} else {
		if left == nil {
			t := env.M3Train(env.FetchMVerb(expr.Monadic.Verb),
				env.FetchDVerb(expr.Monadic.Expression.Monadic.Verb),
				env.BoardTrain(expr.Monadic.Expression.Monadic.Expression,
					left, right).(prism.MonadicFunction), right)
			env.MonadicFunctions[t.Ident()] = &t

			return t
		} else {
			t := env.D3Train(env.FetchDVerb(expr.Monadic.Verb),
				env.FetchDVerb(expr.Monadic.Expression.Monadic.Verb),
				env.BoardTrain(expr.Monadic.Expression.Monadic.Expression,
					left, right).(prism.DyadicFunction), left, right)
			env.DyadicFunctions[t.Ident()] = &t

			return t
		}
	}
}

func Match(e *prism.Type, t *prism.Type) {
	if !(*e).Equals(*t) {
		if !prism.QueryCast(*e, *t) {
			_, err := prism.Delegate(t, e)
			if err != nil {
				panic(*err)
			}
		}
	}
}

// Method for creating the specific function of a Dyadic 3-train with determined types
func (env Environment) D3Train(f, g, h prism.DyadicFunction, APre, BPre prism.Type) prism.DyadicFunction {
	Match(&APre, &f.AlphaType)
	Match(&APre, &h.AlphaType)
	Match(&BPre, &f.OmegaType)
	Match(&BPre, &h.OmegaType)

	if f.Returns.IsAlgebraic() {
		f.Returns = f.Returns.Resolve(f.AlphaType)
	}

	if h.Returns.IsAlgebraic() {
		h.Returns = h.Returns.Resolve(h.AlphaType)
	}

	Match(&f.Returns, &g.AlphaType)
	Match(&h.Returns, &g.OmegaType)

	if g.Returns.IsAlgebraic() {
		g.Returns = g.Returns.Resolve(h.AlphaType /* <- wrong */)
	}

	dy := prism.DyadicFunction{
		Special:     false,
		SkipBuilder: true,
		Inline:      true,
		Name:        prism.Ident{Package: "_", Name: "d3_train_" + fmt.Sprint(env.Iterate())},
		AlphaType:   APre,
		OmegaType:   BPre,
		Returns:     g.Returns,
		PreBody:     nil,
		Body: []prism.Expression{
			prism.MonadicApplication{
				Operator: prism.MonadicFunction{
					Special: false,
					Name:    prism.Ident{Package: "_", Name: "Return"},
					Returns: g.Returns,
				},
				Operand: prism.DyadicApplication{
					Operator: g,
					Left: prism.DyadicApplication{
						Operator: f,
						Left:     prism.Alpha{TypeOf: APre},
						Right:    prism.Omega{TypeOf: BPre},
					},
					Right: prism.DyadicApplication{
						Operator: h,
						Left:     prism.Alpha{TypeOf: APre},
						Right:    prism.Omega{TypeOf: BPre},
					},
				},
			},
		},
	}

	return dy
}

// a (gh) b <=> g (a h b)
func (env Environment) D2Train(g prism.MonadicFunction, h prism.DyadicFunction, APre, BPre prism.Type) prism.DyadicFunction {
	Match(&APre, &h.AlphaType)
	Match(&BPre, &h.OmegaType)

	if h.Returns.IsAlgebraic() {
		h.Returns = h.Returns.Resolve(h.AlphaType)
	}

	Match(&h.Returns, &g.OmegaType)

	if g.Returns.IsAlgebraic() {
		g.Returns = g.Returns.Resolve(h.Returns)
	}

	dy := prism.DyadicFunction{
		Special:     false,
		SkipBuilder: true,
		Inline:      true,
		Name:        prism.Ident{Package: "_", Name: "d2_train_" + fmt.Sprint(env.Iterate())},
		AlphaType:   APre,
		OmegaType:   BPre,
		Returns:     g.Returns,
		PreBody:     nil,
		Body: []prism.Expression{
			prism.MonadicApplication{
				Operator: prism.MonadicFunction{
					Special: false,
					Name:    prism.Ident{Package: "_", Name: "Return"},
					Returns: g.Returns,
				},
				Operand: prism.MonadicApplication{
					Operator: g,
					Operand: prism.DyadicApplication{
						Operator: h,
						Left:     prism.Alpha{TypeOf: APre},
						Right:    prism.Omega{TypeOf: BPre},
					},
				},
			},
		},
	}

	return dy
}

// (fgh) b <=> (f b) g (h b)
func (env Environment) M3Train(f prism.MonadicFunction, g prism.DyadicFunction, h prism.MonadicFunction, BPre prism.Type) prism.MonadicFunction {
	Match(&BPre, &f.OmegaType)
	Match(&BPre, &h.OmegaType)

	if f.Returns.IsAlgebraic() {
		f.Returns = f.Returns.Resolve(f.OmegaType)
	}

	if h.Returns.IsAlgebraic() {
		h.Returns = h.Returns.Resolve(h.OmegaType)
	}

	Match(&f.Returns, &g.AlphaType)
	Match(&h.Returns, &g.OmegaType)

	if g.Returns.IsAlgebraic() {
		g.Returns = g.Returns.Resolve(h.Returns /* <- wrong */)
	}

	dy := prism.MonadicFunction{
		Special:     false,
		SkipBuilder: true,
		Inline:      true,
		Name:        prism.Ident{Package: "_", Name: "m3_train_" + fmt.Sprint(env.Iterate())},
		OmegaType:   BPre,
		Returns:     g.Returns,
		PreBody:     nil,
		Body: []prism.Expression{
			prism.MonadicApplication{
				Operator: prism.MonadicFunction{
					Special: false,
					Name:    prism.Ident{Package: "_", Name: "Return"},
					Returns: g.Returns,
				},
				Operand: prism.DyadicApplication{
					Operator: g,
					Left: prism.MonadicApplication{
						Operator: f,
						Operand:  prism.Omega{TypeOf: BPre},
					},
					Right: prism.MonadicApplication{
						Operator: h,
						Operand:  prism.Omega{TypeOf: BPre},
					},
				},
			},
		},
	}

	return dy
}

// (gh) b <=> g (h b)
func (env Environment) M2Train(g prism.MonadicFunction, h prism.MonadicFunction, BPre prism.Type) prism.MonadicFunction {
	Match(&BPre, &h.OmegaType)

	if h.Returns.IsAlgebraic() {
		h.Returns = h.Returns.Resolve(h.OmegaType)
	}

	Match(&h.Returns, &g.OmegaType)

	if g.Returns.IsAlgebraic() {
		g.Returns = g.Returns.Resolve(h.Returns)
	}

	dy := prism.MonadicFunction{
		Special:     false,
		SkipBuilder: true,
		Inline:      true,
		Name:        prism.Ident{Package: "_", Name: "m2_train_" + fmt.Sprint(env.Iterate())},
		OmegaType:   BPre,
		Returns:     g.Returns,
		PreBody:     nil,
		Body: []prism.Expression{
			prism.MonadicApplication{
				Operator: prism.MonadicFunction{
					Special: false,
					Name:    prism.Ident{Package: "_", Name: "Return"},
					Returns: g.Returns,
				},
				Operand: prism.MonadicApplication{
					Operator: g,
					Operand: prism.MonadicApplication{
						Operator: h,
						Operand:  prism.Omega{TypeOf: BPre},
					},
				},
			},
		},
	}

	return dy
}
