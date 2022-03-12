package palisade

import "github.com/alecthomas/participle/v2/lexer"

type PalisadeResult struct {
	Environmentments []*struct {
		Function  *Function  `parser:"@@"`
		Directive *Directive `parser:"| @@"`
	} `parser:"@@*"`
}

type Directive struct {
	Pos     lexer.Position
	Command *string `parser:"'@' @Ident"`
	Value   *string `parser:"@Ident ';'"`
}

type Ident struct {
	Pos       lexer.Position
	Namespace *string `parser:"(@Ident ':' ':')?"`
	Ident     *string `parser:"@Ident"`
}

type Function struct {
	Pos    lexer.Position
	Dyadic *struct {
		Alpha *Type  `parser:"@@ "`
		Ident *Ident `parser:"@@ "`
		Omega *Type  `parser:"@@ "`
	} `parser:"(@@"`

	Monadic *struct {
		Ident *Ident `parser:"@@"`
		Omega *Type  `parser:"@@"`
	} `parser:"| @@)"`

	Returns *Type         `parser:"'→'  @@  '{'"`
	Body    *[]Expression `parser:"(@@ ';')+ '}'"`
}

type Type struct {
	Pos       lexer.Position
	Primitive *Ident  `parser:"@@"`
	Vector    *Type   `parser:"| '[' @@ ']'"`
	Tuple     []*Type `parser:"| '(' (@@ (',' @@)*)? ')'"`
}
