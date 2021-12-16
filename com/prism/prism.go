package prism

import (
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
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
	ElementType AtomicType
}

type StructType struct {
	FieldTypes []AtomicType
}

type Expression interface {
	Kind() int
	Type() Type
	String() string
	Realise() value.Value
}

type Function struct {
	Name      Ident
	AlphaType Type
	OmegaType Type
	Returns   Type
	Body      *[]Expression
}
