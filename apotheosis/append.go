package apotheosis

import (
	"github.com/sundown/solution/prism"

	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

func (env *Environment) CompileInlineAppend(alpha Value, omega Value) value.Value {
	vec_t := alpha.Type.Realise()
	elm_t := alpha.Type.(prism.VectorType).Type.Realise()
	elm_width := I32(alpha.Type.(prism.VectorType).Type.Width())
	len_a := env.ReadVectorLength(alpha)

	// Length of first vector * size of each element, extended to i64
	sext_len_a_mul := env.Block.NewSExt(env.Block.NewMul(len_a, elm_width), types.I64)
	len_b := env.ReadVectorLength(omega)

	cap_f := env.Block.NewAdd(
		env.ReadVectorCapacity(alpha),
		env.ReadVectorCapacity(omega))

	head := env.Block.NewAlloca(vec_t)

	env.WriteLLVectorLength(Value{head, alpha.Type}, env.Block.NewAdd(len_a, len_b))
	env.WriteLLVectorCapacity(Value{head, alpha.Type}, cap_f)

	body := env.Block.NewCall(
		env.GetCalloc(),
		I32(alpha.Type.(prism.VectorType).Type.Width()),
		cap_f)

	// memcpy(body, alpha, alpha len)
	env.Block.NewCall(env.GetMemcpy(),
		body,
		env.Block.NewBitCast(
			env.Block.NewLoad(
				types.NewPointer(elm_t),
				env.Block.NewGetElementPtr(vec_t, alpha.Value, I32(0), vectorBodyOffset)),
			types.I8Ptr),
		sext_len_a_mul,
		constant.NewBool(false))

	// memcpy(body, omega, omega len)
	env.Block.NewCall(env.GetMemcpy(),
		env.Block.NewIntToPtr(
			env.Block.NewAdd(
				env.Block.NewPtrToInt(body, types.I64),
				sext_len_a_mul),
			types.I8Ptr),
		env.Block.NewBitCast(
			env.Block.NewLoad(
				types.NewPointer(elm_t),
				env.Block.NewGetElementPtr(vec_t, omega.Value, I32(0), vectorBodyOffset)),
			types.I8Ptr),
		env.Block.NewSExt(env.Block.NewMul(len_b, elm_width), types.I64),
		constant.NewBool(false))

	// Point head body at body
	env.Block.NewStore(
		env.Block.NewBitCast(body, types.NewPointer(elm_t)),
		env.Block.NewGetElementPtr(vec_t, head, I32(0), vectorBodyOffset))

	return head
}
