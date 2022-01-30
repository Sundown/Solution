package palisade

type PalisadeResult struct {
	Environmentments []*struct {
		Function  *Function  `@@`
		Directive *Directive `| @@`
	} `@@*`
}

type Directive struct {
	Command *string `"@" @Ident `
	Value   *string `@Ident ";"`
}

type Ident struct {
	Namespace *string `(@Ident ":" ":")?`
	Ident     *string `@Ident`
}

type Function struct {
	Dyadic *struct {
		Alpha *Type  `@@ `
		Ident *Ident `@@ `
		Omega *Type  `@@ `
	} `("Δ"  @@`

	Monadic *struct {
		Ident *Ident `@@`
		Omega *Type  `@@`
	} `| "Δ" @@)`

	Returns *Type         `"→"  @@  ":" `
	Body    *[]Expression `(@@ ";" )+ "∇"`
}

type Type struct {
	Primative *Ident  ` @@`
	Vector    *Type   `| "[" @@ "]"`
	Tuple     []*Type `| "(" (@@ ("," @@)*)? ")"`
}
