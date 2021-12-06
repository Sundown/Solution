package temporal

import (
	"sundown/solution/lexer"
	"sundown/solution/oversight"
)

type Function struct {
	Ident      *Ident
	TakesAlpha *Type
	TakesOmega *Type
	Gives      *Type
	Body       *Expression
	Special    bool
}

func (f *Function) String() string {
	var body, sig string
	if f.Body != nil {
		body = " =\n" + f.Body.String() + "\n"
	} else {
		body = ";\n\n"
	}

	if f.Ident.IsFoundational() {
		sig = *f.Ident.Ident
	} else {
		sig = *f.Ident.Namespace + "::" + *f.Ident.Ident
	}

	return sig + " : " + f.TakesAlpha.String() + ", " + f.TakesOmega.String() + " -> " + f.Gives.String() + body
}

// Name to be used within LLVM IR for ease of reading
func (i *Function) ToLLVMName() string {
	return *i.Ident.Namespace + "::" + *i.Ident.Ident + " " + i.TakesAlpha.String() + ", " + /*  i.TakesOmega.String() + */ "->" + i.Gives.String()
}

// Essentially declaration string
func (f *Function) SigString() string {
	return *f.Ident.Namespace + "::" + *f.Ident.Ident + " : " +
		f.TakesAlpha.String() + ", " + f.TakesOmega.String() + " -> " + f.Gives.String()
}

func (state *State) AnalyseFunction(function *lexer.Ident) (f *Function) {
	f = state.GetFunction(IRIdent(function))

	if f == nil {
		oversight.Error("ident \"" + *function.Ident + "\" is undefined in scope and Foundation.\n" + function.Pos.String()).Exit()
	}

	return f
}

func (state *State) GetFunction(key *Ident) *Function {
	if key.Namespace == nil {
		noun := state.Functions[IdentKey{Namespace: "_", Ident: *key.Ident}]
		if noun == nil {
			return state.Functions[IdentKey{Namespace: *state.PackageIdent, Ident: *key.Ident}]
		} else {
			return noun
		}
	} else {
		if key.Ident == nil {
			return nil
		}

		return state.Functions[IdentKey{Namespace: *key.Namespace, Ident: *key.Ident}]
	}
}

func (state *State) AnalyseFnDef(statement *lexer.FnDef) {
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

	state.CurrentFunction = nil
}

func (state *State) AnalyseFnDecl(statement *lexer.FnSig) {
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
			TakesAlpha: state.AnalyseType(statement.TakesAlpha),
			TakesOmega: state.AnalyseType(statement.TakesOmega),
			Gives:      state.AnalyseType(statement.Gives),
			Special:    false,
		}
	} else {
		oversight.Error(*statement.Ident + " is already declared as " + state.Functions[key].SigString() + ".\n" + statement.Pos.String()).Exit()
	}
}
