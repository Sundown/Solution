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
