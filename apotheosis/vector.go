package apotheosis

import (
	"github.com/sundown/solution/prism"

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

func (env *Environment) compileVector(vector prism.Vector) value.Value {
	leng, cap := CalculateVectorSizes(len(*vector.Body))
	elmType := vector.Type().(prism.VectorType).Type.Realise()
	head_type := vector.Type().Realise()
	head := env.Block.NewAlloca(head_type)

	// Store vector length
	env.writeVectorLength(head, leng, head_type)

	// Store vector capacity
	env.writeVectorCapacity(head, cap, head_type)

	body := env.BuildVectorBody(elmType, cap, vector.Type().(prism.VectorType).Width())

	if len(*vector.Body) > 0 {
		env.PopulateBody(body, elmType, *vector.Body)
	}

	env.writeVectorPointer(head, body, head_type)

	return head
}

// Maps from expression[] to vector in LLVM
func (env *Environment) PopulateBody(
	allocated_body *ir.InstBitCast,
	element_type types.Type,
	expr_vec []prism.Expression) {

	ir_elmType := expr_vec[0].Type()
	for index, element := range expr_vec {
		v := env.compileExpression(&element)

		if _, ok := ir_elmType.(prism.AtomicType); !ok {
			v = env.Block.NewLoad(element_type, v)
		}

		env.Block.NewStore(v,
			env.Block.NewGetElementPtr(
				element_type,
				allocated_body,
				I32(int64(index))))
	}
}

func (env *Environment) writeVectorPointer(head *ir.InstAlloca, body *ir.InstBitCast, head_type types.Type) value.Value {
	env.Block.NewStore(body, env.Block.NewGetElementPtr(head_type, head, I32(0), vectorBodyOffset))
	return head
}

func (env *Environment) BuildVectorBody(typ types.Type, cap int64, width int64) *ir.InstBitCast {
	return env.Block.NewBitCast(env.Block.NewCall(
		env.GetCalloc(),
		I32(width), // Byte size of elements
		I32(cap)),  // How much memory to alloc
		types.NewPointer(typ)) // Cast alloc'd memory to typ
}

func (env *Environment) BuildLLVectorBody(typ types.Type, cap value.Value, width int64) *ir.InstBitCast {
	return env.Block.NewBitCast(env.Block.NewCall(
		env.GetCalloc(),
		I32(width), // Byte size of elements
		cap),       // How much memory to alloc
		types.NewPointer(typ)) // Cast alloc'd memory to typ
}

func (env *Environment) writeVectorLength(vector_struct *ir.InstAlloca, len int64, typ types.Type) {
	env.Block.NewStore(
		I32(len),
		env.Block.NewGetElementPtr(
			typ,
			vector_struct,
			I32(0), vectorLenOffset))
}

func (env *Environment) writeVectorCapacity(vector_struct *ir.InstAlloca, cap int64, typ types.Type) {
	env.Block.NewStore(
		I32(cap),
		env.Block.NewGetElementPtr(
			typ,
			vector_struct,
			I32(0), vectorCapOffset))
}

func (env *Environment) writeLLVectorLength(vec Value, len value.Value) {
	env.Block.NewStore(
		len,
		env.Block.NewGetElementPtr(
			vec.Type.Realise(),
			vec.Value,
			I32(0), vectorLenOffset))
}

func (env *Environment) writeLLVectorCapacity(vec Value, cap value.Value) {
	env.Block.NewStore(
		cap,
		env.Block.NewGetElementPtr(
			vec.Type.Realise(),
			vec.Value,
			I32(0), vectorCapOffset))
}

func (env *Environment) readVectorLength(vec Value) value.Value {
	return env.Block.NewLoad(types.I32,
		env.Block.NewGetElementPtr(
			vec.Type.Realise(),
			vec.Value,
			I32(0), vectorLenOffset))
}

func (env *Environment) readVectorCapacity(vec Value) value.Value {
	return env.Block.NewLoad(types.I32,
		env.Block.NewGetElementPtr(
			vec.Type.Realise(),
			vec.Value,
			I32(0), vectorCapOffset))
}

func (env *Environment) readVectorElement(vec Value, index value.Value) value.Value {
	typ := vec.Type.(prism.VectorType)

	env.ValidateVectorIndex(vec, index)

	elm := env.Block.NewGetElementPtr(
		typ.Type.Realise(),
		env.Block.NewLoad(
			types.NewPointer(typ.Type.Realise()),
			env.Block.NewGetElementPtr(
				vec.Type.Realise(),
				vec.Value,
				I32(0), vectorBodyOffset)),
		index)

	if _, ok := typ.Type.(prism.AtomicType); ok {
		return env.Block.NewLoad(typ.Type.Realise(), elm)
	}

	return elm
}

// Add similar function which splits instructions into 2 blocks, 1 for the body address calc and the other for elm calculate and load
func (env *Environment) UnsafereadVectorElement(vec Value, index value.Value) value.Value {
	typ := vec.Type.(prism.VectorType)

	elm := env.Block.NewGetElementPtr(
		typ.Type.Realise(),
		env.Block.NewLoad(
			types.NewPointer(typ.Type.Realise()),
			env.Block.NewGetElementPtr(
				vec.Type.Realise(),
				vec.Value,
				I32(0), vectorBodyOffset)),
		index)

	if _, ok := typ.Type.(prism.AtomicType); ok {
		return env.Block.NewLoad(typ.Type.Realise(), elm)
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

	leng := env.readVectorLength(vec)

	env.LLVMPanic(bfalse, "Panic: index %d out of bounds [%d]\n", index, leng)

	bfalse.NewUnreachable()

	env.Block.NewCondBr(
		env.Block.NewICmp(
			enum.IPredSLE,
			leng,
			index),
		bfalse, btrue)

	env.Block = btrue
}

func (env Environment) LLVectorFactory(elmType prism.Type, size value.Value) (head *ir.InstAlloca, body *ir.InstBitCast) {
	head = env.Block.NewAlloca(prism.VectorType{Type: elmType}.Realise())
	env.writeLLVectorLength(Value{head, prism.VectorType{Type: elmType}}, size)
	env.writeLLVectorCapacity(Value{head, prism.VectorType{Type: elmType}}, size)
	body = env.BuildLLVectorBody(elmType.Realise(), size, elmType.Width())
	env.writeVectorPointer(head, body, prism.VectorType{Type: elmType}.Realise())

	return
}
