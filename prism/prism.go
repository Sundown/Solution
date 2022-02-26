package prism

import (
	"github.com/llir/llvm/ir/types"
)

type Expression interface {
	Type() Type
	String() string
}

type Type interface {
	Kind() int
	Width() int64
	String() string
	Equals(Type) bool
	Realise() types.Type
	Resolve(Type) Type
	IsAlgebraic() bool
}

type Function interface {
	LLVMise() string
	IsSpecial() bool
	ShouldInline() bool
	Type() Type
	Ident() Ident
	String() string
}

type Ident struct {
	Package string
	Name    string
}
type Morpheme interface {
	_atomicflag()
}

// String function for interface
func (i Ident) String() string {
	if i.Package == "_" {
		return i.Name
	}

	return i.Package + "::" + i.Name
}
