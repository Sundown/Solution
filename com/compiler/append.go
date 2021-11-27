package compiler

import (
	"sundown/solution/temporal"

	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

func (state *State) CompileInlineAppend(typ *temporal.Type, val value.Value) value.Value {
	width := typ.Tuple[0].WidthInBytes()
	vec_head_typ := typ.Tuple[0].AsLLType()
	vec_elem_typ := types.NewPointer(typ.Tuple[0].Vector.AsLLType())

	vec_a := state.Block.NewGetElementPtr(typ.AsLLType(), val, I32(0), I32(0))
	vec_b := state.Block.NewGetElementPtr(typ.AsLLType(), val, I32(0), I32(1))

	len_a := state.ReadVectorLength(typ.Tuple[0], vec_a)
	len_b := state.ReadVectorLength(typ.Tuple[1], vec_b)

	cap_f := state.Block.NewAdd(
		state.ReadVectorCapacity(typ.Tuple[0], vec_a),
		state.ReadVectorCapacity(typ.Tuple[1], vec_b))

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
				state.Block.NewGetElementPtr(vec_head_typ, vec_a, I32(0), vectorBodyOffset)),
			types.I8Ptr),
		len_a, constant.NewBool(false))

	state.Block.NewCall(state.GetMemcpy(),
		state.Block.NewIntToPtr(
			state.Block.NewAdd(len_a, state.Block.NewPtrToInt(body, types.I64)), types.I8Ptr),
		state.Block.NewBitCast(
			state.Block.NewLoad(
				vec_elem_typ,
				state.Block.NewGetElementPtr(vec_head_typ, vec_b, I32(0), vectorBodyOffset)),
			types.I8Ptr),
		len_b, constant.NewBool(false))

	state.WriteVectorPointer(head, vec_head_typ, body)

	return head
}
