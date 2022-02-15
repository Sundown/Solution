package prism

func (a AtomicType) Equals(b Type) bool {
	if t, ok := b.(AtomicType); ok {
		return t.Kind() == a.Kind()
	}

	return false
}

func (a VectorType) Equals(b Type) bool {
	if t, ok := b.(VectorType); ok {
		return a.Type.Equals(t.Type)
	}

	return false
}

func (s StructType) Equals(b Type) (acc bool) {
	panic("Not implemented yet")
}

func (t Int) Equals(b Type) bool {
	return b.Kind() == TypeInt
}

func (t Real) Equals(b Type) bool {
	return b.Kind() == TypeReal
}

func (t Char) Equals(b Type) bool {
	return b.Kind() == TypeChar
}

func (t Bool) Equals(b Type) bool {
	return b.Kind() == TypeBool
}

func (t String) Equals(b Type) bool {
	return b.Kind() == TypeString
}

func (s SumType) Equals(b Type) bool {
	Warn("Comparison of sum types")
	return false
}

func (s GenericType) Equals(b Type) bool {
	Warn("Comparison of generic types")
	return false
}
