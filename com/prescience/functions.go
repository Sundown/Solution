package prescience

import (
	"sundown/solution/oversight"
	"sundown/solution/palisade"
	"sundown/solution/prism"
)

func Init(env *prism.Environment, pali *palisade.PalisadeResult) *prism.Environment {
	if pali == nil {
		oversight.Panic("Palisade state is nil")
	}

	for _, v := range pali.Statements {
		if v.FnDef != nil {
			InvokeFunctionDeclaration(v.FnDef, env)
		}
	}

	return env
}

func InvokeFunctionDeclaration(fd *palisade.FnDef, env *prism.Environment) {
	env.Functions[prism.Intern(*fd.Ident)] = prism.Function{
		Name:      prism.Intern(*fd.Ident),
		AlphaType: env.SubstantiateType(*fd.TakesAlpha),
		OmegaType: env.SubstantiateType(*fd.TakesOmega),
		Returns:   env.SubstantiateType(*fd.Gives),
		Body:      nil, // for now...
	}
}
