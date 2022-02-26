package subtle

import (
	"fmt"

	"github.com/sundown.solution/prism"
)

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

func (env Environment) D2Train(g prism.MonadicFunction, h prism.DyadicFunction, APre, BPre prism.Type) prism.DyadicFunction {
	Match(&APre, &h.AlphaType)
	Match(&BPre, &h.OmegaType)

	if h.Returns.IsAlgebraic() {
		h.Returns = h.Returns.Resolve(h.AlphaType)
	}

	Match(&h.Returns, &g.OmegaType)

	if g.Returns.IsAlgebraic() {
		g.Returns = g.Returns.Resolve(h.AlphaType /* <- wrong */)
	}

	dy := prism.DyadicFunction{
		Special:     false,
		SkipBuilder: true,
		Inline:      true,
		Name:        prism.Ident{Package: "_", Name: "d3_train"},
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

//func (env Environment) M3Train(f, g, h prism.Expression) prism.MonadicFunction
//func (env Environment) M2Train(f, g prism.Expression) prism.MonadicFunction
