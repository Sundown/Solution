package subtle

import (
	"github.com/sundown/solution/prism"
)

func isReturn(m prism.MonadicFunction) bool {
	return m.Name.Package == "_" && m.Name.Name == "←"
}

func (env *Environment) generateMonadicTypes(f prism.MonadicFunction, rType prism.Type) prism.MonadicFunction {
	if g, ok := f.OmegaType.(prism.Universal); ok && g.Has(rType) {
		f.Name.Name = rType.String() + "." + f.Name.Name
		f.OmegaType = rType
		env.analyseMBody(&f)

		if f.Returns.IsAlgebraic() {
			ret := rType
			// Iter until we find a return (←) statement
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

	return f
}
