package compiler

import (
	"sundown/solution/parse"

	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

//
func (state *State) CompileInlineAppend(arg *parse.Expression) value.Value {
	vec_head_typ := arg.Atom.TypeOf.Tuple[0].AsLLType()
	vec_elem_typ := types.NewPointer(arg.Atom.TypeOf.Tuple[0].Vector.AsLLType())
	greater_vector := state.CompileExpression(arg)
	vec_a := state.Block.NewGetElementPtr(
		arg.TypeOf.AsLLType(), greater_vector, I32(0), I32(0))
	vec_b := state.Block.NewGetElementPtr(
		arg.TypeOf.AsLLType(), greater_vector, I32(0), I32(1))
	len_a := state.Block.NewLoad(types.I64,
		state.Block.NewGetElementPtr(vec_head_typ, vec_a, I32(0), vectorLenOffset))
	len_b := state.Block.NewLoad(types.I64,
		state.Block.NewGetElementPtr(vec_head_typ, vec_b, I32(0), vectorLenOffset))

	cap_f := state.Block.NewAdd(
		state.Block.NewLoad(types.I64,
			state.Block.NewGetElementPtr(vec_head_typ, vec_a, I32(0), vectorCapOffset)),
		state.Block.NewLoad(types.I64,
			state.Block.NewGetElementPtr(vec_head_typ, vec_b, I32(0), vectorCapOffset)))

	head := state.Block.NewAlloca(vec_head_typ)
	state.Block.NewStore(
		state.Block.NewAdd(len_a, len_b),
		state.Block.NewGetElementPtr(vec_head_typ, head, I32(0), vectorLenOffset))
	state.Block.NewStore(cap_f,
		state.Block.NewGetElementPtr(vec_head_typ, head, I32(0), vectorCapOffset))

	body := state.Block.NewBitCast(
		state.Block.NewCall(state.GetCalloc(), I32(arg.Atom.TypeOf.Tuple[0].WidthInBytes()),
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
