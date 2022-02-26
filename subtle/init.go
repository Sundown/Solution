package subtle

import (
	"github.com/sundown/solution/palisade"
	"github.com/sundown/solution/prism"
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
		env.analyseDBody(f)
	}

	for _, f := range env.MonadicFunctions {
		env.analyseMBody(f)
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
		if _, ok := env.DyadicFunctions[fn.Name]; ok {
			prism.Panic("Dyadic function " + fn.Name.String() + " is already defined")
		} else {
			env.DyadicFunctions[fn.Name] = &fn
		}
	} else if f.Monadic != nil {
		fn := prism.MonadicFunction{
			Name:      prism.Intern(*f.Monadic.Ident),
			OmegaType: env.SubstantiateType(*f.Monadic.Omega),
			Returns:   env.SubstantiateType(*f.Returns),
			PreBody:   f.Body,
		}
		if _, ok := env.MonadicFunctions[fn.Name]; ok {
			prism.Panic("Monadic function " + fn.Name.String() + " is already defined")
		} else {
			env.MonadicFunctions[fn.Name] = &fn
		}
	}
}

func (env Environment) analyseDBody(f *prism.DyadicFunction) {
	if f.Special || f.SkipBuilder {
		return
	}

	env.CurrentFunctionIR = *f

	for _, expr := range *f.PreBody {
		f.Body = append(f.Body, env.analyseExpression(&expr))
	}
}

func (env Environment) analyseMBody(f *prism.MonadicFunction) {
	if f.Special || f.SkipBuilder {
		return
	}

	env.CurrentFunctionIR = *f

	if len(f.Body) == 0 {
		for _, expr := range *f.PreBody {
			f.Body = append(f.Body, env.analyseExpression(&expr))
		}
	}
}

func (env Environment) analyseExpression(e *palisade.Expression) prism.Expression {
	if e.Monadic != nil {
		if e.Monadic.Expression == nil {
			return env.FetchMVerb(e.Monadic.Verb)
		}
		return env.analyseMonadic(e.Monadic)
	} else if e.Dyadic != nil {
		if e.Dyadic.Expression == nil {
			return env.analysePartial(e.Dyadic)
		}
		return env.analyseDyadic(e.Dyadic)
	} else if e.Morphemes != nil {
		return env.analyseMorphemes(e.Morphemes)
	} else if e.Operator != nil {
		return env.analyseDyadicOperator(e.Operator)
	}

	panic("unreachable")
}
