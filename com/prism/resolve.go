package prism

func (v VectorType) Resolve(t Type) Type {
	return Integrate(v, Derive(v, t))
}

func (s SumType) Resolve(t Type) Type {
	return Integrate(s, Derive(s, t))
}

func (g GenericType) Resolve(t Type) Type {
	return Integrate(g, Derive(g, t))
}

func (s StructType) Resolve(t Type) Type {
	panic("Not implemented yet")
}

func (a AtomicType) Resolve(t Type) Type {
	panic("Unreachable")
}

func (i Int) Resolve(t Type) Type {
	panic("Unreachable")
}

func (i Real) Resolve(t Type) Type {
	panic("Unreachable")
}

func (i Char) Resolve(t Type) Type {
	panic("Unreachable")
}

func (i Bool) Resolve(t Type) Type {
	panic("Unreachable")
}

func (i String) Resolve(t Type) Type {
	panic("Unreachable")
}
