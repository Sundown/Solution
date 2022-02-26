package apotheosis

import (
	"github.com/sundown/solution/prism"

	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

func (env *Environment) CompileInlineAnd(alpha, omega Value) value.Value {
	switch alpha.Type.Kind() {
	case prism.BoolType.ID:
		return env.Block.NewAnd(alpha.Value, omega.Value)
	case prism.RealType.ID:
		return env.Block.NewUIToFP(env.Block.NewAnd(
			env.Block.NewFCmp(enum.FPredOGT, alpha.Value, F64(0)),
			env.Block.NewFCmp(enum.FPredOGT, omega.Value, F64(0))), types.Double)
	case prism.IntType.ID:
		return env.Block.NewZExt(env.Block.NewAnd(
			env.Block.NewICmp(enum.IPredSGT, alpha.Value, I64(0)),
			env.Block.NewICmp(enum.IPredSGT, omega.Value, I64(0))), types.I64)
	case prism.CharType.ID:
		return env.Block.NewZExt(env.Block.NewAnd(
			env.Block.NewICmp(enum.IPredSGT, alpha.Value, constant.NewInt(types.I8, 0)),
			env.Block.NewICmp(enum.IPredSGT, omega.Value, constant.NewInt(types.I8, 0))), types.I8)
	}

	panic("unreachable")
}
func (env *Environment) CompileInlineOr(alpha, omega Value) value.Value {
	switch alpha.Type.Kind() {
	case prism.BoolType.ID:
		return env.Block.NewOr(alpha.Value, omega.Value)
	case prism.RealType.ID:
		return env.Block.NewUIToFP(env.Block.NewOr(
			env.Block.NewFCmp(enum.FPredOGT, alpha.Value, F64(0)),
			env.Block.NewFCmp(enum.FPredOGT, omega.Value, F64(0))), types.Double)
	case prism.IntType.ID:
		return env.Block.NewZExt(env.Block.NewOr(
			env.Block.NewICmp(enum.IPredSGT, alpha.Value, I64(0)),
			env.Block.NewICmp(enum.IPredSGT, omega.Value, I64(0))), types.I64)
	case prism.CharType.ID:
		return env.Block.NewZExt(env.Block.NewOr(
			env.Block.NewICmp(enum.IPredSGT, alpha.Value, constant.NewInt(types.I8, 0)),
			env.Block.NewICmp(enum.IPredSGT, omega.Value, constant.NewInt(types.I8, 0))), types.I8)

	}

	panic("unreachable")
}
