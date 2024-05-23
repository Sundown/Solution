package subtle

import (
	"fmt"

	"github.com/sundown/solution/prism"
)

// The 4 types of trains are defined herein
// These should only be invoked by boardTrain()
//
// https://aplwiki.com/wiki/Tacit_programming#3-trains
//
// d3Train: dyadic 3 train
// 		a (fgh) b <=> (a f b) g (a h b)
//
// d2Train: dyadic 2 train
// 		a (gh) b <=> g (a h b)
//
// m3Train: monadic 3 train
// 		(fgh) b <=> (f b) g (h b)
//
// m2Train: monadic 2 train
// 		(gh) b <=> g (h b)

// a (fgh) b <=> (a f b) g (a h b)
func (env *Environment) d3Train(f, g, h prism.DyadicFunction, APre, BPre prism.Type) prism.DyadicFunction {
	match(&APre, &f.AlphaType)
	match(&APre, &h.AlphaType)
	match(&BPre, &f.OmegaType)
	match(&BPre, &h.OmegaType)

	if f.Returns.IsAlgebraic() {
		f.Returns = f.Returns.Resolve(f.AlphaType)
	}

	if h.Returns.IsAlgebraic() {
		h.Returns = h.Returns.Resolve(h.AlphaType)
	}

	match(&f.Returns, &g.AlphaType)
	match(&h.Returns, &g.OmegaType)

	if g.Returns.IsAlgebraic() {
		g.Returns = g.Returns.Resolve(h.AlphaType /* <- wrong */)
	}

	dy := prism.DyadicFunction{
		Attribute: prism.Attribute{
			Special:     false,
			SkipBuilder: true,
			ForceInline: true,
		},
		Name:      prism.Ident{Package: "_", Name: "d3_train_" + fmt.Sprint(env.Iterate())},
		AlphaType: APre,
		OmegaType: BPre,
		Returns:   g.Returns,
		PreBody:   nil,
		Body: []prism.Expression{
			prism.MonadicApplication{
				Operator: prism.MonadicFunction{
					Name:    prism.Ident{Package: "_", Name: "←"},
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
func (env *Environment) d2Train(g prism.MonadicFunction, h prism.DyadicFunction, APre, BPre prism.Type) prism.DyadicFunction {
	match(&APre, &h.AlphaType)
	match(&BPre, &h.OmegaType)

	if h.Returns.IsAlgebraic() {
		h.Returns = h.Returns.Resolve(h.AlphaType)
	}

	match(&h.Returns, &g.OmegaType)

	if g.Returns.IsAlgebraic() {
		g.Returns = g.Returns.Resolve(h.Returns)
	}

	dy := prism.DyadicFunction{
		Attribute: prism.Attribute{
			Special:     false,
			SkipBuilder: true,
			ForceInline: true,
		},
		Name:      prism.Ident{Package: "_", Name: "d2_train_" + fmt.Sprint(env.Iterate())},
		AlphaType: APre,
		OmegaType: BPre,
		Returns:   g.Returns,
		PreBody:   nil,
		Body: []prism.Expression{
			prism.MonadicApplication{
				Operator: prism.MonadicFunction{
					Name:    prism.Ident{Package: "_", Name: "←"},
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
func (env *Environment) m3Train(f prism.MonadicFunction, g prism.DyadicFunction, h prism.MonadicFunction, BPre prism.Type) prism.MonadicFunction {
	match(&BPre, &f.OmegaType)
	match(&BPre, &h.OmegaType)

	if f.Returns.IsAlgebraic() {
		f.Returns = f.Returns.Resolve(f.OmegaType)
	}

	if h.Returns.IsAlgebraic() {
		h.Returns = h.Returns.Resolve(h.OmegaType)
	}

	match(&f.Returns, &g.AlphaType)
	match(&h.Returns, &g.OmegaType)

	if g.Returns.IsAlgebraic() {
		g.Returns = g.Returns.Resolve(h.Returns /* <- wrong, TODO */)
	}

	dy := prism.MonadicFunction{
		Attribute: prism.Attribute{
			Special:     false,
			SkipBuilder: true,
			ForceInline: true,
		},
		Name:      prism.Ident{Package: "_", Name: "m3_train_" + fmt.Sprint(env.Iterate())},
		OmegaType: BPre,
		Returns:   g.Returns,
		PreBody:   nil,
		Body: []prism.Expression{
			prism.MonadicApplication{
				Operator: prism.MonadicFunction{
					Name:    prism.Ident{Package: "_", Name: "←"},
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
func (env *Environment) m2Train(g prism.MonadicFunction, h prism.MonadicFunction, BPre prism.Type) prism.MonadicFunction {
	match(&BPre, &h.OmegaType)

	if h.Returns.IsAlgebraic() {
		h.Returns = h.Returns.Resolve(h.OmegaType)
	}

	match(&h.Returns, &g.OmegaType)

	if g.Returns.IsAlgebraic() {
		g.Returns = g.Returns.Resolve(h.Returns)
	}

	dy := prism.MonadicFunction{
		Attribute: prism.Attribute{
			Special:     false,
			SkipBuilder: true,
			ForceInline: true,
		},
		Name:      prism.Ident{Package: "_", Name: "m2_train_" + fmt.Sprint(env.Iterate())},
		OmegaType: BPre,
		Returns:   g.Returns,
		PreBody:   nil,
		Body: []prism.Expression{
			prism.MonadicApplication{
				Operator: prism.MonadicFunction{
					Name:    prism.Ident{Package: "_", Name: "←"},
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
