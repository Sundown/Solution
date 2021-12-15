package weave

import (
	"fmt"
	"sundown/solution/palisade"
)

// Dummy environment for testing
var unaries = []string{"u"}
var variables = []string{"v"}

func (state *State) Init(lex *palisade.State) (expr Expression) {
	for _, f := range lex.Statements {
		for _, s := range f.FnDef.Expressions {
			s := SubState{
				mutablePos:     0,
				mutableSubPos:  0,
				Subexpressions: append(s.Lexemes, &palisade.Subexpression{Morpheme: nil, Sub: nil}),
			}

			expr = s.HandleLexeme()
		}
	}

	fmt.Println(expr.String())
	return
}

func (s SubState) HandleLexeme() Expression {
	if s.Cur().Morpheme != nil {
		if s.Cur().Morpheme.Ident != nil {
			if isApplicationIdent(s.Cur().Morpheme.Ident) {
				return Application{
					Operator: *s.Cur().Morpheme.Ident,
					Operand:  s.Forward().HandleLexeme()}
			} else if isVariable(s.Cur().Morpheme.Ident) {
				return Dangle{
					Outer: Ident{Value: *s.Cur().Morpheme.Ident},
					Inner: s.Forward().HandleLexeme()}
			}
		} else if s.Cur().Morpheme.Int != nil {
			return Dangle{
				Outer: Int{Value: *s.Cur().Morpheme.Int},
				Inner: s.Forward().HandleLexeme()}
		} else if s.Cur().Morpheme.String != nil {
			return Dangle{
				Outer: String{Value: *s.Cur().Morpheme.String},
				Inner: s.Forward().HandleLexeme()}
		} else if s.Cur().Morpheme.Real != nil {
			return Dangle{
				Outer: Real{Value: *s.Cur().Morpheme.Real},
				Inner: s.Forward().HandleLexeme(),
			}
		} else if s.Cur().Morpheme.Char != nil {
			return Dangle{
				Outer: Char{Value: *s.Cur().Morpheme.Char},
				Inner: s.Forward().HandleLexeme(),
			}
		} else if s.Cur().Morpheme.Alpha != nil {
			return Dangle{
				Outer: Alpha{},
				Inner: s.Forward().HandleLexeme(),
			}
		} else if s.Cur().Morpheme.Omega != nil {
			return Dangle{
				Outer: Omega{},
				Inner: s.Forward().HandleLexeme(),
			}
		}
	} else if s.Cur().Sub != nil {
		ns := SubState{
			mutablePos:     0,
			mutableSubPos:  0,
			Subexpressions: append(s.Cur().Sub.Lexemes, &palisade.Subexpression{Morpheme: nil, Sub: nil}),
		}

		return Dangle{Outer: Subexpression{Expression: ns.HandleLexeme()}, Inner: s.Forward().HandleLexeme()}
	}

	return EOF{}
}
