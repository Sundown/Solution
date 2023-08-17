package prism

import (
	"github.com/llir/llvm/ir/value"
)

type Callable interface {
	Attrs() Attribute
	Arity() int
}

type MCallable func(val Value) value.Value
type DCallable func(left, right Value) value.Value

type MonadicCallable struct {
	MCallable
	disallowAutoVector bool
	Attribute          Attribute
}

func MakeDC(d DCallable, noAutoVec bool) Callable {
	return DyadicCallable{DCallable: DCallable(d), Attribute: Attribute{DisallowAutoVector: noAutoVec}}
}

func MakeMC(d MCallable, noAutoVec bool) Callable {
	return MonadicCallable{MCallable: MCallable(d), Attribute: Attribute{DisallowAutoVector: noAutoVec}}
}

func (m MonadicCallable) Attrs() Attribute {
	return m.Attribute
}

func (d DyadicCallable) Attrs() Attribute {
	return d.Attribute
}

type DyadicCallable struct {
	DCallable
	disallowAutoVector bool
	Attribute          Attribute
}

func (MonadicCallable) Arity() int {
	return 1
}

func (DyadicCallable) Arity() int {
	return 2
}

func (env Environment) FetchDyadicCallable(v string) Callable {
	if found, ok := env.LLDyadicCallables[v]; ok {
		return found
	}

	return nil
}

func (env Environment) FetchMonadicCallable(v string) Callable {
	if found, ok := env.LLMonadicCallables[v]; ok {
		return found
	}

	return nil
}
