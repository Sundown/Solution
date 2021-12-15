package prescience

import (
	"sundown/solution/oversight"
	"sundown/solution/palisade"
	"sundown/solution/weave"
)

func Init(pali *palisade.State) (ws *weave.State) {
	if pali == nil {
		oversight.Panic("Palisade state is nil")
	}

	for _, v := range pali.Statements {
		if v.FnDef != nil {
			InvokeFunctionDefinition(v.FnDef)
		}
	}

	return
}

func InvokeFunctionDefinition(fd *palisade.FnDef) {
	return
}
