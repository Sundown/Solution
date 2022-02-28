package palisade

type PalisadeResult struct {
	Environmentments []*struct {
		Function  *Function  `parser:"@@"`
		Directive *Directive `parser:"| @@"`
	} `parser:"@@*"`
}

type Directive struct {
	Command *string `parser:"'@' @Ident"`
	Value   *string `parser:"@Ident ';'"`
}

type Ident struct {
	Namespace *string `parser:"(@Ident ':' ':')?"`
	Ident     *string `parser:"@Ident"`
}

type Function struct {
	Dyadic *struct {
		Alpha *Type  `parser:"@@ "`
		Ident *Ident `parser:"@@ "`
		Omega *Type  `parser:"@@ "`
	} `parser:"('Δ'  @@"`

	Monadic *struct {
		Ident *Ident `parser:"@@"`
		Omega *Type  `parser:"@@"`
	} `parser:"| 'Δ' @@)"`

	Returns *Type         `parser:"'→'  @@  ':'"`
	Body    *[]Expression `parser:"(@@ ';')+ '∇'"`
}

type Type struct {
	Primitive *Ident  `parser:"@@"`
	Vector    *Type   `parser:"| '[' @@ ']'"`
	Tuple     []*Type `parser:"| '(' (@@ (',' @@)*)? ')'"`
}
