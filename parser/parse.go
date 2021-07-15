package parser

import (
	"github.com/alecthomas/participle/v2"
)

type Program struct {
	Statements []*struct {
		Directive *Directive `"@" @@`
		FnDecl    *FnDecl    `| @@`
	} `@@*`
}

type Directive struct {
	Class *string `@Ident`
	Instr *struct {
		Ident  *string  `( @Ident`
		String *string  `| @String`
		Number *float64 `| @Float)`
	} `@@ ";"`
}

type Ident struct {
	Namespace *string `(@Ident ":" ":")?`
	Ident     *string `@Ident`
}

type FnDecl struct {
	Ident       *Ident        `@@ ":"`
	Takes       *Type         `@@ "-"`
	Gives       *Type         `">" @@ "="`
	Expressions []*Expression `(@@ ";")*`
}

type Expression struct {
	Application *Application `( @@`
	Type        *Type        `| @@`
	Primary     *Primary     `| @@ )`
}

type TypeName struct {
	Type *string `@("Int" | "Nat" | "Real" | "Bool" | "Str" | "Char" | "Void")`
}

type Type struct {
	Primative *TypeName ` @@`
	Vector    *Type     `| "[" @@ "]"`
	Tuple     []*Type   `| "(" (@@ ("," @@)*)? ")"`
}

type Application struct {
	Function  *Ident      `@@`
	Parameter *Expression `@@`
}

type Primary struct {
	Tuple  []*Expression `	"(" (@@ ("," @@)*)? ")"`
	Vec    []*Expression `| "[" (@@ ("," @@)*)? "]"`
	Int    *int64        `| @Int`
	Real   *float64      `| @Float`
	Bool   *string       `| @("True" | "False")`
	Nil    *string       `| @"Nil"`
	String *string       `| @String`
	Param  *string       `| @"%"`
	Noun   *Ident        `| @@`
}

var Parser = participle.MustBuild(&Program{}, participle.UseLookahead(3))
