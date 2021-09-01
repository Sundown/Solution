package compiler

import (
	"sundown/solution/parse"

	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

func (state *State) CompileInlineAppend(app *parse.Application) value.Value {
	greater_vector := state.CompileExpression(app.Argument)
	vec_head_typ := app.Argument.Atom.Tuple[0].TypeOf.AsLLType()
	vec_elem_typ := app.Argument.Atom.Tuple[0].TypeOf.Vector.AsLLType() // might not work

	vec_a := state.Block.NewGetElementPtr(types.NewPointer(vec_head_typ), greater_vector, I32(0), I32(0))
	vec_b := state.Block.NewGetElementPtr(types.NewPointer(vec_head_typ), greater_vector, I32(0), I32(1))

	len_a := state.Block.NewLoad(types.I64,
		state.Block.NewGetElementPtr(
			vec_head_typ,
			vec_a,
			I32(0), vectorLenOffset))

	len_b := state.Block.NewLoad(types.I64,
		state.Block.NewGetElementPtr(
			vec_head_typ,
			vec_b,
			I32(0), vectorLenOffset))

	cap_a := state.Block.NewLoad(types.I64,
		state.Block.NewGetElementPtr(
			vec_head_typ,
			vec_a,
			I32(0), vectorCapOffset))

	cap_b := state.Block.NewLoad(types.I64,
		state.Block.NewGetElementPtr(
			vec_head_typ,
			vec_b,
			I32(0), vectorCapOffset))

	len_f := state.Block.NewAdd(len_a, len_b)
	cap_f := state.Block.NewAdd(cap_a, cap_b)

	body_a := state.Block.NewLoad(types.NewPointer(vec_elem_typ),
		state.Block.NewGetElementPtr(
			vec_head_typ,
			vec_a,
			I32(0), vectorBodyOffset))

	/* body_b := state.Block.NewLoad(types.NewPointer(vec_elem_typ),
	state.Block.NewGetElementPtr(
		vec_head_typ,
		vec_b,
		I32(0), vectorBodyOffset)) */

	head := state.Block.NewAlloca(vec_head_typ)

	state.Block.NewStore(
		len_f,
		state.Block.NewGetElementPtr(
			vec_head_typ,
			head,
			I32(0), vectorLenOffset))

	state.Block.NewStore(
		cap_f,
		state.Block.NewGetElementPtr(
			vec_head_typ,
			head,
			I32(0), vectorCapOffset))

	body := state.Block.NewBitCast(state.Block.NewCall(
		state.GetCalloc(),
		I32(app.Argument.Atom.Tuple[0].TypeOf.Vector.WidthInBytes()),
		cap_f), types.NewPointer(vec_elem_typ))

	state.Block.NewCall(state.GetMemcpy(), body, body_a, len_a, constant.NewBool(false))
	//state.PopulateBody(body, vec_elem_typ, vector.Vector)

	// Point the vector header to alloc'd body
	//state.WriteVectorPointer(head, head_type, body)

	return nil
}
