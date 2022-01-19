package subtle

import (
	"sundown/solution/palisade"
	"sundown/solution/prism"
)

type Environment struct {
	*prism.Environment
}

func Init(result *palisade.PalisadeResult) prism.Environment {
	env := Environment{prism.NewEnvironment()}

	for _, stmt := range result.Statements {
		// First pass, declare all functions so that they
		// may be referenced before they are defined (text-wise)
		if f := stmt.Function; f != nil {
			// palisade.Function is agnostic to arity
			// containing either monadic or dyadic
			env.InternFunction(*f)
		}
	}

	for _, f := range env.DFunctions {
		env.AnalyseDBody(f)
	}

	return *env.Environment
}

// Intern either monadic or dyadic function header into environment
// will not handle body of function, declarations only
func (env Environment) InternFunction(f palisade.Function) {
	if f.Dyadic != nil {
		fn := prism.DFunction{
			Name:      prism.Intern(*f.Dyadic.Ident),
			AlphaType: env.SubstantiateType(*f.Dyadic.Alpha),
			OmegaType: env.SubstantiateType(*f.Dyadic.Omega),
			Returns:   env.SubstantiateType(*f.Returns),
			PreBody:   f.Body,
		}
		// TODO Perform check that it doesn't already exist
		env.DFunctions[fn.Name] = &fn
	} else if f.Monadic != nil {
		fn := prism.MFunction{
			Name:      prism.Intern(*f.Monadic.Ident),
			OmegaType: env.SubstantiateType(*f.Monadic.Omega),
			Returns:   env.SubstantiateType(*f.Returns),
			PreBody:   f.Body,
		}
		// TODO Perform check that it doesn't already exist
		env.MFunctions[fn.Name] = &fn
	}
}

func (env Environment) AnalyseDBody(f *prism.DFunction) {
	for _, expr := range *f.PreBody {
		f.Body = append(f.Body, env.AnalyseExpression(&expr))
	}
}

func (env Environment) AnalyseExpression(e *palisade.Expression) prism.Expression {
	if e.Monadic != nil {
		return env.AnalyseMonadic(e.Monadic)
	} else if e.Dyadic != nil {
		return env.AnalyseDyadic(e.Dyadic)
	} else if e.Morphemes != nil {
		return env.AnalyseMorphemes(e.Morphemes)
	}

	panic("unreachable")
}

func (env Environment) FetchVerb(v *palisade.Verb) prism.Expression {
	if found, ok := env.MFunctions[prism.Intern(*v.Ident)]; ok {
		return *found
	} else if found, ok := env.DFunctions[prism.Intern(*v.Ident)]; ok {
		return *found
	}

	panic("Verb not found")
}

func (env Environment) AnalyseMorphemes(ms *[]palisade.Morpheme) prism.Expression {
	if len(*ms) == 1 {
		return env.AnalyseMorpheme(&(*ms)[0])
	}

	vec := make([]prism.Expression, len(*ms))
	for i, m := range *ms {
		vec[i] = env.AnalyseMorpheme(&m)
	}

	return prism.Vector{
		ElementType: prism.VectorType{vec[0].Type()},
		Body:        &vec,
	}
}

func (env Environment) AnalyseMorpheme(m *palisade.Morpheme) prism.Expression {
	switch {
	case m.Char != nil:
		return prism.Char{string(*m.Char)}
	case m.Int != nil:
		return prism.Int{*m.Int}
	case m.Real != nil:
		return prism.Real{*m.Real}
	case m.String != nil:
		return prism.String{string(*m.String)}
	case m.Subexpr != nil:
		return env.AnalyseExpression(m.Subexpr)
	}

	panic("Other types not implemented")
}
