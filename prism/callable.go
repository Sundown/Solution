package prism

import (
	"github.com/llir/llvm/ir/value"
)

type Callable interface {
	Arity() int
}
type MCallable func(val Value) value.Value
type DCallable func(left, right Value) value.Value

func (MCallable) Arity() int {
	return 1
}

func (DCallable) Arity() int {
	return 2
}

func (env Environment) FetchDCallable(v string) Callable {
	if found, ok := env.LLDyadicCallables[v]; ok {
		return found
	}

	return nil
}

func (env Environment) FetchMCallable(v string) Callable {
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
