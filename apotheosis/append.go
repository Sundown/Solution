package apotheosis

import (
	"github.com/sundown/solution/prism"

	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

func (env *Environment) compileInlineAppend(alpha prism.Value, omega prism.Value) value.Value {
	vec_t := alpha.Type.Realise()
	elmT := alpha.Type.(prism.VectorType).Type.Realise()
	elmWidth := I32(alpha.Type.(prism.VectorType).Type.Width())
	lenA := env.readVectorLength(alpha)

	// Length of first vector * size of each element, extended to i64
	sext_lenA_mul := env.Block.NewSExt(env.Block.NewMul(lenA, elmWidth), types.I64)
	lenB := env.readVectorLength(omega)

	capF := env.Block.NewAdd(
		env.readVectorCapacity(alpha),
		env.readVectorCapacity(omega))

	head := env.Block.NewAlloca(vec_t)

	env.writeLLVectorLength(prism.Value{Value: head, Type: alpha.Type}, env.Block.NewAdd(lenA, lenB))
	env.writeLLVectorCapacity(prism.Value{Value: head, Type: alpha.Type}, capF)

	body := env.Block.NewCall(
		env.GetCalloc(),
		I32(alpha.Type.(prism.VectorType).Type.Width()),
		capF)

	// memcpy(body, alpha, alpha len)
	env.Block.NewCall(env.GetMemcpy(),
		body,
		env.Block.NewBitCast(
			env.Block.NewLoad(
				types.NewPointer(elmT),
				env.Block.NewGetElementPtr(vec_t, alpha.Value, I32(0), vectorBodyOffset)),
			types.I8Ptr),
		sext_lenA_mul,
		constant.NewBool(false))

	// memcpy(body, omega, omega len)
	env.Block.NewCall(env.GetMemcpy(),
		env.Block.NewIntToPtr(
			env.Block.NewAdd(
				env.Block.NewPtrToInt(body, types.I64),
				sext_lenA_mul),
			types.I8Ptr),
		env.Block.NewBitCast(
			env.Block.NewLoad(
				types.NewPointer(elmT),
				env.Block.NewGetElementPtr(vec_t, omega.Value, I32(0), vectorBodyOffset)),
			types.I8Ptr),
		env.Block.NewSExt(env.Block.NewMul(lenB, elmWidth), types.I64),
		constant.NewBool(false))

	// Point head body at body
	env.Block.NewStore(
		env.Block.NewBitCast(body, types.NewPointer(elmT)),
		env.Block.NewGetElementPtr(vec_t, head, I32(0), vectorBodyOffset))

	return head
}
