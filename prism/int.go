package prism

import (
	"fmt"

	"github.com/llir/llvm/ir/types"
)

type Int struct {
	Value int64
}

// Interface Type comparison
func (t Int) Equals(b Type) bool {
	return b.Kind() == TypeInt
}

// Resolve composes Integrate with Derive,
// Should not be used on concrete types
func (i Int) Resolve(t Type) Type {
	Panic("Unreachable")
	panic(nil)
}

// Type property for interface
func (i Int) Type() Type {
	return IntType
}

// String function for interface
func (i Int) String() string {
	return fmt.Sprintf("%d", i.Value)
}

// Interface Type width for LLVM codegen
func (t Int) Width() int {
	return IntType.WidthInBytes
}

var IntType = AtomicType{
	ID:           TypeInt,
	WidthInBytes: 8,
	Name:         ParseIdent("Int"),
	Actual:       types.I64,
}

// Interface Type algebraic predicate
func (t Int) IsAlgebraic() bool {
	return false
}

func (t Int) _atomicflag() {}
