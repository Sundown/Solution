package palisade

import "sundown/solution/oversight"

type PalisadeResult struct {
	Runtime    *oversight.Runtime
	Statements []*struct {
		Function  *Function  `@@`
		Directive *Directive `| @@`
	} `@@*`
}

type Directive struct {
	Ident *string `":" @Ident`
	Value *string `@String`
}

type Ident struct {
	Namespace *string `(@Ident ":" ":")?`
	Ident     *string `@Ident`
}

type Function struct {
	Dyadic *struct {
		Alpha *Type  `@@`
		Ident *Ident `@@`
		Omega *Type  `@@`
	} `("Δ" @@`

	Monadic *struct {
		Ident *Ident `@@`
		Omega *Type  `@@`
	} `| "Δ" @@)`

	Returns *Type         `"→" @@ ":"`
	Body    *[]Expression `(@@ ";")+ "∇"`
}

type Type struct {
	Primative *Ident  ` @@`
	Vector    *Type   `| "[" @@ "]"`
	Tuple     []*Type `| "(" (@@ ("," @@)*)? ")"`
}
