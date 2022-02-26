package prism

import "github.com/llir/llvm/ir/types"

type Void struct{}

var VoidType = AtomicType{
	ID:           TypeVoid,
	WidthInBytes: 0,
	Name:         ParseIdent("Void"),
	Actual:       types.Void,
}

// Type property for interface
func (Void) Type() Type {
	return VoidType
}

// String function for interface
func (Void) String() string {
	return "Void"
}
