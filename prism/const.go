package prism

const (
	TypeKindAtomic = iota
	TypeKindVector
	TypeKindStruct
	TypeKindSemiDetermined
	TypeKindSemiDeterminedGroup
	KindMapOperator
	KindReduceOperator
	TypeInt
	TypeReal
	TypeChar
	TypeBool
	TypeVoid
	TypeString
)

const CT_OOB = "Index %d overflows %s of length %d"
const CT_Unexpected = "Expected %s, instead found %s"
