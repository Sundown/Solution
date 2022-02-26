package prism

import (
	"fmt"

	"github.com/llir/llvm/ir/types"
)

type Real struct {
	Value float64
}

// Interface prism.Type comparison
func (t Real) Equals(b Type) bool {
	return b.Kind() == TypeReal
}

// Interface prism.Type width for LLVM codegen
func (t Real) Width() int {
	return RealType.WidthInBytes
}

// Type property for interface
func (r Real) Type() Type {
	return RealType
}

var RealType = AtomicType{
	ID:           TypeReal,
	WidthInBytes: 8,
	Name:         ParseIdent("Real"),
	Actual:       types.Double,
}

// String function for interface
func (r Real) String() string {
	return fmt.Sprintf("%f", r.Value)
}

// Resolve composes Integrate with Derive,
// Should not be used on concrete types
func (i Real) Resolve(t Type) Type {
	panic("Unreachable")
}

// Interface prism.Type algebraic predicate
func (t Real) IsAlgebraic() bool {
	return false
}

func (t Real) _atomicflag() { return }
