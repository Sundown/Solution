package prism

import "github.com/llir/llvm/ir/types"

type GenericType struct{}

// Interface prism.Type comparison
func (s GenericType) Equals(b Type) bool {
	Warn("Comparison of generic types")
	return false
}

// Interface prism.Type width for LLVM codegen
func (s GenericType) Width() int64 {
	panic("Impossible")
}

func (s GenericType) String() string {
	return "T"
}

// Resolve composes Integrate with Derive,
// Fills in sum/generic type based on a concrete type
func (g GenericType) Resolve(t Type) Type {
	return Integrate(g, Derive(g, t))
}

func (s GenericType) Realise() types.Type {
	panic("Impossible")
}

func (s GenericType) Kind() int {
	return TypeKindSemiDetermined
}

// Interface prism.Type algebraic predicate
func (s GenericType) IsAlgebraic() bool {
	return true
}
