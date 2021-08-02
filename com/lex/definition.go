package lex

type State struct {
	Statements []*struct {
		Directive *Directive `"@" @@`
		TypeDecl  *TypeDecl  `| @@`
		NounDecl  *NounDecl  `| @@`
		FnDecl    *FnDecl    `| @@`
	} `@@*`
}

type Directive struct {
	Class *string `@Ident`
	Instr *struct {
		Ident  *string  `( @Ident`
		String *string  `| @String`
		Number *float64 `| @Float)`
	} `@@`
}

type Ident struct {
	Namespace *string `(@Ident ":" ":")?`
	Ident     *string `@Ident`
}

type TypeDecl struct {
	Ident *string `@Ident "~"`
	Type  *Type   `@@`
}

type NounDecl struct {
	Ident *string  `@Ident "="`
	Value *Primary `@@ ";"`
}

type FnDecl struct {
	Ident       *string       `@Ident ":"`
	Takes       *Type         `@@ "-"`
	Gives       *Type         `">" @@ "="`
	Expressions []*Expression `(@@ ";")+`
}

type Expression struct {
	Application *Application `( @@`
	Primary     *Primary     `| @@ )`
}

type Type struct {
	Primative *Ident  ` @@`
	Vector    *Type   `| "[" @@ "]"`
	Tuple     []*Type `| "(" (@@ ("," @@)*)? ")"`
}

type Application struct {
	Function  *Ident      `@@`
	Parameter *Expression `@@`
}

type Primary struct {
	Tuple  []*Expression `	"(" (@@ ("," @@)*)? ")"`
	Vec    []*Expression `| "[" (@@ ("," @@)*)? "]"`
	Int    *int64        `| @('-'? Int)`
	Real   *float64      `| @('-'? Float)`
	Nil    *string       `| @"Nil"`
	String *string       `| @String`
	Char   *string       `| @Char`
	Param  *string       `| @"@"`
	Noun   *Ident        `| @@`
}
