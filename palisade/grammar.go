package palisade

type Expression struct {
	// This order is extremely important
	Operator  *Operator `parser:"( @@"`
	Monadic   *Monadic  `parser:"| @@"`
	Dyadic    *Dyadic   `parser:"| @@"`
	Morphemes *Morpheme `parser:"| @@ )"`
}

type Monadic struct {
	Subexpr    *Expression `parser:"(('(' @@ ')')"`
	Verb       *Ident      `parser:"| @@)"`
	Expression *Expression `parser:"@@?"`
}

type Dyadic struct {
	Monadic    *Monadic    `parser:"( @@"`
	Morphemes  *Morpheme   `parser:"| @@ )"`
	Subexpr    *Expression `parser:"(('(' @@ ')')"`
	Verb       *Ident      `parser:"| @@)"`
	Expression *Expression `parser:"@@?"`
}

type Operator struct {
	Subexpr    *Expression `parser:"(('(' @@ ')')"`
	Verb       *Ident      `parser:"| @@)"`
	Operator   *string     `parser:"@Operator"`
	Expression *Expression `parser:"@@"`
}

type Morpheme struct {
	Char   *[]string  `parser:"@Char+"`
	Alpha  *[]string  `parser:"| @Alpha+"`
	Omega  *[]string  `parser:"| @Omega+"`
	Real   *[]float64 `parser:"| @('-'? Float)+"`
	Int    *[]int64   `parser:"| @('-'? Int)+"`
	String *[]string  `parser:"| @String+"`
}
