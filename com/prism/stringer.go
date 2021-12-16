package prism

import (
	"github.com/llir/llvm/ir/types"
)

func (a AtomicType) String() string {
	return a.Name.String()
}

func (v VectorType) String() string {
	return "[" + v.ElementType.String() + "]"
}

func (s StructType) String() (acc string) {
	acc = "("
	for i, v := range s.FieldTypes {
		if i > 0 {
			acc += ", "
		}

		acc += v.String()
	}

	return acc + ")"
}

func (i Ident) String() string {
	if i.Package == "_" {
		return i.Name
	}

	return i.Package + "::" + i.Name
}

func (a AtomicType) Width() int64 {
	return int64(a.WidthInBytes)
}

func (v VectorType) Width() int64 {
	return v.ElementType.Width() + 16
}

func (s StructType) Width() (acc int64) {
	for _, v := range s.FieldTypes {
		acc += v.Width()
	}

	return acc
}

func (a AtomicType) Realise() types.Type {
	return a.Actual
}

func (v VectorType) Realise() types.Type {
	return types.NewStruct(
		types.I64, types.I64,
		types.NewPointer(v.ElementType.Realise()))
}

func (s StructType) Realise() types.Type {
	acc := []types.Type{}
	for _, v := range s.FieldTypes {
		acc = append(acc, v.Realise())
	}

	return types.NewStruct(acc...)
}
