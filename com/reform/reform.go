package reform

import (
	"fmt"
	"sundown/solution/lexer"
)

var unaries = []string{"u"}
var variables = []string{"v"}

type State struct{}

type Expression interface {
	String() string
}

type Application struct {
	Operator string
	Operand  Expression
}

type Dangle struct {
	Outer Expression
	Inner Expression
}

type Int struct {
	Value int64
}

type Ident struct {
	Value string
}

type EOF struct{}

type SubState struct {
	mutablePos     int
	mutableSubPos  int
	Subexpressions []*lexer.Subexpression
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

func (state *State) Init(lex *lexer.State) (expr Expression) {
	for _, s := range lex.Expressions {
		s := SubState{
			mutablePos:     0,
			mutableSubPos:  0,
			Subexpressions: append(s.Singletons, &lexer.Subexpression{nil, nil}),
		}

		expr = s.HandleSingleton()
	}

	fmt.Println(expr.String())
	return
}

func (s SubState) HandleSingleton() Expression {
	if s.Cur().Morpheme != nil {
		if s.Cur().Morpheme.Ident != nil {
			if isApplicationIdent(s.Cur().Morpheme.Ident) {
				return Application{
					Operator: *s.Cur().Morpheme.Ident,
					Operand:  s.Forward().HandleSingleton()}
			} else if isVariable(s.Cur().Morpheme.Ident) {
				return Dangle{
					Outer: Ident{Value: *s.Cur().Morpheme.Ident},
					Inner: s.Forward().HandleSingleton()}
			}
		} else if s.Cur().Morpheme.Int != nil {
			return Dangle{
				Outer: Int{Value: *s.Cur().Morpheme.Int},
				Inner: s.Forward().HandleSingleton()}
		}
	}

	return EOF{}
}

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
