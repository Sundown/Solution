package apotheosis

import (
	"github.com/sundown/solution/prism"

	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

func (env *Environment) compileInlineAdd(alpha, omega Value) value.Value {
	switch alpha.Type.Kind() {
	case prism.RealType.ID:
		return env.Block.NewFAdd(alpha.Value, omega.Value)
	case prism.IntType.ID:
		return env.Block.NewAdd(alpha.Value, omega.Value)
	case prism.CharType.ID:
		return env.Block.NewAdd(alpha.Value, omega.Value)
	case prism.TypeKindVector:
		return env.CombineOf(DCallable(env.compileInlineAdd), alpha, omega)
	}

	panic("unreachable")
}

func (env *Environment) compileInlineSub(alpha, omega Value) value.Value {
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

func (env *Environment) compileInlineMul(alpha, omega Value) value.Value {
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

func (env *Environment) compileInlineDiv(alpha, omega Value) value.Value {
	switch alpha.Type.Kind() {
	case prism.IntType.ID:
		return env.Block.NewFDiv(env.Block.NewSIToFP(alpha.Value, types.Double), env.Block.NewSIToFP(omega.Value, types.Double))
	case prism.RealType.ID:
		return env.Block.NewFDiv(alpha.Value, omega.Value)
	}

	panic("unreachable")
}

func (env *Environment) compileInlineMax(alpha, omega Value) value.Value {
	switch alpha.Type.Kind() {
	case prism.RealType.ID:
		return env.Block.NewCall(env.GetMaxDouble(), alpha.Value, omega.Value)
	case prism.IntType.ID:
		// Branchless max, very important especially in array languages
		// a - ((a-b) & (a-b) >> 31)
		i1 := env.Block.NewSub(alpha.Value, omega.Value)
		return env.Block.NewSub(alpha.Value,
			env.Block.NewAnd(i1, env.Block.NewAShr(i1, I64(31))))
	case prism.CharType.ID:
		// Might work, might not, who cares
		i1 := env.Block.NewSub(alpha.Value, omega.Value)
		return env.Block.NewSub(alpha.Value,
			env.Block.NewAnd(i1, env.Block.NewAShr(i1, I64(31))))
	}

	panic("unreachable")
}

func (env *Environment) compileInlineMin(alpha, omega Value) value.Value {
	switch alpha.Type.Kind() {
	case prism.RealType.ID:
		return env.Block.NewCall(env.GetMinDouble(), alpha.Value, omega.Value)
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

func (env *Environment) compileInlineCeil(omega Value) value.Value {
	switch omega.Type.Kind() {
	case prism.RealType.ID:
		return env.Block.NewSIToFP(
			env.Block.NewAdd(
				env.Block.NewFPToSI(omega.Value, types.I64),
				I64(1)), types.Double)
	case prism.IntType.ID:
		return omega.Value
	}

	panic("unreachable")
}

func (env *Environment) compileInlineFloor(omega Value) value.Value {
	switch omega.Type.Kind() {
	case prism.RealType.ID:
		return env.Block.NewSIToFP(env.Block.NewFPToSI(
			omega.Value, types.I64), types.Double)
	case prism.IntType.ID:
		return omega.Value
	}

	panic("unreachable")
}
