package prism

func (a AtomicType) Width() int64 {
	return int64(a.WidthInBytes)
}

func (v VectorType) Width() int64 {
	return 16
	// (32 + 32 + 64) / 8
	// len + cap + ptr
}

func (s StructType) Width() (acc int64) {
	for _, v := range s.FieldTypes {
		acc += v.Width()
	}

	return acc
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

func (s SumType) Width() int64 {
	panic("Impossible")
}

func (s GenericType) Width() int64 {
	panic("Impossible")
}
