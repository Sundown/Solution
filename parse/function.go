package parse

import (
	"sundown/sunday/lex"
)

type Function struct {
	Ident   *Ident
	Takes   *Type
	Gives   *Type
	Body    *Expression
	Special bool
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

	return sig + " : " + f.Takes.String() + " -> " + f.Gives.String() + body
}

// Name to be used within LLVM IR for ease of reading
func (i *Function) ToLLVMName() string {
	return *i.Ident.Namespace + "::" + *i.Ident.Ident + " " + i.Takes.String() + "->" + i.Gives.String()
}

// Essentially declaration string
func (f *Function) SigString() string {
	return *f.Ident.Namespace + "::" + *f.Ident.Ident + " : " +
		f.Takes.String() + " -> " + f.Gives.String()
}

func (state *State) AnalyseFunction(function *lex.Ident) (f *Function) {
	f = state.GetFunction(IRIdent(function))

	if f == nil {
		panic(*function.Ident + " not found in " + *state.PackageIdent + " or Foundation")
	}

	return f
}

func (state *State) AnalyseBlock(block []*lex.Expression) (b *Expression) {
	var body []*Expression
	for index, elm := range block {
		body[index] = state.AnalyseExpression(elm)
	}

	// TODO: need some way to calculate typeof
	b = &Expression{Block: body}
	return b
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
		return state.Functions[IdentKey{Namespace: *key.Namespace, Ident: *key.Ident}]
	}
}
