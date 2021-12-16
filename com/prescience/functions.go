package prescience

import (
	"sundown/solution/oversight"
	"sundown/solution/palisade"
	"sundown/solution/prism"
)

func Init(pali *palisade.PalisadeResult) (env *prism.Environment) {
	if pali == nil {
		oversight.Panic("Palisade state is nil")
	}

	for _, v := range pali.Statements {
		if v.FnDef != nil {
			InvokeFunctionDeclaration(v.FnDef)
		}
	}

	return
}

func InvokeFunctionDeclaration(fd *palisade.FnDef) {
	f := prism.Function{
		Name: fd.Ident.Intern(),

	}
	return
}
