package prism

import "github.com/llir/llvm/ir/types"

func EqType(a, b Type) bool {
	if a.Any() || b.Any() {
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

func (f DFunction) Type() Type {
	return f.Returns
}

func (f MFunction) Type() Type {
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
