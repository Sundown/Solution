package prism

import (
	"github.com/llir/llvm/ir/types"
)

type VectorType struct {
	Type
}

// IsVector is shorthand for vector interface check
func IsVector(t Type) bool {
	_, ok := t.(VectorType)
	return ok
}

func VectorDepth(t Type) int {
	if !IsVector(t) {
		return 0
	}

	return 1 + VectorDepth(t.(VectorType).Type)
}

func QueryAutoVector(atom, vec Type) bool {
	if !IsVector(vec) {
		return false
	}

	if VectorDepth(atom) == VectorDepth(vec) {
		return false
	}

	t := vec.(VectorType).Type
	_, err := Delegate(&atom, &t)
	if err != nil {
		panic(*err)
	}

	return true
}

// Interface prism.Type algebraic predicate
func (a VectorType) IsAlgebraic() bool {
	return a.Type.IsAlgebraic()
}

func (v VectorType) Realise() types.Type {
	return types.NewStruct(
		types.I32, types.I32,
		types.NewPointer(v.Type.Realise()))
}

type Vector struct {
	ElementType VectorType
	Body        *[]Expression
}

// Resolve composes Integrate with Derive,
// Fills in sum/generic type based on a concrete type
func (v VectorType) Resolve(t Type) Type {
	return VectorType{Type: Integrate(v.Type, Derive(v, t))}
}

// String function for interface
func (v Vector) String() string {
	var s string
	for _, v := range *v.Body {
		s += v.String() + " "
	}

	return s
}

// String function for interface
func (v VectorType) String() string {
	return "[" + v.Type.String() + "]"
}

// Type property for interface
func (v Vector) Type() Type {
	return v.ElementType
}

// Interface prism.Type width for LLVM codegen
func (v VectorType) Width() int64 {
	return 16
	// (32 + 32 + 64) / 8
	// len + cap + ptr
}

// Interface prism.Type comparison
func (a VectorType) Equals(b Type) bool {
	if t, ok := b.(VectorType); ok {
		return a.Type.Equals(t.Type)
	}

	return false
}

func (v VectorType) Kind() int {
	return TypeKindVector
}

func (t Vector) _atomicflag() {}
