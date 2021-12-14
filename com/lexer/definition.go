package lexer

import "github.com/alecthomas/participle/v2/lexer"

type State struct {
	Pos         lexer.Position
	Expressions []*Expression `(@@ ";")+`
	/*Statements []*struct {
		Directive *Directive `"@" @@`
		TypeDecl  *TypeDecl  `| @@`
		//NounDecl  *NounDecl  `| @@`
		FnSig *FnSig `| @@`
		FnDef *FnDef `| @@`
	} `@@*`*/
}

type Directive struct {
	Pos   lexer.Position
	Class *string `@Ident`
	Instr *struct {
		Ident  *string  `( @Ident`
		String *string  `| @String`
		Number *float64 `| @Float)`
	} `@@`
}

type Ident struct {
	//Pos       lexer.Position
	Namespace *string `(@Ident ":" ":")?`
	Ident     *string `@Ident`
}

type TypeDecl struct {
	Pos   lexer.Position
	Ident *string `@Ident "~"`
	Type  *Type   `@@`
}

/* type NounDecl struct {
	Pos   lexer.Position
	Ident *string  `@Ident "="`
	Value *Morpheme `@@ ";"`
} */
// Δ∇
type FnSig struct {
	Pos         lexer.Position
	TakesAlpha  *Type         `"Δ" @@`
	Ident       *Ident        `@Ident `
	TakesOmega  *Type         `@@ "->"`
	Gives       *Type         `@@ + ":"`
	Expressions []*Expression `(@@ ";")+ "∇"`
}

type FnDef struct {
	Pos         lexer.Position
	Ident       *string       `@Ident "="`
	Expressions []*Expression `(@@ ";")+`
}

type Type struct {
	Pos       lexer.Position
	Primative *Ident  ` @@`
	Vector    *Type   `| "[" @@ "]"`
	Tuple     []*Type `| "(" (@@ ("," @@)*)? ")"`
}
