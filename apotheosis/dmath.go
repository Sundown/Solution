package apotheosis

import (
	"github.com/sundown/solution/prism"

	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

func (env *Environment) compileInlineRightTacD(alpha, omega prism.Value) value.Value {
	return omega.Value
}

func (env *Environment) compileInlineRightTacM(omega prism.Value) value.Value {
	return omega.Value
}

func (env *Environment) compileInlineLeftTacD(alpha, omega prism.Value) value.Value {
	return alpha.Value
}

func (env *Environment) compileInlineDAdd(alpha, omega prism.Value) value.Value {
	switch alpha.Type.Kind() {
	case prism.RealType.ID:
		return env.Block.NewFAdd(alpha.Value, omega.Value)
	case prism.IntType.ID:
		return env.Block.NewAdd(alpha.Value, omega.Value)
	case prism.CharType.ID:
		return env.Block.NewAdd(alpha.Value, omega.Value)
	}
	prism.Panic("unreachable")
	panic(nil)
}

func (env *Environment) compileInlineMSub(omega prism.Value) value.Value {
	switch omega.Type.Kind() {
	case prism.RealType.ID:
		return env.Block.NewFSub(f64(0), omega.Value)
	case prism.IntType.ID:
		return env.Block.NewSub(i64(0), omega.Value)
	case prism.CharType.ID:
		return env.Block.NewSub(constant.NewInt(types.I8, 0), omega.Value)
	}

	prism.Panic("unreachable")
	panic(nil)
}

func (env *Environment) compileInlinePow(alpha, omega prism.Value) value.Value {
	switch alpha.Type.Kind() {
	case prism.RealType.ID:
		return env.Block.NewCall(env.getPowReal(), alpha.Value, omega.Value)
	case prism.IntType.ID:
		return env.Block.NewCall(env.getPowInt(),
			env.Block.NewSIToFP(alpha.Value, types.Double),
			env.Block.NewTrunc(omega.Value, types.I32))
	case prism.CharType.ID:
		return env.Block.NewCall(env.getPowInt(),
			env.Block.NewSIToFP(alpha.Value, types.Double),
			env.Block.NewSExt(omega.Value, types.I32))
	}

	prism.Panic("unreachable")
	panic(nil)
}

func (env *Environment) compileInlineExp(omega prism.Value) value.Value {
	switch omega.Type.Kind() {
	case prism.RealType.ID:
		return env.Block.NewCall(env.getPowReal(), f64(2.7182818284590452354), omega.Value)
	case prism.IntType.ID:
		return env.Block.NewCall(env.getPowInt(),
			f64(2.7182818284590452354),
			env.Block.NewTrunc(omega.Value, types.I32))
	case prism.CharType.ID:
		return env.Block.NewCall(env.getPowInt(),
			f64(2.7182818284590452354),
			env.Block.NewSExt(omega.Value, types.I32))
	}

	prism.Panic("unreachable")
	panic(nil)
}

func (env *Environment) compileInlineDSub(alpha, omega prism.Value) value.Value {
	switch alpha.Type.Kind() {
	case prism.RealType.ID:
		return env.Block.NewFSub(alpha.Value, omega.Value)
	case prism.IntType.ID:
		return env.Block.NewSub(alpha.Value, omega.Value)
	case prism.CharType.ID:
		return env.Block.NewSub(alpha.Value, omega.Value)
	}

	prism.Panic("unreachable")
	panic(nil)
}

func (env *Environment) compileInlineMul(alpha, omega prism.Value) value.Value {
	switch alpha.Type.Kind() {
	case prism.RealType.ID:
		return env.Block.NewFMul(alpha.Value, omega.Value)
	case prism.IntType.ID:
		return env.Block.NewMul(alpha.Value, omega.Value)
	case prism.CharType.ID:
		return env.Block.NewMul(alpha.Value, omega.Value)
	}

	prism.Panic("unreachable")
	panic(nil)
}

func (env *Environment) compileInlineDiv(alpha, omega prism.Value) value.Value {
	switch alpha.Type.Kind() {
	case prism.IntType.ID:
		return env.Block.NewFDiv(env.Block.NewSIToFP(alpha.Value, types.Double), env.Block.NewSIToFP(omega.Value, types.Double))
	case prism.RealType.ID:
		return env.Block.NewFDiv(alpha.Value, omega.Value)
	}

	prism.Panic("unreachable")
	panic(nil)
}

func (env *Environment) compileInlineMax(alpha, omega prism.Value) value.Value {
	switch alpha.Type.Kind() {
	case prism.RealType.ID:
		return env.Block.NewCall(env.getMaxDouble(), alpha.Value, omega.Value)
	case prism.IntType.ID:
		// Branchless max, very important especially in array languages
		// a - ((a-b) & (a-b) >> 31)
		i1 := env.Block.NewSub(alpha.Value, omega.Value)
		return env.Block.NewSub(alpha.Value,
			env.Block.NewAnd(i1, env.Block.NewAShr(i1, i64(31))))
	case prism.CharType.ID:
		// Might work, might not, who cares
		i1 := env.Block.NewSub(alpha.Value, omega.Value)
		return env.Block.NewSub(alpha.Value,
			env.Block.NewAnd(i1, env.Block.NewAShr(i1, i64(31))))
	}

	prism.Panic("unreachable")
	panic(nil)
}

func (env *Environment) compileInlineMin(alpha, omega prism.Value) value.Value {
	switch alpha.Type.Kind() {
	case prism.RealType.ID:
		return env.Block.NewCall(env.getMinDouble(), alpha.Value, omega.Value)
	case prism.IntType.ID:
		// Branchless max
		// b + ((a-b) & (a-b) >> 31)
		i1 := env.Block.NewSub(alpha.Value, omega.Value)
		return env.Block.NewAdd(omega.Value,
			env.Block.NewAnd(i1, env.Block.NewAShr(i1, i64(31))))

	case prism.CharType.ID:
		return env.Block.NewMul(alpha.Value, omega.Value)
	}

	prism.Panic("unreachable")
	panic(nil)
}

func (env *Environment) compileInlineCeil(omega prism.Value) value.Value {
	switch omega.Type.Kind() {
	case prism.RealType.ID:
		return env.Block.NewSIToFP(
			env.Block.NewAdd(
				env.Block.NewFPToSI(omega.Value, types.I64),
				i64(1)), types.Double)
	case prism.IntType.ID:
		return omega.Value
	}

	prism.Panic("unreachable")
	panic(nil)
}

func (env *Environment) compileInlineFloor(omega prism.Value) value.Value {
	switch omega.Type.Kind() {
	case prism.RealType.ID:
		return env.Block.NewSIToFP(env.Block.NewFPToSI(
			omega.Value, types.I64), types.Double)
	case prism.IntType.ID:
		return omega.Value
	}

	prism.Panic("unreachable")
	panic(nil)
}
