package parse

import "sundown/solution/lex"

func (state *State) AnalyseFnDef(statement *lex.FnDecl) {
	decl := state.Functions[IdentKey{
		Namespace: *state.PackageIdent,
		Ident:     *statement.Ident,
	}]

	// Somehow function hasn't been declared, or it already has a body
	if decl == nil || decl.Body != nil {
		panic("Not sure how you got here...")
	}

	e := Expression{TypeOf: decl.Gives}

	decl.Body = &e

	state.CurrentFunction = decl

	for _, expr := range statement.Expressions {
		e.Block = append(e.Block, state.AnalyseExpression(expr))
	}
}

func (state *State) AnalyseFnDecl(statement *lex.FnDecl) {
	// Key is used for existential verification and/or definition
	key := IdentKey{
		Namespace: *state.PackageIdent,
		Ident:     *statement.Ident,
	}

	if state.Functions[key] == nil {
		state.Functions[key] = &Function{
			Ident: &Ident{
				Namespace: state.PackageIdent,
				Ident:     statement.Ident,
			},
			Takes:   state.AnalyseType(statement.Takes),
			Gives:   state.AnalyseType(statement.Gives),
			Special: false,
		}
	} else {
		panic(*statement.Ident + " is already declared")
	}
}
