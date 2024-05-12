package apotheosis

import (
	"github.com/sundown/solution/prism"

	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

func (env *Environment) newInlineRightTacD(alpha, omega prism.Value) value.Value {
	return omega.Value
}

func (env *Environment) newInlineRightTacM(omega prism.Value) value.Value {
	return omega.Value
}

func (env *Environment) newInlineLeftTacD(alpha, omega prism.Value) value.Value {
	return alpha.Value
}

func (env *Environment) newInlineDAdd(alpha, omega prism.Value) value.Value {
	switch alpha.Type.Kind() {
	case prism.RealType.ID:
		return env.Block.NewFAdd(alpha.Value, omega.Value)
	case prism.IntType.ID:
		return env.Block.NewAdd(alpha.Value, omega.Value)
	case prism.CharType.ID:
		return env.Block.NewAdd(alpha.Value, omega.Value)
	}
	prism.Panic("unreachable")
	panic("Unknown error")
}

func (env *Environment) newInlineMSub(omega prism.Value) value.Value {
	switch omega.Type.Kind() {
	case prism.RealType.ID:
		return env.Block.NewFSub(f64(0), omega.Value)
	case prism.IntType.ID:
		return env.Block.NewSub(i64(0), omega.Value)
	case prism.CharType.ID:
		return env.Block.NewSub(constant.NewInt(types.I8, 0), omega.Value)
	}

	prism.Panic("unreachable")
	panic("Unknown error")
}

func (env *Environment) newInlinePow(alpha, omega prism.Value) value.Value {
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
	panic("Unknown error")
}

func (env *Environment) newInlineExp(omega prism.Value) value.Value {
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
	panic("Unknown error")
}

func (env *Environment) newInlineDSub(alpha, omega prism.Value) value.Value {
	switch alpha.Type.Kind() {
	case prism.RealType.ID:
		return env.Block.NewFSub(alpha.Value, omega.Value)
	case prism.IntType.ID:
		return env.Block.NewSub(alpha.Value, omega.Value)
	case prism.CharType.ID:
		return env.Block.NewSub(alpha.Value, omega.Value)
	}

	prism.Panic("unreachable")
	panic("Unknown error")
}

func (env *Environment) newInlineMul(alpha, omega prism.Value) value.Value {
	switch alpha.Type.Kind() {
	case prism.RealType.ID:
		return env.Block.NewFMul(alpha.Value, omega.Value)
	case prism.IntType.ID:
		return env.Block.NewMul(alpha.Value, omega.Value)
	case prism.CharType.ID:
		return env.Block.NewMul(alpha.Value, omega.Value)
	}

	prism.Panic("unreachable")
	panic("Unknown error")
}

func (env *Environment) newInlineDiv(alpha, omega prism.Value) value.Value {
	switch alpha.Type.Kind() {
	case prism.IntType.ID:
		return env.Block.NewFDiv(env.Block.NewSIToFP(alpha.Value, types.Double), env.Block.NewSIToFP(omega.Value, types.Double))
	case prism.RealType.ID:
		// TODO APO this is messy and needs to apply to all mathematical functions
		// Perhaps this is a good time to take a look at casting and how it works in Subtle too
		if omega.Type.Kind() == prism.IntType.ID {
			return env.Block.NewFDiv(alpha.Value, env.Block.NewSIToFP(omega.Value, types.Double))
		}
		return env.Block.NewFDiv(alpha.Value, omega.Value)
	}

	prism.Panic("unreachable")
	panic("Unknown error")
}

func (env *Environment) newInlineMax(alpha, omega prism.Value) value.Value {
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
	panic("Unknown error")
}

func (env *Environment) newInlineMin(alpha, omega prism.Value) value.Value {
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
	panic("Unknown error")
}

func (env *Environment) newInlineCeil(omega prism.Value) value.Value {
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
	panic("Unknown error")
}

func (env *Environment) newInlineFloor(omega prism.Value) value.Value {
	switch omega.Type.Kind() {
	case prism.RealType.ID:
		return env.Block.NewSIToFP(env.Block.NewFPToSI(
			omega.Value, types.I64), types.Double)
	case prism.IntType.ID:
		return omega.Value
	}

	prism.Panic("unreachable")
	panic("Unknown error")
}
