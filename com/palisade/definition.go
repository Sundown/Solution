package palisade

import "github.com/alecthomas/participle/v2/lexer"

type PalisadeResult struct {
	Pos        lexer.Position
	Statements []*struct {
		//	Directive *Directive `"@" @@`
		TypeDecl *TypeDecl `@@`
		//NounDecl  *NounDecl  `| @@`
		FnDef *FnDef `| @@`
	} `@@*`
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
	Pos       lexer.Position
	Namespace *string `(@Ident ":" ":")?`
	Ident     *string `@Ident`
}

type TypeDecl struct {
	Pos   lexer.Position
	Ident *Ident `@@ "~"`
	Type  *Type  `@@`
}

/* type NounDecl struct {
	Pos   lexer.Position
	Ident *string  `@Ident "="`
	Value *Morpheme `@@ ";"`
} */

type FnDef struct {
	Pos         lexer.Position
	TakesAlpha  *Type         `"Δ" @@`
	Ident       *Ident        `@@`
	TakesOmega  *Type         `@@ "→"`
	Gives       *Type         `@@ ":"`
	Expressions []*Expression `(@@ ";")+ "∇"`
}

type Type struct {
	Pos       lexer.Position
	Primative *Ident  ` @@`
	Vector    *Type   `| "[" @@ "]"`
	Tuple     []*Type `| "(" (@@ ("," @@)*)? ")"`
}
