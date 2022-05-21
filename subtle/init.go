package subtle

import (
	"github.com/sundown/solution/palisade"
	"github.com/sundown/solution/prism"
)

// Environment inherits from prism.Environment
type Environment struct {
	*prism.Environment
}

// Parse is entry point of Subtle,
// receiving environment created by Prism where Palisade stage is complete
func Parse(penv *prism.Environment) *prism.Environment {
	env := Environment{penv}

	tempEntry := ""

	for _, stmt := range env.LexResult.Environmentments {
		// First pass, declare all functions so that they
		// may be referenced before they are defined (text-wise)
		// Handle compiler directives
		if d := stmt.Directive; d != nil {
			switch *d.Command {
			case "Package":
				env.Output = *d.Value
			case "Entry":
				tempEntry = *d.Value
			}
		}

		if f := stmt.Function; f != nil {
			// palisade.Function is agnostic to arity
			// containing either monadic or dyadic
			env.internFunction(*f)
		}
	}

	for _, f := range env.DyadicFunctions {
		env.analyseDBody(f)
	}

	for _, f := range env.MonadicFunctions {
		env.analyseMBody(f)
	}

	if fn, ok := env.MonadicFunctions[prism.Ident{Package: "_", Name: tempEntry}]; ok {
		env.EntryFunction = *fn
	} else if fn, ok := env.MonadicFunctions[prism.Ident{Package: env.Output, Name: "Main"}]; ok {
		env.EntryFunction = *fn
	} else {
		prism.Panic("No entry function found")
	}

	return env.Environment
}

func (env Environment) internFunction(f palisade.Function) {
	var alpha, omega, sigma prism.Type
	dyadic := false
	var ident prism.Ident
	var body []palisade.Expression

	if p := f.TypedFunction; p != nil {
		if p.Dyadic != nil {
			dyadic = true
			ident = env.AwareIntern(*p.Dyadic.Ident)
			alpha = env.SubstantiateType(*p.Dyadic.Alpha)
			omega = env.SubstantiateType(*p.Dyadic.Omega)
		} else if p.Monadic != nil {
			ident = env.AwareIntern(*p.Monadic.Ident)
			omega = env.SubstantiateType(*p.Monadic.Omega)
		}

		sigma = env.SubstantiateType(*p.Returns)
	} else if p := f.AmbiguousFunction; p != nil {
		dyadic = env.determineArity(&f)
		alpha = prism.Universal{}
		omega = prism.Universal{}
		sigma = prism.Universal{}
		ident = env.AwareIntern(*p.Ident)
	}

	if f.Tacit != nil {
		body = []palisade.Expression{*f.Tacit}
	} else if f.Body != nil {
		body = *f.Body
	}

	if dyadic {
		fn := prism.DyadicFunction{
			Name:      ident,
			AlphaType: alpha,
			OmegaType: omega,
			Returns:   sigma,
			PreBody:   &body,
		}

		if _, ok := env.DyadicFunctions[fn.Name]; ok {
			prism.Panic("Dyadic function " + fn.Name.String() + " is already defined")
		} else {
			env.DyadicFunctions[fn.Name] = &fn
		}
	} else {
		fn := prism.MonadicFunction{
			Name:      ident,
			OmegaType: omega,
			Returns:   sigma,
			PreBody:   &body,
		}

		if _, ok := env.MonadicFunctions[fn.Name]; ok {
			prism.Panic("Monadic function " + fn.Name.String() + " is already defined")
		} else {
			env.MonadicFunctions[fn.Name] = &fn
		}
	}
}

func (env Environment) analyseDBody(f *prism.DyadicFunction) {
	if _, okr := f.OmegaType.(prism.Universal); okr || f.Attrs().Special || f.Attrs().SkipBuilder {
		return
	}

	if _, okl := f.AlphaType.(prism.Universal); okl {
		return
	}

	t := env.CurrentFunctionIR
	env.CurrentFunctionIR = *f

	if len(f.Body) == 0 {
		for _, expr := range *f.PreBody {
			f.Body = append(f.Body, env.analyseExpression(&expr))
		}
	} else {
		panic("Body already filled somehow")
	}

	env.CurrentFunctionIR = t
}

func (env Environment) analyseMBody(f *prism.MonadicFunction) {
	if _, ok := f.OmegaType.(prism.Universal); ok || f.Attrs().Special || f.Attrs().SkipBuilder {
		return
	}

	t := env.CurrentFunctionIR
	env.CurrentFunctionIR = *f

	if len(f.Body) == 0 {
		for _, expr := range *f.PreBody {
			f.Body = append(f.Body, env.analyseExpression(&expr))
		}
	}

	env.CurrentFunctionIR = t
}

func (env Environment) analyseExpression(e *palisade.Expression) prism.Expression {
	if e.Monadic != nil {
		if e.Monadic.Expression == nil {
			return env.FetchMVerb(e.Monadic.Applicable.Verb)
		}
		return env.analyseMonadic(e.Monadic)
	} else if e.Dyadic != nil {
		return env.analyseDyadic(e.Dyadic)
	} else if e.Morphemes != nil {
		return env.analyseMorphemes(e.Morphemes)
	}

	prism.Panic("unreachable")
	panic(nil)
}
