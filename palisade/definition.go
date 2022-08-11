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
	TypedFunction *struct {
		Dyadic *struct {
			Alpha *Type  `parser:"@@ "`
			Ident *Ident `parser:"@@ "`
			Omega *Type  `parser:"@@ "`
		} `parser:"(@@"`
		Monadic *struct {
			Ident *Ident `parser:"@@"`
			Omega *Type  `parser:"@@"`
		} `parser:"| @@)"`
		Returns *Type `parser:"'→'  @@ "`
	} `parser:"(@@ |"`

	AmbiguousFunction *struct {
		Ident   *Ident `parser:"@@"`
		Returns *Type  `parser:"('→' @@)?"`
	} `parser:"@@)"`

	Body  *[]Expression `parser:"(('{' (@@ ';')+ '}')"`
	Tacit *Expression   `parser:"| (@@ ';'))"`
}

type Type struct {
	Primitive *Ident  `parser:"@@"`
	Vector    *Type   `parser:"| '[' @@ ']'"`
	Tuple     []*Type `parser:"| '(' (@@ (',' @@)*)? ')'"`
}
