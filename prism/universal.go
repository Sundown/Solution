package prism

import "github.com/llir/llvm/ir/types"

type Universal struct{}

func (g Universal) Kind() int {
	return TypeKindGroup
}

func (g Universal) IsAlgebraic() bool {
	return true
}
func (g Universal) Width() int64 {
	panic("Impossible")
}
func (g Universal) Realise() types.Type {
	panic("Impossible")
}

// Resolve composes Integrate with Derive,
// Fills in sum/generic type based on a concrete type
func (s Universal) Resolve(t Type) Type {
	return t
}

func (Universal) String() string {
	return "T"
}

func (g Universal) Equals(t Type) bool {
	return t.Kind() == TypeKindGroup && t.(Group).Universal()
}

func (Universal) Universal() bool {
	return true
}

func (g Universal) Union(t Group) Group {
	return g
}

func (g Universal) Intersection(t Group) Group {
	return t
}

func (group Universal) Has(typ Type) bool {
	return true
}
