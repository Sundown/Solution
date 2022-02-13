package prism

import (
	"sundown/solution/palisade"

	"github.com/llir/llvm/ir/types"
)

type Ident struct {
	Package string
	Name    string
}

type Function interface {
	LLVMise() string
	IsSpecial() bool
	Type() Type
	Ident() Ident
	String() string
}

func (d DyadicFunction) IsSpecial() bool {
	return d.Special
}

func (d MonadicFunction) IsSpecial() bool {
	return d.Special
}

func (d DyadicFunction) Ident() Ident {
	return d.Name
}

func (m MonadicFunction) Ident() Ident {
	return m.Name
}

type Vector struct {
	ElementType VectorType
	Body        *[]Expression
}

type Expression interface {
	Type() Type
	String() string
}

type Void struct{}

type DyadicOperator struct {
	Operator int
	Left     Expression
	Right    Expression
	Returns  Type
}

type Morpheme interface {
	_atomicflag()
}

type Cast struct {
	Value  Expression
	ToType Type
}

type DyadicFunction struct {
	Special   bool
	Name      Ident
	AlphaType Type
	OmegaType Type
	Returns   Type
	PreBody   *[]palisade.Expression
	Body      []Expression
}

type MonadicFunction struct {
	Special   bool
	Name      Ident
	OmegaType Type

	Returns Type
	PreBody *[]palisade.Expression
	Body    []Expression
}

type MApplication struct {
	Operator MonadicFunction
	Operand  Expression
	// Algebraic T refers to which type in given application
	T_Refers Type
}

type DApplication struct {
	Operator DyadicFunction
	Left     Expression
	Right    Expression
	// Algebraic T refers to which type in given application
	// Possibility of requiring A and B algebraic types once tuples are present
	T_Refers Type
}

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
	Type
}

type StructType struct {
	FieldTypes []Type
}
type GenericType struct{}

type SumType struct {
	Types []Type
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

type Alpha struct {
	TypeOf Type
}
type Omega struct {
	TypeOf Type
}

type EOF struct{}
