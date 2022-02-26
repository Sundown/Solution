package prism

import "github.com/llir/llvm/ir/types"

type StructType struct {
	FieldTypes []Type
}

// Interface prism.Type comparison
func (s StructType) Equals(b Type) (acc bool) {
	panic("Not implemented yet")
}

// Interface prism.Type algebraic predicate
func (s StructType) IsAlgebraic() (acc bool) {
	panic("Not implemented yet")
}

// Interface prism.Type width for LLVM codegen
func (s StructType) Width() (acc int64) {
	for _, v := range s.FieldTypes {
		acc += v.Width()
	}

	return acc
}

func (s StructType) Realise() types.Type {
	acc := []types.Type{}
	for _, v := range s.FieldTypes {
		acc = append(acc, v.Realise())
	}

	return types.NewStruct(acc...)
}

// String function for interface
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

func (s StructType) Kind() int {
	return TypeKindStruct
}

// Resolve composes Integrate with Derive,
// Fills in sum/generic type based on a concrete type
func (s StructType) Resolve(t Type) Type {
	panic("Not implemented yet")
}
