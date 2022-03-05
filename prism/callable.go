package prism

import (
	"github.com/llir/llvm/ir/value"
)

type Callable interface {
	Arity() int
	NoAutoVector() bool
}
type MCallable func(val Value) value.Value
type DCallable func(left, right Value) value.Value

type MonadicCallable struct {
	MCallable
	disallowAutoVector bool
}

func MakeDC(d DCallable, noAutoVec bool) Callable {
	return DyadicCallable{DCallable: DCallable(d), disallowAutoVector: noAutoVec}
}

func MakeMC(d MCallable, noAutoVec bool) Callable {
	return MonadicCallable{MCallable: MCallable(d), disallowAutoVector: noAutoVec}
}

func (m MonadicCallable) NoAutoVector() bool {
	return m.disallowAutoVector
}

func (m DyadicCallable) NoAutoVector() bool {
	return m.disallowAutoVector
}

func (m MonadicFunction) NoAutoVector() bool {
	return m.disallowAutoVector
}

func (m DyadicFunction) NoAutoVector() bool {
	return m.disallowAutoVector
}

type DyadicCallable struct {
	DCallable
	disallowAutoVector bool
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

func (DyadicFunction) Arity() int {
	return 2
}

func (MonadicFunction) Arity() int {
	return 1
}
