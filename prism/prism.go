package prism

import (
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

// For Apotheosis
type Value struct {
	Value value.Value
	Type  Type
}

func Val(v value.Value, t Type) Value {
	return Value{Value: v, Type: t}
}

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
	Attrs() Attribute
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
