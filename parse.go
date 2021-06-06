package main

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/alecthomas/participle/v2"
)

type Program struct {
	Pos lexer.Position

	Expression []*Expression `@@*`
}

type Expression struct {
	Pos lexer.Position

	Application *Application `  @@`
	Primary     *Primary     "| @@"
}

type Application struct {
	Pos lexer.Position

	Op    string        `"(" @Ident`
	Atoms []*Expression ` @@* ")"`
}

type Primary struct {
	Pos lexer.Position

	Nat    *int64   `  @Int`
	Real   *float64 `| @Float`
	String *string  `| @String`
	Noun   *string  `| @Ident`
	Bool   *bool    `| ( @"true" | "false" )`
	Nil    bool     `| @"nil"`
}

var parser = participle.MustBuild(&Program{}, participle.UseLookahead(2))
