package ir

import "sundown/sunday/parser"

func (state *State) AnalyseStatement(statement *parser.FnDecl) (s *Function) {
	takes, gives := AnalyseType(statement.Takes), AnalyseType(statement.Gives)
	e := Expression{TypeOf: gives}
	for _, expr := range statement.Expressions {
		e.Block = append(e.Block, state.AnalyseExpression(expr))
	}

	s = &Function{
		Ident: &Ident{
			Namespace: *state.PackageIdent,
			Ident:     *statement.Ident,
		},
		Takes: takes,
		Gives: gives,
		Body:  &e,
	}

	return s
}
