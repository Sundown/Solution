package apotheosis

import (
	"sundown/solution/prism"

	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/value"
)

func (env *Environment) CompileInlineAdd(alpha, omega Value) value.Value {
	switch alpha.Type.Kind() {
	case prism.RealType.ID:
		return env.Block.NewFAdd(alpha.Value, omega.Value)
	case prism.IntType.ID:
		return env.Block.NewAdd(alpha.Value, omega.Value)
	case prism.CharType.ID:
		return env.Block.NewAdd(alpha.Value, omega.Value)
	}

	panic("unreachable")
}

func (env *Environment) CompileInlineSub(alpha, omega Value) value.Value {
	switch alpha.Type.Kind() {
	case prism.RealType.ID:
		return env.Block.NewFSub(alpha.Value, omega.Value)
	case prism.IntType.ID:
		return env.Block.NewSub(alpha.Value, omega.Value)
	case prism.CharType.ID:
		return env.Block.NewSub(alpha.Value, omega.Value)
	}

	panic("unreachable")
}

func (env *Environment) CompileInlineMul(alpha, omega Value) value.Value {
	switch alpha.Type.Kind() {
	case prism.RealType.ID:
		return env.Block.NewFMul(alpha.Value, omega.Value)
	case prism.IntType.ID:
		return env.Block.NewMul(alpha.Value, omega.Value)
	case prism.CharType.ID:
		return env.Block.NewMul(alpha.Value, omega.Value)
	}

	panic("unreachable")
}

func (env *Environment) CompileInlineDiv(alpha, omega Value) value.Value {
	switch alpha.Type.Kind() {
	case prism.RealType.ID:
		return env.Block.NewFDiv(alpha.Value, omega.Value)
	}

	panic("unreachable")
}

func (env *Environment) CompileInlineMax(alpha, omega Value) value.Value {
	switch alpha.Type.Kind() {
	case prism.RealType.ID:
		env.Block.NewFCmp(enum.FPredOGT, alpha.Value, omega.Value)
	case prism.IntType.ID:
		// Branchless max, very important especially in array lang like this
		// a - ((a-b) & (a-b) >> 31)
		i1 := env.Block.NewSub(alpha.Value, omega.Value)
		return env.Block.NewSub(alpha.Value,
			env.Block.NewAnd(i1, env.Block.NewAShr(i1, I64(31))))
	case prism.CharType.ID:
		return env.Block.NewMul(alpha.Value, omega.Value)
	}

	panic("unreachable")
}
func (env *Environment) CompileInlineMin(alpha, omega Value) value.Value {
	switch alpha.Type.Kind() {
	case prism.RealType.ID:
		env.Block.NewFCmp(enum.FPredOGT, alpha.Value, omega.Value)
	case prism.IntType.ID:
		// Branchless max
		// b + ((a-b) & (a-b) >> 31)
		i1 := env.Block.NewSub(alpha.Value, omega.Value)
		return env.Block.NewAdd(omega.Value,
			env.Block.NewAnd(i1, env.Block.NewAShr(i1, I64(31))))

	case prism.CharType.ID:
		return env.Block.NewMul(alpha.Value, omega.Value)
	}

	panic("unreachable")
}
