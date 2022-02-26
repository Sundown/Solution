package prism

import "github.com/llir/llvm/ir/types"

// Interface prism.Type width for LLVM codegen
func (s SumType) Width() int64 {
	panic("Impossible")
}

func (s SumType) String() (res string) {
	for i, t := range s.Types {
		if i > 0 {
			res += " | "
		}
		res += t.String()
	}

	return
}

// Resolve composes Integrate with Derive,
// Fills in sum/generic type based on a concrete type
func (s SumType) Resolve(t Type) Type {
	return Integrate(s, Derive(s, t))
}

type SumType struct {
	Types []Type
}

func (s SumType) Realise() types.Type {
	panic("Impossible")
}

func (s SumType) Kind() int {
	return TypeKindSemiDeterminedGroup
}

// Interface prism.Type algebraic predicate
func (s SumType) IsAlgebraic() bool {
	return true
}

// Interface prism.Type comparison
func (s SumType) Equals(b Type) bool {
	Warn("Comparison of sum types")
	return false
}
