package reform

import "sundown/solution/lexer"

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

func (s SubState) Next() *lexer.Subexpression {
	s.mutablePos++
	return s.Subexpressions[s.mutablePos]
}

func (s SubState) Forward() SubState {
	s.mutablePos++
	return s
}

func (s SubState) Peek() *lexer.Subexpression {
	return s.Subexpressions[s.mutablePos+1]
}

func (s SubState) Cur() *lexer.Subexpression {
	return s.Subexpressions[s.mutablePos]
}
