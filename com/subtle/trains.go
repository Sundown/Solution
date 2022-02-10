package subtle

import "sundown/solution/prism"

func (env Environment) D3Train(f, g, h prism.DyadicFunction) prism.DyadicFunction {
	A := &f.AlphaType
	B := &f.OmegaType
	X := &f.Returns

	C := &h.AlphaType
	D := &h.OmegaType
	Y := &h.Returns

	E := &g.AlphaType
	F := &g.OmegaType
	Z := &g.Returns

	Match(A, C)
	Match(B, D)
	Match(X, E)
	Match(Y, F)

	return prism.DyadicFunction{
		Special:   false,
		Name:      prism.Ident{Package: "_", Name: "d3_train"},
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
				Operand: prism.DApplication{
					Operator: g,
					Left: prism.DApplication{
						Operator: f,
						Left:     prism.Alpha{TypeOf: *A},
						Right:    prism.Omega{TypeOf: *B},
					},
					Right: prism.DApplication{
						Operator: h,
						Left:     prism.Alpha{TypeOf: *A},
						Right:    prism.Omega{TypeOf: *B},
					},
				},
			},
		},
	}
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
func (env Environment) M3Train(f, g, h prism.Expression) prism.MonadicFunction
func (env Environment) M2Train(f, g prism.Expression) prism.MonadicFunction
