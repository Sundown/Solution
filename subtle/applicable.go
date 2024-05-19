package subtle

import (
	"github.com/sundown/solution/palisade"
	"github.com/sundown/solution/prism"
)

// Everything wrong with the compiler starts in this function
func (env *Environment) analysePrimeApplicable(app palisade.Applicable, lType, rType prism.Type) prism.Function {
	var function prism.Function
	if app.Verb != nil {
		if lType == nil {
			f := env.generateMonadicTypes(env.FetchMVerb(app.Verb), rType)
			env.MonadicFunctions[f.Name] = &f
			function = f
		} else {
			f := env.generateDyadicTypes(env.FetchDVerb(app.Verb), lType, rType)

			env.DyadicFunctions[f.Name] = &f
			function = f
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
		} else {
			panic("Dyadic operators not implemented")
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
