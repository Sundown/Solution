package prism

import "github.com/llir/llvm/ir/types"

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

func (Void) Type() Type {
	return VoidType
}

func (i Int) Type() Type {
	return IntType
}

func (r Real) Type() Type {
	return RealType
}

func (c Cast) Type() Type {
	return c.ToType
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

func (a Alpha) Type() Type {
	return a.TypeOf
}

func (a Omega) Type() Type {
	return a.TypeOf
}

func (do DyadicOperator) Type() Type {
	switch do.Operator {
	case KindMapOperator:
		return VectorType{do.Left.Type()}
	case KindFoldlOperator:
		return do.Left.Type()
	}

	panic("Need to impl type")
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
