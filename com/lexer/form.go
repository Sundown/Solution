package lexer

type Expression struct {
	Singletons []*Subexpression `@@+`
}

type Subexpression struct {
	Morpheme *Morpheme   `(@@`
	Sub      *Expression `| ("(" @@ ")"))`
}

type Morpheme struct {
	Ident *string `@Ident`
	Int   *int64  `| @Int`
}
