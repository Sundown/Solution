package weave

import "sundown/solution/palisade"

func isApplicationIdent(s *string) bool {
	for _, v := range unaries {
		if v == *s {
			return true
		}
	}
	return false
}

func isVariable(s *string) bool {
	for _, v := range variables {
		if v == *s {
			return true
		}
	}
	return false
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
