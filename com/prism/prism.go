package prism

import (
	"github.com/llir/llvm/ir/types"
)

type Environment struct {
	Functions map[Ident]Function
	Types     map[Ident]Type
}

type Ident struct {
	Package string
	Name    string
}

const (
	TypeKindAtomic = iota
	TypeKindVector
	TypeKindStruct
	KindFunction
	TypeInt
	TypeReal
	TypeChar
	TypeBool
	TypeVoid
	TypeString
)

type Type interface {
	Kind() int
	Width() int64
	String() string
	Realise() types.Type
}

type AtomicType struct {
	ID           int
	WidthInBytes int
	Name         Ident
	Actual       types.Type
}

type VectorType struct {
	ElementType Type
}

type StructType struct {
	FieldTypes []Type
}

type Expression interface {
	Type() Type
	String() string
}

type Subexpression struct {
	Expression Expression
}

type Function struct {
	Name      Ident
	AlphaType Type
	OmegaType Type
	Returns   Type
	PreBody   *[]Expression
	Body      *[]Expression
}

type Monadic struct {
	Operator Function
	Operand  Expression
}

type Dyadic struct {
	Operator Function
	Left     Expression
	Right    Expression
}

type Application struct {
	Operator Function
	Operand  Expression
}

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

type Bool struct {
	Value bool
}

type Alpha struct{}
type Omega struct{}
type EOF struct{}
