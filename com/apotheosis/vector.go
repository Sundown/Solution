package apotheosis

import (
	"sundown/solution/prism"

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

func (env *Environment) CompileVector(vector prism.Vector) value.Value {
	leng, cap := CalculateVectorSizes(len(*vector.Body))
	elm_type := vector.Type().(prism.VectorType).Type.Realise()
	head_type := vector.Type().Realise()
	head := env.Block.NewAlloca(head_type)

	// Store vector length
	env.WriteVectorLength(head, leng, head_type)

	// Store vector capacity
	env.WriteVectorCapacity(head, cap, head_type)

	body := env.BuildVectorBody(elm_type, cap, vector.Type().(prism.VectorType).Width())

	if len(*vector.Body) > 0 {
		env.PopulateBody(body, elm_type, *vector.Body)
	}

	env.WriteVectorPointer(head, body, head_type)

	return head
}

// Maps from expression[] to vector in LLVM
func (env *Environment) PopulateBody(
	allocated_body *ir.InstBitCast,
	element_type types.Type,
	expr_vec []prism.Expression) {

	ir_elm_type := expr_vec[0].Type()
	for index, element := range expr_vec {
		v := env.CompileExpression(&element)

		if _, ok := ir_elm_type.(prism.AtomicType); !ok {
			v = env.Block.NewLoad(element_type, v) // TODO might be backwards
		}

		env.Block.NewStore(v,
			env.Block.NewGetElementPtr(
				element_type,
				allocated_body,
				I32(int64(index))))
	}
}

func (env *Environment) WriteVectorPointer(head *ir.InstAlloca, body *ir.InstBitCast, head_type types.Type) {
	pt := env.Block.NewGetElementPtr(
		head_type, head, I32(0), vectorBodyOffset)

	env.Block.NewStore(
		body,
		pt)
}

func (env *Environment) BuildVectorBody(typ types.Type, cap int64, width int64) *ir.InstBitCast {
	return env.Block.NewBitCast(env.Block.NewCall(
		env.GetCalloc(),
		I32(width), // Byte size of elements
		I32(cap)),  // How much memory to alloc
		types.NewPointer(typ)) // Cast alloc'd memory to typ
}

func (env *Environment) WriteVectorLength(vector_struct *ir.InstAlloca, len int64, typ types.Type) {
	env.Block.NewStore(
		I64(len),
		env.Block.NewGetElementPtr(
			typ,
			vector_struct,
			I32(0), vectorLenOffset))
}

func (env *Environment) WriteVectorCapacity(vector_struct *ir.InstAlloca, cap int64, typ types.Type) {
	env.Block.NewStore(
		I64(cap),
		env.Block.NewGetElementPtr(
			typ,
			vector_struct,
			I32(0), vectorCapOffset))
}

func (env *Environment) ReadVectorLength(vec Value) value.Value {
	return env.Block.NewLoad(types.I64,
		env.Block.NewGetElementPtr(
			vec.Type.Realise(),
			vec.Value,
			I32(0), vectorLenOffset))
}

func (env *Environment) ReadVectorCapacity(vec Value) value.Value {
	return env.Block.NewLoad(types.I64,
		env.Block.NewGetElementPtr(
			vec.Type.Realise(),
			vec.Value,
			I32(0), vectorCapOffset))
}

func (env *Environment) ReadVectorElement(vec Value, index value.Value) value.Value {
	if _, ok := vec.Type.(prism.VectorType); !ok {
		prism.Panic(
			prism.CT_Unexpected,
			prism.Yellow("vector"),
			prism.Yellow(vec.Type.String()))
	}

	env.ValidateVectorIndex(vec, index)

	elm := env.Block.NewGetElementPtr(
		vec.Type.(prism.VectorType).Type.Realise(),
		env.Block.NewLoad(
			types.NewPointer(vec.Type.(prism.VectorType).Type.Realise()),
			env.Block.NewGetElementPtr(
				vec.Type.Realise(),
				vec.Value,
				I32(0), vectorBodyOffset)), index)

	if _, ok := vec.Type.(prism.VectorType).Type.(prism.AtomicType); ok {
		return env.Block.NewLoad(vec.Type.(prism.VectorType).Type.Realise(), elm)
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

func (env *Environment) ValidateVectorIndex(vec Value, index value.Value) {
	btrue := env.CurrentFunction.NewBlock("")
	bfalse := env.CurrentFunction.NewBlock("")

	leng := env.Block.NewLoad(types.I64,
		env.Block.NewGetElementPtr(
			vec.Type.Realise(),
			vec.Value,
			I32(0), vectorLenOffset))

	env.LLVMPanic(bfalse, "Panic: index %d out of bounds [%d]\n", index, leng)

	bend := env.CurrentFunction.NewBlock("")
	btrue.NewBr(bend)
	bfalse.NewUnreachable()

	env.Block.NewCondBr(
		env.Block.NewICmp(
			enum.IPredSLE,
			leng,
			index),
		bfalse, btrue)

	env.Block = bend
}
