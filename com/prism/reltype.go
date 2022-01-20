package prism

import (
	"sundown/solution/palisade"
)

// These are redundant I think
func (a AtomicType) Kind() int {
	return TypeKindAtomic
}

func (v VectorType) Kind() int {
	return TypeKindVector
}

func (s StructType) Kind() int {
	return TypeKindStruct
}

func (a AtomicType) Any() bool {
	return a.AnyType
}

func (v VectorType) Any() bool {
	return v.AnyType
}

func (s StructType) Any() bool {
	return s.AnyType
}

func (v Vector) Kind() int {
	return TypeKindVector
}

// ...

func (env Environment) SubstantiateType(t palisade.Type) Type {
	if t.Primative != nil {
		if ptr := env.Types[Intern(*t.Primative)]; ptr != nil {
			return ptr
		}
	} else if t.Vector != nil {
		return VectorType{
			AnyType: false,
			Type:    env.SubstantiateType(*t.Vector),
		}
	} else if t.Tuple != nil {
		acc := make([]Type, len(t.Tuple))
		for _, cur := range t.Tuple {
			acc = append(acc, env.SubstantiateType(*cur))
		}

		return StructType{FieldTypes: acc}
	}

	Panic("Unknown type")
	return nil
}
