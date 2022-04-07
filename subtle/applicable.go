package subtle

import (
	"github.com/sundown/solution/palisade"
	"github.com/sundown/solution/prism"
)

/* type Applicable struct {

    Subexpr  *Expression `parser:"(('(' @@ ')')"`
    Verb     *Ident      `parser:"| @@)"`
    Operator *Operator   `parser:"@@?"`
} */

func (env *Environment) analysePrimeApplicable(app palisade.Applicable, lType, rType prism.Type) prism.Function {
	var function prism.Function
	if app.Verb != nil {
		if lType == nil {
			function = env.FetchMVerb(app.Verb)
		} else {
			function = env.FetchDVerb(app.Verb)
		}
	} else if app.Subexpr != nil {
		// Monadic/dyadic cases are handled within train system

		function = env.boardTrain(app.Subexpr, lType, rType)
	}
	return function
}

func (env *Environment) analyseApplicable(app palisade.Applicable, lType, rType prism.Type) prism.Function {
	var function prism.Function
	if app.Operator != nil {
		if lType == nil {
			function = env.monadicOperatorToFunction(env.analyseMonadicOperator(app, rType))
		}
		// TODO implement dyadic operators
		/* else {
			function = env.analyseDyadicOperator(app.Operator, function)
		} */
	} else {
		function = env.analysePrimeApplicable(app, lType, rType)
	}

	return function
}

// Inductive cast cardinality mismatch within algebraic group
