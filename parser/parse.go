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

type FnDecl struct {
	Ident       *string       `@Ident ":"`
	Type        *Expression   `@@ "="`
	Expressions []*Expression `(@@ ";")*`
}

type Expression struct {
	Application *Application `( @@`
	Type        *Type        `| @@`
	Primary     *Primary     `| @@ )`
	Op          *string      `( @(":"":" | "-"">" | "." )`
	Binary      *Expression  `@@)?`
}

type TypeName struct {
	Type *string `@("Int" | "Nat" | "Real" | "Bool" | "Str" | "Char" | "Void")`
}

type Type struct {
	Primative *TypeName ` @@`
	Vector    *Type     `| "[" @@ "]"`
	Struct    []*Type   `| "(" (@@ ("," @@)*)? ")"`
}

type Application struct {
	Function  *string     `@Ident`
	Parameter *Expression `@@`
}

type Primary struct {
	Tuple  []*Expression `  "(" (@@ ("," @@)*)? ")"`
	Vec    []*Expression `| "[" (@@ ("," @@)*)? "]"`
	Int    *int64        `| @Int`
	Real   *float64      `| @Float`
	Bool   *string       `| @("True" | "False")`
	Nil    *string       `| @"Nil"`
	String *string       `| @String`
	Param  *string       `| @"%"`
	Noun   *string       `| @Ident`
}

var Parser = participle.MustBuild(&Program{}, participle.UseLookahead(2))
