package ir

import (
	"sundown/sunday/parser"
)

type Function struct {
	Ident *Ident
	Takes *Type
	Gives *Type
	Body  *Expression
}

func (f *Function) String() string {
	var body string
	if f.Body != nil {
		body = " =\n" + f.Body.String() + "\n"
	} else {
		body = ";\n"
	}

	return f.Ident.Namespace + "::" + f.Ident.Ident + " : " +
		f.Takes.String() + " -> " + f.Gives.String() + body
}

func (f *Function) SigString() string {
	return f.Ident.Namespace + "::" + f.Ident.Ident + " : " +
		f.Takes.String() + " -> " + f.Gives.String()
}

func (state *State) AnalyseFunction(function *parser.Ident) (f *Function) {
	var namespace string
	if function.Namespace == nil {
		namespace = *state.PackageIdent
	} else {
		namespace = *function.Namespace
	}

	// Attempt 1: is function in local package?
	f = state.GetFunction(&namespace, function.Ident)
	if f == nil {
		str := "_"
		// Attempt 2: is function part of foundation?
		f = state.GetFunction(&str, function.Ident)
	}

	if f == nil {
		panic(namespace + "::" + *function.Ident + " not found in " + *state.PackageIdent + " or Foundation")
	}

	return f
}

func (state *State) AnalyseBlock(block []*parser.Expression) (b *Expression) {
	var body []*Expression
	for index, elm := range block {
		body[index] = state.AnalyseExpression(elm)
	}

	/* TODO: need some way to calculate typeof */
	b = &Expression{Block: body}
	return b
}
