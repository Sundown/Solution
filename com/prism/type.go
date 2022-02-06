package prism

import "github.com/llir/llvm/ir/types"

const (
	TypeKindAtomic = iota
	TypeKindVector
	TypeKindStruct
	TypeKindSemiDetermined
	TypeKindSemiDeterminedGroup
	TypeKindDyadicFunction
	KindMapOperator
	KindFoldlOperator
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
	Type
}

type StructType struct {
	FieldTypes []Type
}

func (s GenericType) Kind() int {
	return TypeKindSemiDetermined
}

func (s SumType) Kind() int {
	return TypeKindSemiDeterminedGroup
}

func (s SumType) Width() int64 {
	panic("Impossible")
}

func (s GenericType) Width() int64 {
	panic("Impossible")
}

func (s GenericType) Realise() types.Type {
	panic("Impossible")
}

type GenericType struct{}

func (s GenericType) String() string {
	return "T"
}

func (s SumType) String() (res string) {
	for i, t := range s.Types {
		if i > 0 {
			res += " | "
		}
		res += t.String()
	}

	return
}

func (s SumType) Realise() types.Type {
	panic("Impossible")
}

type SumType struct {
	Types []Type
}

func PrimativeTypeEq(a, b Type) bool {
	if a.Kind() == TypeKindSemiDetermined || b.Kind() == TypeKindSemiDetermined {
		return false
	}

	/* p := func(x SumType, y Type) bool {
		for _, t := range x.Types {
			if PrimativeTypeEq(t, y) {
				return true
			}
		}

		return false
	} */

	/* 	if s, ok := a.(SumType); ok && p(s, b) {
	   		return true
	   	} else if s, ok := b.(SumType); ok && p(s, a) {
	   		return true
	   	} */
	if _, ok := a.(SumType); ok {
		return false
	} else if _, ok := b.(SumType); ok {
		return false
	}
	if a.Kind() != b.Kind() {
		return false
	}

	switch a.(type) {
	case AtomicType:
		return a.(AtomicType).ID == b.(AtomicType).ID
	case VectorType:
		return PrimativeTypeEq(a.(VectorType).Type, b.(VectorType).Type)
	case StructType:
		// TODO
		// ... other kinds
	}

	return false
}

func (v Vector) Type() Type {
	return v.ElementType
}

func (f DyadicFunction) Type() Type {
	return f.Returns
}

func (f MonadicFunction) Type() Type {
	return f.Returns
}

func (m MApplication) Type() Type {
	return m.Operator.Type()
}

func (d DApplication) Type() Type {
	return d.Operator.Type()
}

func (i Int) Type() Type {
	return IntType
}

func (r Real) Type() Type {
	return RealType
}

func (c Char) Type() Type {
	return CharType
}

func (b Bool) Type() Type {
	return BoolType
}

func (s String) Type() Type {
	return StringType
}

func (t Int) Width() int {
	return IntType.WidthInBytes
}

func (t Real) Width() int {
	return RealType.WidthInBytes
}

func (t Char) Width() int {
	return CharType.WidthInBytes
}

func (t Bool) Width() int {
	return BoolType.WidthInBytes
}

func (t String) Width() int {
	return StringType.WidthInBytes
}

var (
	IntType = AtomicType{
		ID:           TypeInt,
		WidthInBytes: 8,
		Name:         ParseIdent("Int"),
		Actual:       types.I64,
	}
	RealType = AtomicType{
		ID:           TypeReal,
		WidthInBytes: 8,
		Name:         ParseIdent("Real"),
		Actual:       types.Double,
	}
	CharType = AtomicType{
		ID:           TypeChar,
		WidthInBytes: 1,
		Name:         ParseIdent("Char"),
		Actual:       types.I8,
	}
	StringType = AtomicType{
		ID:           TypeString,
		WidthInBytes: 12, // TODO
		Name:         ParseIdent("String"),
		Actual:       types.I8Ptr,
	}
	BoolType = AtomicType{
		ID:           TypeBool,
		WidthInBytes: 1,
		Name:         ParseIdent("Bool"),
		Actual:       types.I1,
	}
	VoidType = AtomicType{
		ID:           TypeVoid,
		WidthInBytes: 0,
		Name:         ParseIdent("Void"),
		Actual:       types.Void,
	}
)
