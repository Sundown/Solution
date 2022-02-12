package subtle

import (
	"sundown/solution/prism"
)

// Method for creating the specific function of a Dyadic 3-train with determined types
func (env Environment) D3Train(f, g, h prism.DyadicFunction, a, b prism.Expression) prism.DyadicFunction {
	APre := a.Type()
	BPre := b.Type()

	Match(&APre, &f.AlphaType)
	Match(&APre, &h.AlphaType)
	Match(&BPre, &f.OmegaType)
	Match(&BPre, &h.OmegaType)

	if prism.PredicateGenericType(f.Returns) {
		f.Returns = prism.IntegrateGenericType(f.AlphaType, f.Returns)
	}

	if prism.PredicateGenericType(h.Returns) {
		h.Returns = prism.IntegrateGenericType(h.AlphaType, h.Returns)
	}

	Match(&f.Returns, &g.AlphaType)
	Match(&h.Returns, &g.OmegaType)

	if prism.PredicateGenericType(g.Returns) {
		g.Returns = prism.IntegrateGenericType(h.AlphaType /* <- wrong */, g.Returns)
	}

	dy := prism.DyadicFunction{
		Special:   false,
		Name:      prism.Ident{Package: "_", Name: "d3_train"},
		AlphaType: APre,
		OmegaType: BPre,
		Returns:   g.Returns,
		PreBody:   nil,
		Body: []prism.Expression{
			prism.MApplication{
				Operator: prism.MonadicFunction{
					Special: false,
					Name:    prism.Ident{Package: "_", Name: "Return"},
					Returns: g.Returns,
				},
				Operand: prism.DApplication{
					Operator: g,
					Left: prism.DApplication{
						Operator: f,
						Left:     prism.Alpha{TypeOf: APre},
						Right:    prism.Omega{TypeOf: BPre},
					},
					Right: prism.DApplication{
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

func (env Environment) D2Train(g prism.MonadicFunction, h prism.DyadicFunction) prism.DyadicFunction {
	A := &h.AlphaType
	B := &h.OmegaType
	Y := &h.Returns

	C := &g.OmegaType
	Z := &g.Returns

	Match(Y, C)

	return prism.DyadicFunction{
		Special:   false,
		Name:      prism.Ident{Package: "_", Name: "d2_train"},
		AlphaType: *A,
		OmegaType: *B,
		Returns:   *Z,
		PreBody:   nil,
		Body: []prism.Expression{
			prism.MApplication{
				Operator: prism.MonadicFunction{
					Special: false,
					Name:    prism.Ident{Package: "_", Name: "Return"},
					Returns: *Z,
				},
				Operand: prism.MApplication{
					Operator: g,
					Operand: prism.DApplication{
						Operator: h,
						Left:     prism.Alpha{TypeOf: *A},
						Right:    prism.Omega{TypeOf: *B},
					},
				},
			},
		},
	}
}

//func (env Environment) M3Train(f, g, h prism.Expression) prism.MonadicFunction
//func (env Environment) M2Train(f, g prism.Expression) prism.MonadicFunction
