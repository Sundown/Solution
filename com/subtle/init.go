package subtle

import (
	"sundown/solution/palisade"
	"sundown/solution/prism"
)

type Environment struct {
	*prism.Environment
}

func Parse(penv *prism.Environment) *prism.Environment {
	env := Environment{penv}

	for _, stmt := range env.LexResult.Environmentments {
		// First pass, declare all functions so that they
		// may be referenced before they are defined (text-wise)
		// Handle compiler directives
		if d := stmt.Directive; d != nil {
			if *d.Command == "Package" {
				env.Output = *d.Value
			}
		}

		if f := stmt.Function; f != nil {
			// palisade.Function is agnostic to arity
			// containing either monadic or dyadic
			env.InternFunction(*f)
		}
	}

	for _, f := range env.DyadicFunctions {
		env.AnalyseDBody(f)
	}

	for _, f := range env.MonadicFunctions {
		env.AnalyseMBody(f)
	}

	return env.Environment
}

// Intern either monadic or dyadic function header into environment
// will not handle body of function, declarations only
func (env Environment) InternFunction(f palisade.Function) {
	if f.Dyadic != nil {
		fn := prism.DyadicFunction{
			Name:      prism.Intern(*f.Dyadic.Ident),
			AlphaType: env.SubstantiateType(*f.Dyadic.Alpha),
			OmegaType: env.SubstantiateType(*f.Dyadic.Omega),
			Returns:   env.SubstantiateType(*f.Returns),
			PreBody:   f.Body,
		}
		// TODO Perform check that it doesn't already exist
		env.DyadicFunctions[fn.Name] = &fn
	} else if f.Monadic != nil {
		fn := prism.MonadicFunction{
			Name:      prism.Intern(*f.Monadic.Ident),
			OmegaType: env.SubstantiateType(*f.Monadic.Omega),
			Returns:   env.SubstantiateType(*f.Returns),
			PreBody:   f.Body,
		}
		// TODO Perform check that it doesn't already exist
		env.MonadicFunctions[fn.Name] = &fn
	}
}

func (env Environment) AnalyseDBody(f *prism.DyadicFunction) {
	if f.Special {
		return
	}

	env.CurrentFunctionIR = *f

	for _, expr := range *f.PreBody {
		f.Body = append(f.Body, env.AnalyseExpression(&expr))
	}
}

func (env Environment) AnalyseMBody(f *prism.MonadicFunction) {
	if f.Special {
		return
	}

	env.CurrentFunctionIR = *f

	for _, expr := range *f.PreBody {
		f.Body = append(f.Body, env.AnalyseExpression(&expr))
	}
}

func (env Environment) AnalyseExpression(e *palisade.Expression) prism.Expression {
	if e.Monadic != nil {
		return env.AnalyseMonadic(e.Monadic)
	} else if e.Dyadic != nil {
		if *e.Dyadic.Verb.Ident.Ident == "Map" {
			return env.AnalyseDyadicOperator(e.Dyadic)
		}

		return env.AnalyseDyadic(e.Dyadic)
	} else if e.Morphemes != nil {
		return env.AnalyseMorphemes(e.Morphemes)
	}

	panic("unreachable")
}

func (env Environment) FetchVerb(v *palisade.Ident) prism.Expression {
	if found, ok := env.MonadicFunctions[prism.Intern(*v)]; ok {
		return *found
	} else if found, ok := env.DyadicFunctions[prism.Intern(*v)]; ok {
		return *found
	}

	panic("Verb " + *v.Ident + " not found")
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
		ElementType: prism.VectorType{Type: vec[0].Type()},
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
	case m.Ident != nil:
		return env.FetchVerb(m.Ident)
	case m.Alpha != nil:
		if f, ok := env.CurrentFunctionIR.(prism.DyadicFunction); ok {
			return prism.Alpha{
				TypeOf: f.AlphaType,
			}
		} else {
			panic("Alpha in monadic function")

		}
	case m.Omega != nil:
		if f, ok := env.CurrentFunctionIR.(prism.DyadicFunction); ok {
			return prism.Omega{
				TypeOf: f.OmegaType,
			}
		} else if f, ok := env.CurrentFunctionIR.(prism.MonadicFunction); ok {
			return prism.Omega{
				TypeOf: f.OmegaType,
			}
		} else {
			panic("Unreachable")

		}
	}

	panic("Other types not implemented")
}
