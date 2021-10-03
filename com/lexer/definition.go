package lexer

import "github.com/alecthomas/participle/v2/lexer"

type State struct {
	Pos        lexer.Position
	Statements []*struct {
		Directive *Directive `"@" @@`
		TypeDecl  *TypeDecl  `| @@`
		NounDecl  *NounDecl  `| @@`
		FnDecl    *FnDecl    `| @@`
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
	Ident *string  `@Ident "="`
	Value *Primary `@@ ";"`
}

type FnDecl struct {
	Pos         lexer.Position
	Ident       *string       `@Ident ":"`
	Takes       *Type         `@@ "-"`
	Gives       *Type         `">" @@ "="`
	Expressions []*Expression `(@@ ";")+`
}

type Expression struct {
	Pos         lexer.Position
	Application *Application `( @@`
	Primary     *Primary     `| @@ )`
}

type Type struct {
	Pos       lexer.Position
	Primative *Ident  ` @@`
	Vector    *Type   `| "[" @@ "]"`
	Tuple     []*Type `| "(" (@@ ("," @@)*)? ")"`
}

type Application struct {
	Pos       lexer.Position
	Function  *Ident      `@@`
	Parameter *Expression `@@`
}

type Primary struct {
	Pos    lexer.Position
	Tuple  []*Expression `	"(" (@@ ("," @@)*)? ")"`
	Vec    []*Expression `| "[" (@@ ("," @@)*)? "]"`
	Int    *int64        `| @('-'? Int)`
	Real   *float64      `| @('-'? Float)`
	Nil    *string       `| @"Nil"`
	String *string       `| @String`
	Char   *string       `| @Char`
	Param  *string       `| @"@"`
	Noun   *Ident        `| @@`
}
