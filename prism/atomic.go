package prism

import "github.com/llir/llvm/ir/types"

// Interface prism.Type width for LLVM codegen
func (a AtomicType) Width() int64 {
	return int64(a.WidthInBytes)
}

// String function for interface
func (a AtomicType) String() string {
	return a.Name.String()
}

// Resolve composes Integrate with Derive,
// Should not be used on concrete types
func (a AtomicType) Resolve(t Type) Type {
	Panic("Unreachable")
	return nil
}

type AtomicType struct {
	ID           int
	WidthInBytes int
	Name         Ident
	Actual       types.Type
}

func (a AtomicType) Realise() types.Type {
	return a.Actual
}

// Interface prism.Type algebraic predicate
func (a AtomicType) IsAlgebraic() bool {
	return false
}

// Interface prism.Type comparison
func (a AtomicType) Equals(b Type) bool {
	if t, ok := b.(AtomicType); ok {
		return t.Kind() == a.Kind()
	}

	return false
}

func (a AtomicType) Kind() int {
	return a.ID
}
