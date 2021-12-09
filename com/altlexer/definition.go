package altlexer

import "github.com/alecthomas/participle/v2/lexer"

type State struct {
	Pos        lexer.Position
	Statements []*struct {
		Directive *Directive `"@" @@`
		TypeDecl  *TypeDecl  `| @@`
		NounDecl  *NounDecl  `| @@`
		FnSig     *FnSig     `| @@`
		FnDef     *FnDef     `| @@`
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
	Ident *string `@Ident "~"`
	Type  *Type   `@@`
}

type NounDecl struct {
	Pos   lexer.Position
	Ident *string   `@Ident "="`
	Value *Morpheme `@@ ";"`
}

type FnSig struct {
	Pos        lexer.Position
	Ident      *string `@Ident ":"":"`
	TakesAlpha *Type   `@@ ","`
	TakesOmega *Type   `@@ "-"`
	Gives      *Type   `">" @@`
}

type FnDef struct {
	Pos         lexer.Position
	Ident       *string       `@Ident "="`
	Expressions []*Expression `(@@ ";")+`
}

type Expression struct {
	Pos         lexer.Position
	Application *Application `( @@`
	Morpheme    *Morpheme    `| @@ )`
}

type Type struct {
	Pos       lexer.Position
	Primative *Ident  ` @@`
	Vector    *Type   `| "[" @@ "]"`
	Tuple     []*Type `| "(" (@@ ("," @@)*)? ")"`
}

type Application struct {
	Pos            lexer.Position
	Function       *Ident      `@@`
	ParameterAlpha *Expression `@@ ","`
	ParameterOmega *Expression `@@`
}

type Morpheme struct {
	Pos        lexer.Position
	Tuple      []*Expression `	"(" (@@ ("," @@)*)? ")"`
	Vec        []*Expression `| "[" (@@ ("," @@)*)? "]"`
	Int        *int64        `| @('-'? Int)`
	Real       *float64      `| @('-'? Float)`
	Nil        *string       `| @"Nil"`
	String     *string       `| @String`
	Char       *string       `| @Char`
	ParamAlpha *string       `| @"Alpha"`
	ParamOmega *string       `| @"Omega"`
	Noun       *Ident        `| @@`
}
