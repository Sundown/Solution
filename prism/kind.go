package prism

const (
	TypeKindAtomic = iota
	TypeKindVector
	TypeKindStruct
	TypeKindSemiDetermined
	TypeKindSemiDeterminedGroup
	TypeKindDyadicFunction
	KindMapOperator
	KindFoldlOperator
	TypeInt
	TypeReal
	TypeChar
	TypeBool
	TypeVoid
	TypeString
)

func (a AtomicType) Kind() int {
	return a.ID
}

func (v VectorType) Kind() int {
	return TypeKindVector
}

func (s StructType) Kind() int {
	return TypeKindStruct
}

func (v Vector) Kind() int {
	return TypeKindVector
}

func (s GenericType) Kind() int {
	return TypeKindSemiDetermined
}

func (s SumType) Kind() int {
	return TypeKindSemiDeterminedGroup
}
