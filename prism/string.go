package prism

import "github.com/llir/llvm/ir/types"

var StringType = AtomicType{
	ID:           TypeString,
	WidthInBytes: 12, // TODO
	Name:         ParseIdent("String"),
	Actual:       types.I8Ptr,
}

// Interface Type comparison
func (t String) Equals(b Type) bool {
	return b.Kind() == TypeString
}

// Resolve composes Integrate with Derive,
// Should not be used on concrete types
func (i String) Resolve(t Type) Type {
	Panic("Unreachable")
	panic(nil)
}

type String struct {
	Value string
}

// String function for interface
func (s String) String() string {
	return "\"" + s.Value + "\""
}

// Type property for interface
func (s String) Type() Type {
	return StringType
}

// Interface Type width for LLVM codegen
func (t String) Width() int {
	return StringType.WidthInBytes
}

// Interface Type algebraic predicate
func (t String) IsAlgebraic() bool {
	return false
}

func (t String) _atomicflag() { return }
