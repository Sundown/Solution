package prescience

import (
	"sundown/solution/oversight"
	"sundown/solution/palisade"
	"sundown/solution/prism"
)

// TODO
// 1. intenr type defs
// 2. intern atom defs
// Checking always that they're not reserved
// so future stages have an easier time

func Init(env *prism.Environment, pali *palisade.PalisadeResult) *prism.Environment {
	if pali == nil {
		oversight.Panic("Palisade state is nil")
	}

	// nouns and type decls later

	for _, v := range pali.Statements {
		if v.FnDef != nil {
			InvokeFunctionDeclaration(v.FnDef, env)
		}
	}

	return env
}

func InvokeFunctionDeclaration(fd *palisade.FnDef, env *prism.Environment) {
	// TODO: doesn't handle unaries
	env.Functions[prism.Intern(*fd.Ident)] = &prism.Function{
		Name:      prism.Intern(*fd.Ident),
		AlphaType: env.SubstantiateType(*fd.TakesAlpha),
		OmegaType: env.SubstantiateType(*fd.TakesOmega),
		Returns:   env.SubstantiateType(*fd.Gives),
		PreBody:   nil,
		Body:      nil,
	}
}
