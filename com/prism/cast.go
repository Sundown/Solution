package prism

func DelegateCast(from Expression, to Type) Expression {
	if !QueryCast(from.Type(), to) {
		panic("Can't cast " + from.Type().String() + " to " + to.String())
	}

	switch from.Type().(type) {
	case AtomicType:
		switch from.Type().Kind() {
		case TypeInt:
			return Cast{Value: from, ToType: IntType}
		case TypeReal:
			return Cast{Value: from, ToType: RealType}
		}
	case VectorType:
		if v, ok := to.(VectorType); ok &&
			QueryCast(from.Type().(VectorType).Type, v.Type) {
			return Cast{Value: from, ToType: to}
		}
	}

	panic("not implemented")
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
			case TypeInt, TypeReal, TypeString:
				return true
			}
		case TypeString:
			switch to.Kind() {
			case TypeBool, TypeInt, TypeReal, TypeString:
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
