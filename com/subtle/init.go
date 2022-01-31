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
		if e.Monadic.Expression.Monadic != nil {
			if *e.Monadic.Expression.Monadic.Verb.Ident == "/" {
				return env.AnalyseDyadicOperator(e.Monadic)
			}
		}

		return env.AnalyseMonadic(e.Monadic)
	} else if e.Dyadic != nil {
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

func (env Environment) AnalyseMorphemes(ms *palisade.Morpheme) prism.Expression {
	mor := env.AnalyseMorpheme(ms)
	if vec, ok := mor.(prism.Vector); ok {
		if len(*vec.Body) == 1 {
			return (*vec.Body)[0]
		}
	}

	return mor
}

func (env Environment) AnalyseMorpheme(m *palisade.Morpheme) prism.Expression {
	switch {
	case m.Char != nil:
		vec := make([]prism.Expression, len(*m.Char))
		for i, c := range *m.Char {
			vec[i] = prism.Char{Value: string(c[0])}
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
		vec := make([]prism.Expression, len(*m.String))
		for i, c := range *m.String {
			vec[i] = prism.String{Value: c}
		}

		return prism.Vector{
			ElementType: prism.VectorType{Type: prism.StringType},
			Body:        &vec,
		}
	case m.Subexpr != nil:
		vec := make([]prism.Expression, len(*m.Subexpr))
		var typ prism.Type
		for i, c := range *m.Subexpr {
			vec[i] = env.AnalyseExpression(&c)
			if i == 0 {
				typ = vec[i].Type()
			} else {
				if typ != vec[i].Type() {
					panic("Type mismatch")
				}
			}
		}
		return prism.Vector{
			ElementType: prism.VectorType{Type: typ},
			Body:        &vec,
		}
	case m.Alpha != nil:
		if len(*m.Alpha) == 1 {
			if f, ok := env.CurrentFunctionIR.(prism.DyadicFunction); ok {
				return prism.Alpha{
					TypeOf: f.AlphaType,
				}
			}
		}

		panic("Unreachable")

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
			panic("Unreachable")

		}
	}

	panic("Other types not implemented")
}
