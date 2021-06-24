package parser

import (
	"github.com/alecthomas/participle/v2"
)

type Program struct {
	Expression []*HighExpression `@@*`
}

type Block struct {
	Expression []*Expression `"=" @@* ";"`
}

type Directive struct {
	Class *string `@Ident`
	Instr *string `@Ident ";"`
}

type FnDecl struct {
	Ident *Ident `@@ ":"`
	Type  *Type  `":" @@`
	Block *Block `@@`
}

type HighExpression struct {
	Directive  *Directive  `"@" @@`
	Expression *Expression `| @@`
}

type Expression struct {
	FnDecl      *FnDecl      `@@`
	Application *Application `| @@`
	Type        *Type        `| @@`
	Primary     *Primary     `| @@`
	Block       *Block       `| @@`
}

type TypeName struct {
	Type *string `@("Int" | "Nat" | "Real" | "Bool" | "Str" | "Char" | "Void")`
}

type Type struct {
	Takes *TypeName `(@@ "-"`
	Gives *TypeName `">" @@)`
}

type Application struct {
	Op    *Ident      `@@`
	Atoms *Expression `@@`
}

type Ident struct {
	Ident *string `@Ident`
}

type Primary struct {
	Tuple  []*Expression `  "(" (@@ ("," @@)*)? ")"`
	Vec    []*Expression `| "[" (@@ ("," @@)*)? "]"`
	Int    *int64        `| @Int`
	Real   *float64      `| @Float`
	Bool   *string       `| @("True" | "False")`
	String *string       `| @String`
	Param  *string       `| @"%"`
	Noun   *Ident        `| @@`
	Nil    bool          `| @"Nil"`
}

var Parser = participle.MustBuild(&Program{}, participle.UseLookahead(2))
