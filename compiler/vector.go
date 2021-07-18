package compiler

import (
	"math"
	"sundown/sunday/parse"

	"github.com/alecthomas/repr"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

func (state *State) CompileVector(vector *parse.Atom) value.Value {
	if vector.Vector == nil {
		panic("Unreachable")
	}

	leng, cap := CalculateVectorSizes(vector.Vector)

	elm_type := vector.Vector[0].TypeOf.AsLLType()
	head_type := vector.TypeOf.AsLLType()

	head := state.BuildVectorHeader(head_type)
	// Store vector length
	state.WriteVectorLength(head, leng, head_type)

	// Store vector capacity
	state.WriteVectorCapacity(head, cap, head_type)

	body := state.BuildVectorBody(elm_type, cap, 8)
	repr.Println(body)
	// Point the vector header to alloc'd body
	//state.WriteVectorPointer(head, elm_type, body)

	state.WriteVectorPointer(head, head_type, body)

	return head
}

func (state *State) WriteVectorPointer(vector_struct *ir.InstAlloca, typ types.Type, body *ir.InstBitCast) {
	state.Block.NewStore(
		body,
		state.Block.NewGetElementPtr(typ, vector_struct, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 2)))
}

func (state *State) BuildVectorHeader(typ types.Type) *ir.InstAlloca {
	return state.Block.NewAlloca(typ)
}

func (state *State) BuildVectorBody(typ types.Type, cap int64, width int64) *ir.InstBitCast {
	return state.Block.NewBitCast(state.Block.NewCall(
		state.Calloc(),
		constant.NewInt(types.I32, width), // Byte size of elements
		constant.NewInt(types.I32, cap)),  // How much memory to alloc
		types.NewPointer(typ)) // Cast alloc'd memory to typ
}

func (state *State) WriteVectorLength(vector_struct *ir.InstAlloca, len int64, typ types.Type) {
	state.Block.NewStore(
		constant.NewInt(types.I32, len),
		state.Block.NewGetElementPtr(typ, vector_struct, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0)))
}

func (state *State) WriteVectorCapacity(vector_struct *ir.InstAlloca, cap int64, typ types.Type) {
	state.Block.NewStore(
		constant.NewInt(types.I32, cap),
		state.Block.NewGetElementPtr(typ, vector_struct, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 1)))
}

func CalculateVectorSizes(vector []*parse.Expression) (leng int64, cap int64) {
	leng = int64(len(vector))
	if leng == 0 {
		cap = 8
	} else {
		// Round upto the next power of 2
		// TODO: broken lol
		cap = int64(math.Floor(math.Log2(float64(leng)) + 1))
	}

	return leng, cap
}
