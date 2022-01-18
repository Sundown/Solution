package palisade

type Expression struct {
	Monadic  *Monadic    `( @@`
	Dyadic   *Dyadic     `| @@`
	Morpheme *[]Morpheme `| @@+ )`
}

type Monadic struct {
	Verb       *Verb       `@@`
	Expression *Expression `@@`
}

type Dyadic struct {
	Monadic    *Monadic    `( @@`
	Morphemes  *[]Morpheme `| (@@+) )`
	Verb       *Verb       `@@`
	Expression *Expression `@@?` // possibly broken, leave for now
}

type Verb struct {
	Ident *Ident `@@`
}

type Morpheme struct {
	Char    *string     `@Char`
	Alpha   *string     `| @"α"`
	Omega   *string     `| @"ω"`
	Ident   *Ident      `| "$"@@`
	Int     *int64      `| @('-'? Int)`
	Real    *float64    `| @('-'? Float)`
	String  *string     `| @String`
	Subexpr *Expression `| "(" @@ ")"`
}
