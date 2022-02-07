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

	// TODO fix dumb
	for _, stmt := range env.LexResult.Environmentments {
		if d := stmt.Directive; d != nil {
			if *d.Command == "Entry" {
				fn, ok := env.MonadicFunctions[prism.Ident{Package: "_", Name: *d.Value}]
				if !ok {
					panic("Entry function not found")
				}

				env.EntryFunction = *fn
			}
		}
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

	if len(f.Body) == 0 {
		for _, expr := range *f.PreBody {
			f.Body = append(f.Body, env.AnalyseExpression(&expr))
		}
	}
}

func (env Environment) AnalyseExpression(e *palisade.Expression) prism.Expression {
	if e.Monadic != nil {
		if e.Monadic.Expression.Monadic != nil {
			if e.Monadic.Expression.Monadic.Verb != nil &&
				(*e.Monadic.Expression.Monadic.Verb.Ident == "/" || *e.Monadic.Expression.Monadic.Verb.Ident == "Map") {
				return env.AnalyseDyadicOperator(e.Monadic)
			}
		}

		return env.AnalyseMonadic(e.Monadic)
	} else if e.Dyadic != nil {
		if e.Dyadic.Expression == nil {
			return env.AnalysePartial(e.Dyadic)
		}
		return env.AnalyseDyadic(e.Dyadic)
	} else if e.Morphemes != nil {
		return env.AnalyseMorphemes(e.Morphemes)
	}

	panic("unreachable")
}

func (env Environment) FetchDVerb(v *palisade.Ident) prism.Expression {
	if found, ok := env.DyadicFunctions[prism.Intern(*v)]; ok {
		return *found
	}

	panic("Verb " + *v.Ident + " not found")
}

func (env Environment) FetchMVerb(v *palisade.Ident) prism.Expression {
	if found, ok := env.MonadicFunctions[prism.Intern(*v)]; ok {
		return *found
	}
	panic("Verb " + *v.Ident + " not found")
}

func (env Environment) FetchVerb(v *palisade.Ident) prism.Expression {
	if found, ok := env.MonadicFunctions[prism.Intern(*v)]; ok {
		return *found
	} else if found, ok := env.DyadicFunctions[prism.Intern(*v)]; ok {
		return *found
	}

	panic("Verb " + *v.Ident + " not found")
}
