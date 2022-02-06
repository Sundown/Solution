package palisade

type Expression struct {
	Monadic   *Monadic  `( @@` // This order is extremely important unfortunately.
	Dyadic    *Dyadic   `| @@`
	Morphemes *Morpheme `| @@ )`
}

type Monadic struct {
	Subexpr    *Expression `(("(" @@ ")")`
	Verb       *Ident      `| @@)`
	Expression *Expression `@@`
}

type Dyadic struct {
	Monadic    *Monadic    `( @@`
	Morphemes  *Morpheme   `| @@ )`
	Verb       *Ident      `@@`
	Expression *Expression `@@?`
}

type Morpheme struct {
	Char   *[]string  `@Char+`
	Alpha  *[]string  `| @Alpha+`
	Omega  *[]string  `| @Omega+`
	Real   *[]float64 `| @('-'? Float)+`
	Int    *[]int64   `| @('-'? Int)+`
	String *[]string  `| @String+`
}
