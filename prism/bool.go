package prism

import "github.com/llir/llvm/ir/types"

var BoolType = AtomicType{
	ID:           TypeBool,
	WidthInBytes: 1,
	Name:         ParseIdent("Bool"),
	Actual:       types.I1,
}

// Interface Type comparison
func (t Bool) Equals(b Type) bool {
	return b.Kind() == TypeBool
}

// Interface Type width for LLVM codegen
func (t Bool) Width() int {
	return BoolType.WidthInBytes
}

// Type property for interface
func (b Bool) Type() Type {
	return BoolType
}

// String function for interface
func (b Bool) String() string {
	if b.Value {
		return "True"
	}

	return "False"
}

// Resolve composes Integrate with Derive,
// Should not be used on concrete types
func (i Bool) Resolve(t Type) Type {
	Panic("Unreachable")
	return nil
}

type Bool struct {
	Value bool
}

// Interface Type algebraic predicate
func (t Bool) IsAlgebraic() bool {
	return false
}

func (t Bool) _atomicflag() { return }
