package prism

import (
	"github.com/sundown/solution/palisade"

	"github.com/llir/llvm/ir/types"
)

type Expression interface {
	Type() Type
	String() string
}

type Type interface {
	Kind() int
	Width() int64
	String() string
	Equals(Type) bool
	Realise() types.Type
	Resolve(Type) Type
	IsAlgebraic() bool
}

type Function interface {
	LLVMise() string
	IsSpecial() bool
	ShouldInline() bool
	Type() Type
	Ident() Ident
	String() string
}

type Ident struct {
	Package string
	Name    string
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
	Special     bool
	SkipBuilder bool
	Inline      bool
	Name        Ident
	AlphaType   Type
	OmegaType   Type
	Returns     Type
	PreBody     *[]palisade.Expression
	Body        []Expression
}

type MonadicFunction struct {
	Special     bool
	SkipBuilder bool
	Inline      bool
	Name        Ident
	OmegaType   Type

	Returns Type
	PreBody *[]palisade.Expression
	Body    []Expression
}

type MonadicApplication struct {
	Operator MonadicFunction
	Operand  Expression
}

type DyadicApplication struct {
	Operator DyadicFunction
	Left     Expression
	Right    Expression
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
