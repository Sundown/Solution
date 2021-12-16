package weave

import (
	"sundown/solution/palisade"
	"sundown/solution/prism"
)

func (state *SubState) isApplicationIdent(i *prism.Ident) bool {
	_, ok := state.Env.Functions[*i]
	return ok
}

func (s SubState) Next() *palisade.Subexpression {
	s.mutablePos++
	return s.Subexpressions[s.mutablePos]
}

func (s SubState) Forward() SubState {
	s.mutablePos++
	return s
}

func (s SubState) Peek() *palisade.Subexpression {
	return s.Subexpressions[s.mutablePos+1]
}

func (s SubState) Cur() *palisade.Subexpression {
	return s.Subexpressions[s.mutablePos]
}
