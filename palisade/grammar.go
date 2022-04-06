package palisade

import "github.com/alecthomas/participle/v2/lexer"

type Expression struct {
	Pos       lexer.Position
	Monadic   *Monadic  `parser:"( @@"`
	Dyadic    *Dyadic   `parser:"| @@"`
	Morphemes *Morpheme `parser:"| @@ )"`
}

type Monadic struct {
	Pos        lexer.Position
	Applicable *Applicable `parser:"@@"`
	Expression *Expression `parser:"@@?"`
}

type Dyadic struct {
	Pos        lexer.Position
	Monadic    *Monadic    `parser:"( @@"`
	Morphemes  *Morpheme   `parser:"| @@ )"`
	Applicable *Applicable `parser:"@@"`
	Expression *Expression `parser:"@@?"`
}

type Operator struct {
	Pos      lexer.Position
	Operator *string   `parser:"@Operator"`
	Next     *Operator `parser:"@@?"`
}

type Applicable struct {
	Pos      lexer.Position
	Subexpr  *Expression `parser:"(('(' @@ ')')"`
	Verb     *Ident      `parser:"| @@)"`
	Operator *Operator   `parser:"@@?"`
}

type Morpheme struct {
	Pos    lexer.Position
	Char   *[]string  `parser:"@Char+"`
	Alpha  *[]string  `parser:"| @Alpha+"`
	Omega  *[]string  `parser:"| @Omega+"`
	Real   *[]float64 `parser:"| @('-'? Float)+"`
	Int    *[]int64   `parser:"| @('-'? Int)+"`
	String *[]string  `parser:"| @String+"`
}
