package apotheosis

import (
	"github.com/sundown/solution/prism"

	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

func (env *Environment) compileInlineNot(omega prism.Value) value.Value {
	switch omega.Type.Kind() {
	case prism.BoolType.ID:
		return env.Block.NewAnd(omega.Value, constant.NewInt(types.I1, 0))
	case prism.RealType.ID:
		return env.Block.NewFCmp(enum.FPredOEQ, omega.Value, f64(0))
	case prism.IntType.ID:
		return env.Block.NewICmp(enum.IPredEQ, omega.Value, i64(0))
	case prism.CharType.ID:
		return env.Block.NewICmp(enum.IPredEQ, omega.Value, constant.NewInt(types.I8, 0))
	}

	prism.Panic("unreachable")
	panic(nil)
}

func (env *Environment) compileInlineAnd(alpha, omega prism.Value) value.Value {
	switch alpha.Type.Kind() {
	case prism.BoolType.ID:
		return env.Block.NewAnd(alpha.Value, omega.Value)
	case prism.RealType.ID:
		return env.Block.NewUIToFP(env.Block.NewAnd(
			env.Block.NewFCmp(enum.FPredOGT, alpha.Value, f64(0)),
			env.Block.NewFCmp(enum.FPredOGT, omega.Value, f64(0))), types.Double)
	case prism.IntType.ID:
		return env.Block.NewZExt(env.Block.NewAnd(
			env.Block.NewICmp(enum.IPredSGT, alpha.Value, i64(0)),
			env.Block.NewICmp(enum.IPredSGT, omega.Value, i64(0))), types.I64)
	case prism.CharType.ID:
		return env.Block.NewZExt(env.Block.NewAnd(
			env.Block.NewICmp(enum.IPredSGT, alpha.Value, constant.NewInt(types.I8, 0)),
			env.Block.NewICmp(enum.IPredSGT, omega.Value, constant.NewInt(types.I8, 0))), types.I8)
	}

	prism.Panic("unreachable")
	panic(nil)
}
func (env *Environment) compileInlineOr(alpha, omega prism.Value) value.Value {
	switch alpha.Type.Kind() {
	case prism.BoolType.ID:
		return env.Block.NewOr(alpha.Value, omega.Value)
	case prism.RealType.ID:
		return env.Block.NewUIToFP(env.Block.NewOr(
			env.Block.NewFCmp(enum.FPredOGT, alpha.Value, f64(0)),
			env.Block.NewFCmp(enum.FPredOGT, omega.Value, f64(0))), types.Double)
	case prism.IntType.ID:
		return env.Block.NewZExt(env.Block.NewOr(
			env.Block.NewICmp(enum.IPredSGT, alpha.Value, i64(0)),
			env.Block.NewICmp(enum.IPredSGT, omega.Value, i64(0))), types.I64)
	case prism.CharType.ID:
		return env.Block.NewZExt(env.Block.NewOr(
			env.Block.NewICmp(enum.IPredSGT, alpha.Value, constant.NewInt(types.I8, 0)),
			env.Block.NewICmp(enum.IPredSGT, omega.Value, constant.NewInt(types.I8, 0))), types.I8)

	}

	prism.Panic("unreachable")
	panic(nil)
}
