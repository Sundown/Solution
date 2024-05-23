package subtle

import (
	"github.com/sundown/solution/prism"
)

func isReturn(m prism.MonadicFunction) bool {
	return m.Name.Package == "_" && m.Name.Name == "←"
}

func determineFunctionReturn(f prism.Function, def prism.Type) prism.Type {
	var body *[]prism.Expression

	if m, ok := f.(prism.MonadicFunction); ok {
		body = &m.Body
	} else if d, ok := f.(prism.DyadicFunction); ok {
		body = &d.Body
	}

	if f.Type().IsAlgebraic() {
		// Iter until we find a return (←) statement
		for _, expr := range *body {
			if res, ok := expr.(prism.MonadicApplication); ok {
				if res.Operator.Name.Package == "_" && res.Operator.Name.Name == "←" {
					return f.Type().Resolve(res.Operand.Type())
				}
			}
		}

		return def

	}

	return f.Type()
}

func (env *Environment) generateMonadicTypes(f prism.MonadicFunction, rType prism.Type) prism.MonadicFunction {
	if g, ok := f.OmegaType.(prism.Group); ok && g.Has(rType) {
		f.Name.Name = rType.String() + "." + f.Name.Name
		f.OmegaType = rType

		env.analyseMBody(&f)

		f.Returns = determineFunctionReturn(f, rType)
	}

	// TODO (elsewhere) return shouldn't be allowed to be algebraic if alpha and omega are concrete

	return f
}

func (env *Environment) generateDyadicTypes(f prism.DyadicFunction, lType, rType prism.Type) prism.DyadicFunction {
	didGeneric := false
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

	env.analyseDBody(&f)

	if didGeneric {
		f.Returns = determineFunctionReturn(f, rType)
	}

	return f
}
