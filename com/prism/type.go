package prism

// These are redundant I think
func (a *AtomicType) Kind() int {
	return TypeKindAtomic
}

func (v *VectorType) Kind() int {
	return TypeKindVector
}

func (s *StructType) Kind() int {
	return TypeKindStruct
}

// ...
