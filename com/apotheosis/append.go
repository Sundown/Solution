package apotheosis

import (
	"sundown/solution/prism"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

// self explanatory
func (env *Environment) CompileInlineAppend(alpha Value, omega Value) (head value.Value) {
	vec_head_typ := alpha.Type.Realise()
	vec_elem_typ := types.NewPointer(alpha.Type.(prism.VectorType).Realise())

	len_a := env.ReadVectorLength(alpha)
	len_b := env.ReadVectorLength(omega)
	cap_f := env.Block.NewAdd(env.ReadVectorCapacity(alpha),
		env.ReadVectorCapacity(omega))
	head = env.Block.NewAlloca(vec_head_typ)
	env.Block.NewStore(env.Block.NewAdd(len_a, len_b),
		env.Block.NewGetElementPtr(vec_head_typ, head, I32(0), vectorLenOffset))
	env.Block.NewStore(cap_f, env.Block.NewGetElementPtr(
		vec_head_typ, head, I32(0), vectorCapOffset))
	body := env.Block.NewBitCast(env.Block.NewCall(env.GetCalloc(),
		I32(alpha.Type.Width()), env.Block.NewTrunc(cap_f, types.I32)), vec_elem_typ)
	env.Block.NewCall(env.GetMemcpy(), env.Block.NewBitCast(body, types.I8Ptr),
		env.Block.NewBitCast(env.Block.NewLoad(vec_elem_typ,
			env.Block.NewGetElementPtr(vec_head_typ, alpha.Value, I32(0), vectorBodyOffset)),
			types.I8Ptr), len_a, constant.NewBool(false))
	env.Block.NewCall(env.GetMemcpy(), env.Block.NewIntToPtr(
		env.Block.NewAdd(len_a, env.Block.NewPtrToInt(body, types.I64)), types.I8Ptr),
		env.Block.NewBitCast(env.Block.NewLoad(vec_elem_typ,
			env.Block.NewGetElementPtr(vec_head_typ, omega.Value, I32(0), vectorBodyOffset)),
			types.I8Ptr), len_b, constant.NewBool(false))
	env.WriteVectorPointer(head.(*ir.InstAlloca), body, alpha.Type.Realise())

	return
}
