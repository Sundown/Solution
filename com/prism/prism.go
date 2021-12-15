package prism

import "github.com/llir/llvm/ir/types"

type Ident struct {
	Package string
	Name    string
}

const (
	TypeKindAtomic = iota
	TypeKindVector
	TypeKindStruct
	TypeInt
	TypeReal
	TypeChar
)

type Type struct {
	Name      Ident
	Substance interface {
		Kind() int
		Width() int64
		String() string
		Realise() types.Type
	}
}

type AtomicType struct {
	ID     int
	Width  int
	Actual types.Type
}

type VectorType struct {
	ElementType AtomicType
}

type StructType struct {
	FieldTypes []AtomicType
}
