package prism

import "github.com/llir/llvm/ir/types"

const (
	TypeKindAtomic = iota
	TypeKindVector
	TypeKindStruct
	TypeKindSome
	KinDyadicFunction
	TypeInt
	TypeReal
	TypeChar
	TypeBool
	TypeVoid
	TypeString
)

type Type interface {
	Any() bool
	SemiDetermined() bool
	Kind() int
	Width() int64
	String() string
	Realise() types.Type
}

type AtomicType struct {
	AnyType      bool
	ID           int
	WidthInBytes int
	Name         Ident
	Actual       types.Type
}

type VectorType struct {
	Type
	AnyType bool
}

type StructType struct {
	AnyType    bool
	FieldTypes []Type
}

func (s SemiDeterminedType) SemiDetermined() bool {
	return true
}

func (s VectorType) SemiDetermined() bool {
	return false
}
func (s TypeGroup) SemiDetermined() bool {
	return false
}
func (s AtomicType) SemiDetermined() bool {
	return false
}
func (s StructType) SemiDetermined() bool {
	return false
}

func (s TypeGroup) Any() bool {
	return false
}

func (s SemiDeterminedType) Any() bool {
	return false
}

func (s TypeGroup) Kind() int {
	return TypeKindSome
}

func (s TypeGroup) Width() int64 {
	panic("Impossible")
}

func (s SemiDeterminedType) Width() int64 {
	panic("Impossible")
}

func (s SemiDeterminedType) Realise() types.Type {
	panic("Impossible")
}

type SemiDeterminedType struct{}

func (s SemiDeterminedType) String() string {
	return "T"
}

func (s TypeGroup) String() (res string) {
	for i, t := range s.Types {
		if i > 0 {
			res += " | "
		}
		res += t.String()
	}

	return
}

func (s TypeGroup) Realise() types.Type {
	panic("Impossible")
}

type TypeGroup struct {
	Types []Type
}

func EqType(a, b Type) bool {
	if a.Any() || b.Any() {
		return true
	}

	p := func(x TypeGroup, y Type) bool {
		for _, t := range x.Types {
			if EqType(t, y) {
				return true
			}
		}

		return false
	}

	if s, ok := a.(TypeGroup); ok && p(s, b) {
		return true
	} else if s, ok := b.(TypeGroup); ok && p(s, a) {
		return true
	}

	if a.Kind() != b.Kind() {
		return false
	}

	switch a.Kind() {
	case TypeKindAtomic:
		return a.(AtomicType).ID == b.(AtomicType).ID
	case TypeKindVector:
		return EqType(a.(VectorType).Type, b.(VectorType).Type)
	case TypeKindStruct:
		// TODO
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
