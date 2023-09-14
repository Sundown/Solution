package subtle

import (
	"github.com/sundown/solution/palisade"
	"github.com/sundown/solution/prism"
)

func isReturn(m prism.MonadicFunction) bool {
	return m.Name.Package == "_" && m.Name.Name == "←"
}

// Everything wrong with the compiler starts in this function
func (env *Environment) analysePrimeApplicable(app palisade.Applicable, lType, rType prism.Type) prism.Function {
	var function prism.Function
	if app.Verb != nil {
		if lType == nil {
			f := env.FetchMVerb(app.Verb)
			if g, ok := f.OmegaType.(prism.Universal); ok && g.Has(rType) {
				f.Name.Name = rType.String() + "." + f.Name.Name
				f.OmegaType = rType
				env.analyseMBody(&f)

				if f.Returns.IsAlgebraic() {
					ret := rType
					for _, expr := range f.Body {
						// Iter until we find a return (←) statement
						if res, ok := expr.(prism.MonadicApplication); ok {
							if isReturn(res.Operator) {
								ret = res.Operand.Type()
								break
							}
						}
					}
					f.Returns = f.Returns.Resolve(ret)
				}

			}

			env.MonadicFunctions[f.Name] = &f
			function = f
		} else {
			didGeneric := false
			f := env.FetchDVerb(app.Verb)
			if g, ok := f.OmegaType.(prism.Universal); ok && g.Has(rType) {
				f.Name.Name = rType.String() + "." + f.Name.Name
				f.OmegaType = rType

				didGeneric = true
			}

			if g, ok := f.AlphaType.(prism.Universal); ok && g.Has(rType) {
				f.Name.Name = lType.String() + "." + f.Name.Name
				f.AlphaType = lType

				didGeneric = true
			}

			if didGeneric {
				env.analyseDBody(&f)
				if f.Returns.IsAlgebraic() {
					ret := rType
					for _, expr := range f.Body {
						// Iter until we find a return (←) statement
						if res, ok := expr.(prism.MonadicApplication); ok {
							if isReturn(res.Operator) {
								ret = res.Operand.Type()
								break
							}
						}
					}
					f.Returns = f.Returns.Resolve(ret)
				}
			}

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
