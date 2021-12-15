package weave

import (
	"sundown/solution/palisade"
	"sundown/solution/prism"
)

type State struct {
	TypeDefinitions map[prism.Ident]prism.Type
}

type Expression interface {
	String() string
}

type Subexpression struct {
	Expression Expression
}

type Application struct {
	Operator string
	Operand  Expression
}

// Dangle type is for LHS arguments to dyadic operators
// which are left unbound to the RHS until a later analysis pass.
type Dangle struct {
	Outer Expression
	Inner Expression
}

type Int struct {
	Value int64
}

type Real struct {
	Value float64
}

type String struct {
	Value string
}

type Char struct {
	Value string
}

type Alpha struct{}
type Omega struct{}

type Ident struct {
	Value string
}

type EOF struct{}

type SubState struct {
	mutablePos     int
	mutableSubPos  int
	Subexpressions []*palisade.Subexpression
}
