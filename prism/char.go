package prism

import "github.com/llir/llvm/ir/types"

var CharType = AtomicType{
	ID:           TypeChar,
	WidthInBytes: 1,
	Name:         ParseIdent("Char"),
	Actual:       types.I8,
}

// Interface prism.Type comparison
func (t Char) Equals(b Type) bool {
	return b.Kind() == TypeChar
}

// Interface prism.Type width for LLVM codegen
func (t Char) Width() int {
	return CharType.WidthInBytes
}

// Type property for interface
func (c Char) Type() Type {
	return CharType
}

type Char struct {
	Value string
}

// String function for interface
func (c Char) String() string {
	return "'" + string(c.Value) + "'"
}

// Resolve composes Integrate with Derive,
// Should not be used on concrete types
func (i Char) Resolve(t Type) Type {
	panic("Unreachable")
}

// Interface prism.Type algebraic predicate
func (t Char) IsAlgebraic() bool {
	return false
}

func (t Char) _atomicflag() { return }
