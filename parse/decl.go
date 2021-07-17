package parse

import "sundown/sunday/lex"

func (state *State) AnalyseFnDef(statement *lex.FnDecl) {
	decl := state.Functions[IdentKey{Namespace: *state.PackageIdent, Ident: *statement.Ident}]
	if decl == nil || decl.Body != nil {
		panic("Not sure how you got here...")
	}

	takes, gives := state.AnalyseType(statement.Takes), state.AnalyseType(statement.Gives)
	e := Expression{TypeOf: gives}
	for _, expr := range statement.Expressions {
		e.Block = append(e.Block, state.AnalyseExpression(expr))
	}

	state.Functions[IdentKey{Namespace: *state.PackageIdent, Ident: *statement.Ident}] = &Function{
		Ident: &Ident{
			Namespace: state.PackageIdent,
			Ident:     statement.Ident,
		},
		Takes: takes,
		Gives: gives,
		Body:  &e,
	}

}

func (state *State) AnalyseFnDecl(statement *lex.FnDecl) {
	if state.Functions[IdentKey{Namespace: *state.PackageIdent, Ident: *statement.Ident}] == nil {
		state.Functions[IdentKey{Namespace: *state.PackageIdent, Ident: *statement.Ident}] = &Function{
			Ident: &Ident{
				Namespace: state.PackageIdent,
				Ident:     statement.Ident,
			},
			Takes: state.AnalyseType(statement.Takes) /* -> */, Gives: state.AnalyseType(statement.Gives),
		}
	} else {
		panic(*statement.Ident + " is already declared")
	}
}
