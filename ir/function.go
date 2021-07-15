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
	var ident string
	if f.Ident.Namespace != nil {
		ident = *f.Ident.Namespace + "::" + *f.Ident.Ident
	} else {
		ident = *f.Ident.Ident
	}

	return ident + " : " + f.Takes.String() + " -> " + f.Gives.String() + " =\n" + f.Body.String() + "\n"
}

func AnalyseFunction(function *parser.Ident) (f *Function) {
	f = &Function{
		Ident: &Ident{
			Namespace: function.Namespace,
			Ident:     function.Ident,
		},
		Takes: &Type{Atomic: "Int"},
		Gives: &Type{Atomic: "Int"}}
	return f
}

func AnalyseBlock(block []*parser.Expression) (b *Expression) {
	var body []*Expression
	for index, elm := range block {
		body[index] = AnalyseExpression(elm)
	}

	/* TODO: need some way to calculate typeof */
	b = &Expression{Block: body}
	return b
}
