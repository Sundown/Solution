package apotheosis

import (
	"sundown/solution/subtle"

	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

func (state *State) CompileInlineAppend(
	typAlpha *subtle.Type, valAlpha value.Value,
	typOmega *subtle.Type, valOmega value.Value) value.Value {
	width := typAlpha.WidthInBytes()

	vec_head_typ := typAlpha.AsLLType()
	vec_elem_typ := types.NewPointer(typAlpha.Vector.AsLLType())

	len_a := state.ReadVectorLength(typAlpha, valAlpha)
	len_b := state.ReadVectorLength(typOmega, valOmega)

	cap_f := state.Block.NewAdd(
		state.ReadVectorCapacity(typAlpha, valAlpha),
		state.ReadVectorCapacity(typOmega, valOmega))

	head := state.Block.NewAlloca(vec_head_typ)

	state.Block.NewStore(
		state.Block.NewAdd(len_a, len_b),
		state.Block.NewGetElementPtr(vec_head_typ, head, I32(0), vectorLenOffset))

	state.Block.NewStore(cap_f,
		state.Block.NewGetElementPtr(vec_head_typ, head, I32(0), vectorCapOffset))

	body := state.Block.NewBitCast(
		state.Block.NewCall(state.GetCalloc(), I32(width),
			state.Block.NewTrunc(cap_f, types.I32)),
		vec_elem_typ)

	state.Block.NewCall(state.GetMemcpy(), state.Block.NewBitCast(body, types.I8Ptr),
		state.Block.NewBitCast(
			state.Block.NewLoad(
				vec_elem_typ,
				state.Block.NewGetElementPtr(vec_head_typ, valAlpha, I32(0), vectorBodyOffset)),
			types.I8Ptr),
		len_a, constant.NewBool(false))

	state.Block.NewCall(state.GetMemcpy(),
		state.Block.NewIntToPtr(
			state.Block.NewAdd(len_a, state.Block.NewPtrToInt(body, types.I64)), types.I8Ptr),
		state.Block.NewBitCast(
			state.Block.NewLoad(
				vec_elem_typ,
				state.Block.NewGetElementPtr(vec_head_typ, valOmega, I32(0), vectorBodyOffset)),
			types.I8Ptr),
		len_b, constant.NewBool(false))

	state.WriteVectorPointer(head, vec_head_typ, body)

	return head
}
