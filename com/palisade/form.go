package palisade

type Expression struct {
	Lexemes []*Subexpression `@@+`
}

type Subexpression struct {
	Morpheme *Morpheme   `(@@`
	Sub      *Expression `| ("(" @@ ")"))`
	Term     *bool
}

type Morpheme struct {
	Char   *string  `@Char`
	Alpha  *string  `| @"α"`
	Omega  *string  `| @"ω"`
	Ident  *Ident   `| @@`
	Int    *int64   `| @('-'? Int)`
	Real   *float64 `| @('-'? Float)`
	String *string  `| @String`
}
