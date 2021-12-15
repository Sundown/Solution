package apotheosis

import (
	"sundown/solution/oversight"
	"sundown/solution/subtle"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

var (
	vectorLenOffset  = I32(0)
	vectorCapOffset  = I32(1)
	vectorBodyOffset = I32(2)
)

func (state *State) CompileVector(vector *subtle.Morpheme) value.Value {
	if vector.Vector == nil {
		panic("Unreachable")
	}

	leng, cap := CalculateVectorSizes(len(vector.Vector))
	elm_type := vector.TypeOf.Vector.AsLLType()
	head_type := vector.TypeOf.AsLLType()
	head := state.Block.NewAlloca(head_type)

	// Store vector length
	state.WriteVectorLength(head, leng, head_type)

	// Store vector capacity
	state.WriteVectorCapacity(head, cap, head_type)

	body := state.BuildVectorBody(elm_type, cap, vector.TypeOf.Vector.WidthInBytes())

	if len(vector.Vector) > 0 {
		state.PopulateBody(body, elm_type, vector.Vector)
	}

	// Probably broken for 0-length vectors because this pointer is gonna be garbage
	state.WriteVectorPointer(head, head_type, body)

	return head
}

// Maps from expression[] to vector in LLVM
func (state *State) PopulateBody(
	allocated_body *ir.InstBitCast,
	element_type types.Type,
	expr_vec []*subtle.Expression) {

	ir_elm_type := expr_vec[0].TypeOf
	for index, element := range expr_vec {
		v := state.CompileExpression(element)

		// I'm 50% sure this is wrong
		if ir_elm_type.Atomic == nil {
			v = state.Block.NewLoad(element_type, v)
		}

		state.Block.NewStore(v,
			state.Block.NewGetElementPtr(
				element_type,
				allocated_body,
				I32(int64(index))))
	}
}

func (state *State) WriteVectorPointer(
	vector_header *ir.InstAlloca,
	vector_header_type types.Type,
	constructed_body *ir.InstBitCast) {
	state.Block.NewStore(
		constructed_body,
		state.Block.NewGetElementPtr(
			vector_header_type, vector_header, I32(0), vectorBodyOffset))
}

func (state *State) BuildVectorBody(typ types.Type, cap int64, width int64) *ir.InstBitCast {
	return state.Block.NewBitCast(state.Block.NewCall(
		state.GetCalloc(),
		I32(width), // Byte size of elements
		I32(cap)),  // How much memory to alloc
		types.NewPointer(typ)) // Cast alloc'd memory to typ
}

func (state *State) WriteVectorLength(vector_struct *ir.InstAlloca, len int64, typ types.Type) {
	state.Block.NewStore(
		I64(len),
		state.Block.NewGetElementPtr(
			typ,
			vector_struct,
			I32(0), vectorLenOffset))
}

func (state *State) WriteVectorCapacity(vector_struct *ir.InstAlloca, cap int64, typ types.Type) {
	state.Block.NewStore(
		I64(cap),
		state.Block.NewGetElementPtr(
			typ,
			vector_struct,
			I32(0), vectorCapOffset))
}

func (state *State) ReadVectorLength(typ *subtle.Type, vec value.Value) value.Value {
	return state.Block.NewLoad(types.I64,
		state.Block.NewGetElementPtr(
			typ.AsLLType(),
			vec,
			I32(0), vectorLenOffset))
}

func (state *State) ReadVectorCapacity(typ *subtle.Type, vec value.Value) value.Value {
	return state.Block.NewLoad(types.I64,
		state.Block.NewGetElementPtr(
			typ.AsLLType(),
			vec,
			I32(0), vectorCapOffset))
}

func (state *State) ReadVectorElement(head_typ *subtle.Type, src value.Value, index value.Value) value.Value {
	if head_typ.Vector == nil {
		oversight.Panic(
			oversight.CT_Unexpected,
			oversight.Yellow("vector"),
			oversight.Yellow(head_typ.String()))
	}

	state.ValidateVectorIndex(head_typ, src, index)

	elm := state.Block.NewGetElementPtr(
		head_typ.Vector.AsLLType(), state.Block.NewLoad(
			types.NewPointer(head_typ.Vector.AsLLType()),
			state.Block.NewGetElementPtr(
				head_typ.AsLLType(),
				src,
				I32(0), vectorBodyOffset)), index)

	if head_typ.Vector.Atomic != nil {
		return state.Block.NewLoad(head_typ.Vector.AsLLType(), elm)
	}

	return elm
}

func CalculateVectorSizes(l int) (leng int64, cap int64) {
	leng = int64(l)
	if leng < 4 {
		cap = 8
	} else {
		cap = 2 * leng
	}

	return leng, cap
}

func (state *State) ValidateVectorIndex(typ *subtle.Type, src value.Value, index value.Value) {
	btrue := state.CurrentFunction.NewBlock("")
	bfalse := state.CurrentFunction.NewBlock("")

	leng := state.Block.NewLoad(types.I64,
		state.Block.NewGetElementPtr(
			typ.AsLLType(),
			src,
			I32(0), vectorLenOffset))

	state.LLVMPanic(bfalse, "Panic: index %d out of bounds [%d]\n", index, leng)

	bend := state.CurrentFunction.NewBlock("")
	btrue.NewBr(bend)
	bfalse.NewUnreachable()

	state.Block.NewCondBr(
		state.Block.NewICmp(
			enum.IPredSLE,
			leng,
			index),
		bfalse, btrue)

	state.Block = bend
}
