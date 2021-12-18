package weave

import (
	"sundown/solution/palisade"
	"sundown/solution/prism"
)

func Init(env *prism.Environment, lex *palisade.PalisadeResult) *prism.Environment {
	for _, f := range lex.Statements {
		exprs := []prism.Expression{}
		for _, s := range f.FnDef.Expressions {
			s := SubState{
				Env:            env,
				mutablePos:     0,
				mutableSubPos:  0,
				Subexpressions: append(s.Lexemes, &palisade.Subexpression{Morpheme: nil, Sub: nil}),
			}

			exprs = append(exprs, s.HandleLexeme())
		}

		fn, ok := env.Functions[prism.Intern(*f.FnDef.Ident)]
		if ok {
			fn.PreBody = &exprs
		}
	}

	return env
}

func (s SubState) HandleLexeme() prism.Expression {
	if s.Cur().Morpheme != nil {
		if s.Cur().Morpheme.Ident != nil {
			id := prism.Intern(*s.Cur().Morpheme.Ident)
			if s.isApplicationIdent(&id) {
				return prism.Application{
					Operator: *s.Env.GetFunction(prism.Intern(*s.Cur().Morpheme.Ident)),
					Operand:  s.Forward().HandleLexeme()}
			} /*else if isVariable(s.Cur().Morpheme.Ident) {
				return Dangle{
					Outer: Ident{Value: *s.Cur().Morpheme.Ident},
					Inner: s.Forward().HandleLexeme()}
			}*/ // TODO: implement
		} else if s.Cur().Morpheme.Int != nil {
			return prism.Dangle{
				Outer: prism.Int{Value: *s.Cur().Morpheme.Int},
				Inner: s.Forward().HandleLexeme()}
		} else if s.Cur().Morpheme.String != nil {
			return prism.Dangle{
				Outer: prism.String{Value: *s.Cur().Morpheme.String},
				Inner: s.Forward().HandleLexeme()}
		} else if s.Cur().Morpheme.Real != nil {
			return prism.Dangle{
				Outer: prism.Real{Value: *s.Cur().Morpheme.Real},
				Inner: s.Forward().HandleLexeme(),
			}
		} else if s.Cur().Morpheme.Char != nil {
			return prism.Dangle{
				Outer: prism.Char{Value: *s.Cur().Morpheme.Char},
				Inner: s.Forward().HandleLexeme(),
			}
		} else if s.Cur().Morpheme.Alpha != nil {
			return prism.Dangle{
				Outer: prism.Alpha{},
				Inner: s.Forward().HandleLexeme(),
			}
		} else if s.Cur().Morpheme.Omega != nil {
			return prism.Dangle{
				Outer: prism.Omega{},
				Inner: s.Forward().HandleLexeme(),
			}
		}
	} else if s.Cur().Sub != nil {
		ns := SubState{
			mutablePos:     0,
			mutableSubPos:  0,
			Subexpressions: append(s.Cur().Sub.Lexemes, &palisade.Subexpression{Morpheme: nil, Sub: nil}),
		}

		return prism.Dangle{Outer: prism.Subexpression{Expression: ns.HandleLexeme()}, Inner: s.Forward().HandleLexeme()}
	}

	return prism.EOF{}
}
