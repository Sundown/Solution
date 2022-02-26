package prism

func (a AtomicType) IsAlgebraic() bool {
	return false
}

func (a VectorType) IsAlgebraic() bool {
	return a.Type.IsAlgebraic()
}

func (s StructType) IsAlgebraic() (acc bool) {
	panic("Not implemented yet")
}

func (t Int) IsAlgebraic() bool {
	return false
}

func (t Real) IsAlgebraic() bool {
	return false
}

func (t Char) IsAlgebraic() bool {
	return false
}

func (t Bool) IsAlgebraic() bool {
	return false
}

func (t String) IsAlgebraic() bool {
	return false
}

func (s SumType) IsAlgebraic() bool {
	return true
}

func (s GenericType) IsAlgebraic() bool {
	return true
}
