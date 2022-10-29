package apotheosis

import (
	"github.com/sundown/solution/prism"

	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

func (env *Environment) newInlineAppend(alpha prism.Value, omega prism.Value) value.Value {
	alphaLen, alphaCap := env.readVectorSizes(alpha)
	omegaLen, _ := env.readVectorSizes(omega)

	totalLen := env.Block.NewAdd(alphaLen, omegaLen)

	lentimessize := env.Block.NewSExt(env.Block.NewMul(omegaLen, i64(alpha.Type.(prism.VectorType).Type.Width())), types.I64)
	elementType := alpha.Type.(prism.VectorType).Type.Realise()
	vectorType := alpha.Type.Realise()

	alphaBody := env.Block.NewLoad(
		types.NewPointer(elementType),
		env.Block.NewGetElementPtr(vectorType, alpha.Value, i32(0), vectorBodyOffset))

	omegaBody := env.Block.NewLoad(
		types.NewPointer(elementType),
		env.Block.NewGetElementPtr(vectorType, omega.Value, i32(0), vectorBodyOffset))

	new := env.newBlock(env.CurrentFunction)
	cur := env.newBlock(env.CurrentFunction)

	// (|a| + |w| <= ||a||)
	env.Block.NewCondBr(env.Block.NewICmp(enum.IPredSLE, totalLen, alphaCap), cur, new)

	env.Block = new

	env.Block.NewCall(env.getRealloc(), alphaBody, env.Block.NewSExt(totalLen, types.I64))
	env.Block.NewBr(cur)

	env.Block = cur

	env.Block.NewCall(env.getMemcpy(),
		env.Block.NewIntToPtr(
			env.Block.NewAdd(
				env.Block.NewPtrToInt(alphaBody, types.I64),
				env.Block.NewSExt(env.Block.NewMul(alphaLen, i64(alpha.Type.(prism.VectorType).Type.Width())), types.I64)),
			types.I8Ptr),
		env.Block.NewBitCast(
			omegaBody,
			types.I8Ptr),
		lentimessize,
		constant.NewBool(false))

	env.writeLLVectorLength(alpha, totalLen)

	exit := env.newBlock(env.CurrentFunction)
	env.Block.NewBr(exit)
	env.Block = exit
	return alpha.Value
}
