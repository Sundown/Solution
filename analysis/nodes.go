package analysis

type Program struct {
	Statements []*Statement
	Directives []*Directive
}

type Statement struct {
	Ident string
	Value *Expression
}

type Directive struct {
	Class       string
	Instruction struct {
		String, Ident string
		Number        float64
	}
}

type Type struct {
	Atomic string
	Vector *Type
	Struct []*Type
}

type Function struct {
	Name  string
	Takes *Expression
	Gives *Expression
}

type Application struct {
	TypeOf   *Type
	Function *Function
	Argument *Expression
}

type Atom struct {
	TypeOf *Type
	Struct []*Expression
	Vector []*Expression
	Int    int64
	Nat    uint64
	Real   float64
	Bool   bool
	String string
	Noun   string
	Param  uint
}

type Expression struct {
	TypeOf              *Type
	Binary, Application *Application
	Atom                *Atom
	Type                *Type
	Block               []*Expression
}
