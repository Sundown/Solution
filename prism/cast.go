package prism

type Cast struct {
	Value  Expression
	ToType Type
}

// Type property for interface
func (c Cast) Type() Type {
	return c.ToType
}

// String function for interface
func (c Cast) String() string {
	return "<" + c.ToType.String() + ">" + c.Value.String()
}

func DelegateCast(from Expression, to Type) Expression {
	if !QueryCast(from.Type(), to) {
		Panic("Can't cast " + from.Type().String() + " to " + to.String())
	}

	switch from.Type().(type) {
	case AtomicType:
		return Cast{Value: from, ToType: to}
	case VectorType:
		if v, ok := to.(VectorType); ok &&
			QueryCast(from.Type().(VectorType).Type, v.Type) {
			return Cast{Value: from, ToType: to}
		}
	}

	Panic("not implemented")
	panic("unlabelled error")
}

func QueryCast(from, to Type) bool {
	switch from.(type) {
	case AtomicType:
		switch from.Kind() {
		case TypeBool:
			switch to.Kind() {
			case TypeBool, TypeChar, TypeInt, TypeReal, TypeString:
				return true
			}
		case TypeChar:
			switch to.Kind() {
			case TypeChar, TypeInt, TypeReal, TypeString:
				return true
			}
		case TypeInt:
			switch to.Kind() {
			case TypeInt, TypeReal, TypeString:
				return true
			}
		case TypeReal:
			switch to.Kind() {
			case TypeReal, TypeString:
				return true
			}
		case TypeString:
			switch to.Kind() {
			case TypeString:
				return true
			}
		}
	case VectorType:
		if v, ok := to.(VectorType); ok {
			return QueryCast(from.(VectorType).Type, v.Type)
		}
	}

	return false
}
