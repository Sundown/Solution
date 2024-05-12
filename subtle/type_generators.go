package subtle

import (
	"github.com/sundown/solution/palisade"
	"github.com/sundown/solution/prism"
)

func (env *Environment) generateMonadicTypes(app palisade.Applicable, rType prism.Type, function *prism.Function) {
	f := env.FetchMVerb(app.Verb)
	if g, ok := f.OmegaType.(prism.Universal); ok && g.Has(rType) {
		// For each application of an algebraic function, a new function with concrete
		// types is created based on the type of its operand
		f.Name.Name = rType.String() + "." + f.Name.Name
		f.OmegaType = rType

		// Fill out the body, input type is already known
		env.analyseMBody(&f)

		// Iterate through expressions inside the function body until a 'return' is found,
		// determine the return type from the type of the operand
		if f.Returns.IsAlgebraic() {
			ret := rType
			for _, expr := range f.Body {
				if res, ok := expr.(prism.MonadicApplication); ok {
					if !isReturn(res.Operator) {
						continue
					}

					ret = res.Operand.Type()
					break
				}
			}
			f.Returns = f.Returns.Resolve(ret)
		}

	}

	// Induct new instance of the function, will have a unique name based on input type
	env.MonadicFunctions[f.Name] = &f
	*function = f
}

func (env *Environment) generateDyadicTypes(app *palisade.Applicable, rType prism.Type, lType prism.Type, function *prism.Function) {
	didGeneric := false
	f := env.FetchDVerb(app.Verb)
	// For each application of an algebraic function, a new function with concrete
	// types is created based on the types of its operands
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
		// Fill out the body, input types are already known
		env.analyseDBody(&f)

		// Iterate through expressions inside the function body until a 'return' is found,
		// determine the return type from the type of the operand
		if f.Returns.IsAlgebraic() {
			ret := rType
			for _, expr := range f.Body {
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

	// Induct new instance of the function, will have a unique name based on input types
	env.DyadicFunctions[f.Name] = &f
	*function = f
}

func isReturn(m prism.MonadicFunction) bool {
	return m.Name.Package == "_" && m.Name.Name == "←"
}
