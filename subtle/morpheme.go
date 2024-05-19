package subtle

import (
	"github.com/sundown/solution/palisade"
	"github.com/sundown/solution/prism"
)

func (env Environment) analyseMorphemes(ms *palisade.Morpheme) prism.Expression {
	mor := env.analyseMorpheme(ms)
	if vec, ok := mor.(prism.Vector); ok {
		if len(*vec.Body) == 1 {
			return (*vec.Body)[0]
		}
	}

	return mor
}

func (env Environment) analyseMorpheme(m *palisade.Morpheme) prism.Expression {
	switch {
	case m.Char != nil:
		vec := make([]prism.Expression, len(*m.Char))
		for i, c := range *m.Char {
			vec[i] = prism.Char{Value: string(c[1])}
		}

		return prism.Vector{
			ElementType: prism.VectorType{Type: prism.CharType},
			Body:        &vec,
		}
	case m.Int != nil:
		vec := make([]prism.Expression, len(*m.Int))
		for i, c := range *m.Int {
			vec[i] = prism.Int{Value: c}
		}

		return prism.Vector{
			ElementType: prism.VectorType{Type: prism.IntType},
			Body:        &vec,
		}
	case m.Real != nil:
		vec := make([]prism.Expression, len(*m.Real))
		for i, c := range *m.Real {
			vec[i] = prism.Real{Value: c}
		}

		return prism.Vector{
			ElementType: prism.VectorType{Type: prism.RealType},
			Body:        &vec,
		}
	case m.String != nil:
		outer := make([]prism.Expression, len(*m.String))
		for i, str := range *m.String {
			vec := make([]prism.Expression, len(str)+1)
			for inner, ch := range str {
				vec[inner] = prism.Char{Value: string(ch)}

			}

			// TODO get rid of this, use -n functions in LLVM instead
			vec[len(str)] = prism.Char{Value: "\000"} // null termination for C utils

			outer[i] = prism.Vector{
				ElementType: prism.VectorType{Type: prism.CharType},
				Body:        &vec,
			}
		}

		return prism.Vector{
			ElementType: prism.VectorType{Type: prism.StringType},
			Body:        &outer,
		}
	case m.Alpha != nil:
		if len(*m.Alpha) == 1 {
			if f, ok := env.CurrentFunctionIR.(prism.DyadicFunction); ok {
				return prism.Alpha{
					TypeOf: f.AlphaType,
				}
			}
		}

		prism.Panic("Unreachable")

	case m.Omega != nil:
		if len(*m.Omega) == 1 {
			if f, ok := env.CurrentFunctionIR.(prism.DyadicFunction); ok {
				return prism.Omega{
					TypeOf: f.OmegaType,
				}
			} else if f, ok := env.CurrentFunctionIR.(prism.MonadicFunction); ok {
				return prism.Omega{
					TypeOf: f.OmegaType,
				}
			}
		} else {
			prism.Panic("Unreachable")

		}

	}

	prism.Panic("Other types not implemented")
	panic("unlabelled error")
}
