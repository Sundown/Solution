package lexer

type Expression struct {
	Singletons []*Subexpression `@@+`
}

type Subexpression struct {
	Morpheme *Morpheme   `(@@`
	Sub      *Expression `| ("(" @@ ")"))`
	Term     *bool
}

type Morpheme struct {
	Ident  *string  `@Ident`
	Int    *int64   `| @('-'? Int)`
	Real   *float64 `| @('-'? Float)`
	String *string  `| @String`
	Char   *string  `| @Char`
	Alpha  *string  `| @"Alpha"`
	Omega  *string  `| @"Omega"`
}
