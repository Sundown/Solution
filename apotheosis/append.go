package apotheosis

import (
	"github.com/sundown/solution/prism"

	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

func (env *Environment) compileInlineAppend(alpha prism.Value, omega prism.Value) value.Value {
	vecT := alpha.Type.Realise()
	elmT := alpha.Type.(prism.VectorType).Type.Realise()
	elmWidth := i32(alpha.Type.(prism.VectorType).Type.Width())
	lenA := env.readVectorLength(alpha)

	// Length of first vector * size of each element, extended to i64
	sextLenAMul := env.Block.NewSExt(env.Block.NewMul(lenA, elmWidth), types.I64)
	lenB := env.readVectorLength(omega)

	capF := env.Block.NewAdd(
		env.readVectorCapacity(alpha),
		env.readVectorCapacity(omega))

	head := env.Block.NewAlloca(vecT)

	env.writeLLVectorLength(prism.Value{Value: head, Type: alpha.Type}, env.Block.NewAdd(lenA, lenB))
	env.writeLLVectorCapacity(prism.Value{Value: head, Type: alpha.Type}, capF)

	body := env.Block.NewCall(
		env.getCalloc(),
		i32(alpha.Type.(prism.VectorType).Type.Width()),
		capF)

	// memcpy(body, alpha, |alpha|)
	env.Block.NewCall(env.getMemcpy(),
		body,
		env.Block.NewBitCast(
			env.Block.NewLoad(
				types.NewPointer(elmT),
				env.Block.NewGetElementPtr(vecT, alpha.Value, i32(0), vectorBodyOffset)),
			types.I8Ptr),
		sextLenAMul,
		constant.NewBool(false))

	// memcpy(body, omega, omega len)
	env.Block.NewCall(env.getMemcpy(),
		env.Block.NewIntToPtr(
			env.Block.NewAdd(
				env.Block.NewPtrToInt(body, types.I64),
				sextLenAMul),
			types.I8Ptr),
		env.Block.NewBitCast(
			env.Block.NewLoad(
				types.NewPointer(elmT),
				env.Block.NewGetElementPtr(vecT, omega.Value, i32(0), vectorBodyOffset)),
			types.I8Ptr),
		env.Block.NewSExt(env.Block.NewMul(lenB, elmWidth), types.I64),
		constant.NewBool(false))

	// Point head body at body
	env.Block.NewStore(
		env.Block.NewBitCast(body, types.NewPointer(elmT)),
		env.Block.NewGetElementPtr(vecT, head, i32(0), vectorBodyOffset))

	return head
}
