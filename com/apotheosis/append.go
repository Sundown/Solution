package apotheosis

import (
	"sundown/solution/prism"

	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

// self explanatory
func (state *State) CompileInlineAppend(alpha Value, omega Value) (head value.Value) {
	vec_head_typ := alpha.Type.Realise()
	vec_elem_typ := types.NewPointer(alpha.Type.(prism.VectorType).Realise())

	len_a := state.ReadVectorLength(alpha)
	len_b := state.ReadVectorLength(omega)
	cap_f := state.Block.NewAdd(state.ReadVectorCapacity(alpha),
		state.ReadVectorCapacity(omega))
	head = state.Block.NewAlloca(vec_head_typ)
	state.Block.NewStore(state.Block.NewAdd(len_a, len_b),
		state.Block.NewGetElementPtr(vec_head_typ, head, I32(0), vectorLenOffset))
	state.Block.NewStore(cap_f, state.Block.NewGetElementPtr(
		vec_head_typ, head, I32(0), vectorCapOffset))
	body := state.Block.NewBitCast(state.Block.NewCall(state.GetCalloc(),
		I32(alpha.Type.Width()), state.Block.NewTrunc(cap_f, types.I32)), vec_elem_typ)
	state.Block.NewCall(state.GetMemcpy(), state.Block.NewBitCast(body, types.I8Ptr),
		state.Block.NewBitCast(state.Block.NewLoad(vec_elem_typ,
			state.Block.NewGetElementPtr(vec_head_typ, alpha.Value, I32(0), vectorBodyOffset)),
			types.I8Ptr), len_a, constant.NewBool(false))
	state.Block.NewCall(state.GetMemcpy(), state.Block.NewIntToPtr(
		state.Block.NewAdd(len_a, state.Block.NewPtrToInt(body, types.I64)), types.I8Ptr),
		state.Block.NewBitCast(state.Block.NewLoad(vec_elem_typ,
			state.Block.NewGetElementPtr(vec_head_typ, omega.Value, I32(0), vectorBodyOffset)),
			types.I8Ptr), len_b, constant.NewBool(false))
	state.WriteVectorPointer(Value{head, alpha.Type}, body)

	return
}
