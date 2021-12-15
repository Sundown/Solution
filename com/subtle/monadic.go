package subtle

import "sundown/solution/palisade"

type MonadicApplication struct {
	TypeOf *Type
}

func (state *State) AnalyseMonadicApplication(monadic *palisade.Monadic) (m *MonadicApplication) {
	if monadic.Morpheme != nil {
		morph := state.AnalyseMorpheme(monadic.Morpheme)

	}
}
