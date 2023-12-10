package palisade

type Expression struct {
	Bool      *[]string `parser:"( @( 'true' | 'false' )"`
	Monadic   *Monadic  `parser:"| @@"`
	Dyadic    *Dyadic   `parser:"| @@"`
	Morphemes *Morpheme `parser:"| @@ )"`
}

type Monadic struct {
	Applicable *Applicable `parser:"@@"`
	Expression *Expression `parser:"@@?"`
}

type Dyadic struct {
	Monadic    *Monadic    `parser:"( @@"`
	Morphemes  *Morpheme   `parser:"| @@ )"`
	Applicable *Applicable `parser:"@@"`
	Expression *Expression `parser:"@@?"`
}

type Operator struct {
	Operator *string   `parser:"@Operator"`
	Next     *Operator `parser:"@@?"`
}

type Applicable struct {
	Subexpr  *Expression `parser:"(('(' @@ ')')"`
	Verb     *Ident      `parser:"| @@)"`
	Operator *Operator   `parser:"@@?"`
}

// TODO allow 1 2 omega 3 ...
type Morpheme struct {
	Char   *[]string  `parser:"@Char+"`
	Alpha  *[]string  `parser:"| @Alpha+"`
	Omega  *[]string  `parser:"| @Omega+"`
	Real   *[]float64 `parser:"| @('-'? Float)+"`
	Int    *[]int64   `parser:"| @('-'? Int)+"`
	String *[]string  `parser:"| @String+"`
}
