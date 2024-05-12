package subtle

import (
	"github.com/sundown/solution/palisade"
	"github.com/sundown/solution/prism"
)

func isReturn(m prism.MonadicFunction) bool {
	return m.Name.Package == "_" && m.Name.Name == "←"
}

func (env *Environment) analysePrimeApplicable(app palisade.Applicable, lType, rType prism.Type) prism.Function {
	var function prism.Function
	if app.Verb != nil {
		if lType == nil {
			// Iter until we find a return (←) statement
			env.generateMonadicTypes(app, rType, &function)
		} else {
			// Iter until we find a return (←) statement
			env.generateDyadicTypes(&app, rType, lType, &function)
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
