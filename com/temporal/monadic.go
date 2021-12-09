package temporal

import "sundown/solution/lexer"

type MonadicApplication struct {
	TypeOf *Type
}

func (state *State) AnalyseMonadicApplication(monadic *lexer.Monadic) (m *MonadicApplication) {
	if monadic.Morpheme != nil {
		morph := state.AnalyseMorpheme(monadic.Morpheme)
		
	}
}
