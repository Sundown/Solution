package prism

import "github.com/llir/llvm/ir/types"

func EqType(a, b Type) bool {
	if a.Kind() != b.Kind() {
		return false
	}

	switch a.Kind() {
	case TypeKindAtomic:
		return a.(AtomicType).ID == b.(AtomicType).ID
	case TypeKindVector:
		return EqType(a.(VectorType).ElementType, b.(VectorType).ElementType)
	case TypeKindStruct:
		// TODO
	}

	return false
}

func (f Function) Type() Type {
	return f.Returns
}

func (m Monadic) Type() Type {
	return m.Operator.Type()
}

func (d Dyadic) Type() Type {
	return d.Operator.Type()
}

func (a Application) Type() Type {
	return a.Operator.Type()
}

func (d Dangle) Type() Type {
	if ok := d.Inner; ok != nil {
		return d.Inner.Type()
	}
	// Just hanging out, no pun intended
	return d.Outer.Type()
}

func (i Int) Type() Type {
	return AtomicType{
		ID:           TypeInt,
		WidthInBytes: 8,
		Name:         ParseIdent("Int"),
		Actual:       types.I64,
	}
}

func (r Real) Type() Type {
	return AtomicType{
		ID:           TypeReal,
		WidthInBytes: 8,
		Name:         ParseIdent("Real"),
		Actual:       types.Double,
	}
}

func (c Char) Type() Type {
	return AtomicType{
		ID:           TypeChar,
		WidthInBytes: 1,
		Name:         ParseIdent("Char"),
		Actual:       types.I8,
	}
}

func (s String) Type() Type {
	return AtomicType{
		ID:           TypeString,
		WidthInBytes: 1,
		Name:         ParseIdent("String"),
		Actual:       types.I8,
	}
}
