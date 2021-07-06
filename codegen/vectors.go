package codegen

import (
	"fmt"
	"math"
	"sundown/sunday/parser"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

func BuildVectorType(typ types.Type) *types.StructType {
	return types.NewStruct(
		types.I64,             // Length
		types.I64,             // Capacity
		types.NewPointer(typ)) // Body
}

func (state *State) BuildVectorHeader(typ types.Type) *ir.InstAlloca {
	return state.block.NewAlloca(BuildVectorType(typ))
}

func (state *State) BuildVectorBody(typ types.Type, cap int64, width int64) *ir.InstBitCast {
	return state.block.NewBitCast(state.block.NewCall(
		state.fns["calloc"],
		constant.NewInt(types.I64, width), // Byte size of elements
		constant.NewInt(types.I64, cap)),  // How much memory to alloc
		types.NewPointer(typ)) // Cast alloc'd memory to typ
}

func (state *State) WriteVectorLength(vector_struct *ir.InstAlloca, len int64, typ types.Type) {
	state.block.NewStore(constant.NewInt(types.I64, len),
		state.block.NewGetElementPtr(
			BuildVectorType(typ), vector_struct, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0)))
}

func (state *State) WriteVectorCapacity(vector_struct *ir.InstAlloca, cap int64, typ types.Type) {
	state.block.NewStore(constant.NewInt(types.I64, cap),
		state.block.NewGetElementPtr(
			BuildVectorType(typ), vector_struct, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 1)))
}

func (state *State) WriteVectorPointer(vector_struct *ir.InstAlloca, typ types.Type, body *ir.InstBitCast) {
	state.block.NewStore(body, state.block.NewGetElementPtr(
		BuildVectorType(typ), vector_struct, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 2)))
}

func (state *State) MakeVector(vector []*parser.Expression) (value.Value, types.Type) {
	element_type := GenPrimaryType(vector[0].Primary)
	fmt.Println(element_type)
	// Round length *up* to the nearest power of 2
	capacity := int64(math.Floor(math.Log2(float64(len(vector))) + 1))
	// No point in tiny vectors
	if capacity < 8 {
		capacity = 8
	}

	// Vector structure
	vector_header := state.BuildVectorHeader(element_type)

	// Allocagte vectory body
	vector_body := state.BuildVectorBody(element_type, capacity, 4)

	// Store vector length
	state.WriteVectorLength(vector_header, int64(len(vector)), element_type)

	// Store vector capacity
	state.WriteVectorCapacity(vector_header, capacity, element_type)

	// Point the vector header to alloc'd body
	state.WriteVectorPointer(vector_header, element_type, vector_body)

	for index, element := range vector {
		v, _ := state.Compile(element)

		state.block.NewStore(v,
			state.block.NewGetElementPtr(
				element_type,
				vector_body,
				constant.NewInt(types.I32, int64(index))))

	}

	return vector_header, element_type
}
