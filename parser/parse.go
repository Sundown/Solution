package parser

import (
	"github.com/alecthomas/participle/v2"
)

type Program struct {
	Expression []*Expression `@@*`
}

type Block struct {
	Expression []*Expression `":" @@* ";"`
}

type FnDecl struct {
	Ident *Ident `@@ "="`
	Type  *Type  `@@`
	Block *Block `@@`
}

type Expression struct {
	FnDecl      *FnDecl      `@@`
	Application *Application `| @@`
	Type        *Type        `| @@`
	Primary     *Primary     `| @@`
	Block       *Block       `| @@`
}

type TypeName struct {
	Type *string `@("int" | "nat" | "real" | "bool" | "str" | "char" | "void")`
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
	String *string       `| @String`
	Param  *string       `| @"%"`
	Noun   *Ident        `| @@`
	Bool   *string       `| @("true" | "false")`
	Nil    bool          `| @"void"`
}

var Parser = participle.MustBuild(&Program{}, participle.UseLookahead(2))
