package apotheosis

import (
	"github.com/sundown/solution/prism"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

var (
	vectorLenOffset  = i32(0)
	vectorCapOffset  = i32(1)
	vectorBodyOffset = i32(2)
)

func (env *Environment) compileVector(vector prism.Vector) value.Value {
	leng, cap := calculateVectorSizes(len(*vector.Body))
	elmType := vector.Type().(prism.VectorType).Type.Realise()
	headType := vector.Type().Realise()
	head := env.Block.NewAlloca(headType)

	// Store vector length
	env.writeVectorLength(head, leng, headType)

	// Store vector capacity
	env.writeVectorCapacity(head, cap, headType)

	body := env.buildVectorBody(elmType, cap, vector.Type().(prism.VectorType).Width())

	if len(*vector.Body) > 0 {
		env.populateBody(body, elmType, *vector.Body)
	}

	env.writeVectorPointer(head, body, headType)

	return head
}

// Maps from expression[] to vector in LLVM
func (env *Environment) populateBody(
	allocatedBody *ir.InstBitCast,
	elementType types.Type,
	exprVec []prism.Expression) {

	irElmType := exprVec[0].Type()
	for index, element := range exprVec {
		v := env.compileExpression(&element)

		if _, ok := irElmType.(prism.AtomicType); !ok {
			v = env.Block.NewLoad(elementType, v)
		}

		env.Block.NewStore(v,
			env.Block.NewGetElementPtr(
				elementType,
				allocatedBody,
				i32(int64(index))))
	}
}

func (env *Environment) writeVectorPointer(head *ir.InstAlloca, body *ir.InstBitCast, headType types.Type) value.Value {
	env.Block.NewStore(body, env.Block.NewGetElementPtr(headType, head, i32(0), vectorBodyOffset))
	return head
}

func (env *Environment) buildVectorBody(typ types.Type, cap int64, width int64) *ir.InstBitCast {
	return env.Block.NewBitCast(env.Block.NewCall(
		env.getCalloc(),
		i32(width), // Byte size of elements
		i32(cap)),  // How much memory to alloc
		types.NewPointer(typ)) // Cast alloc'd memory to typ
}

func (env *Environment) buildLLVectorBody(typ types.Type, cap value.Value, width int64) *ir.InstBitCast {
	return env.Block.NewBitCast(env.Block.NewCall(
		env.getCalloc(),
		i32(width), // Byte size of elements
		cap),       // How much memory to alloc
		types.NewPointer(typ)) // Cast alloc'd memory to typ
}

func (env *Environment) writeVectorLength(vectorStruct *ir.InstAlloca, len int64, typ types.Type) {
	env.Block.NewStore(
		i32(len),
		env.Block.NewGetElementPtr(
			typ,
			vectorStruct,
			i32(0), vectorLenOffset))
}

func (env *Environment) writeVectorCapacity(vectorStruct *ir.InstAlloca, cap int64, typ types.Type) {
	env.Block.NewStore(
		i32(cap),
		env.Block.NewGetElementPtr(
			typ,
			vectorStruct,
			i32(0), vectorCapOffset))
}

func (env *Environment) writeLLVectorLength(vec prism.Value, len value.Value) {
	env.Block.NewStore(
		len,
		env.Block.NewGetElementPtr(
			vec.Type.Realise(),
			vec.Value,
			i32(0), vectorLenOffset))
}

func (env *Environment) writeLLVectorCapacity(vec prism.Value, cap value.Value) {
	env.Block.NewStore(
		cap,
		env.Block.NewGetElementPtr(
			vec.Type.Realise(),
			vec.Value,
			i32(0), vectorCapOffset))
}

func (env *Environment) readVectorLength(vec prism.Value) value.Value {

	return env.Block.NewLoad(types.I32,
		env.Block.NewGetElementPtr(
			vec.Type.Realise(),
			vec.Value,
			i32(0), vectorLenOffset))
}

func (env *Environment) readVectorCapacity(vec prism.Value) value.Value {
	return env.Block.NewLoad(types.I32,
		env.Block.NewGetElementPtr(
			vec.Type.Realise(),
			vec.Value,
			i32(0), vectorCapOffset))
}

func (env *Environment) readVectorElement(vec prism.Value, index value.Value) value.Value {
	typ := vec.Type.(prism.VectorType)

	env.validateVectorIndex(vec, index)

	elm := env.Block.NewGetElementPtr(
		typ.Type.Realise(),
		env.Block.NewLoad(
			types.NewPointer(typ.Type.Realise()),
			env.Block.NewGetElementPtr(
				vec.Type.Realise(),
				vec.Value,
				i32(0), vectorBodyOffset)),
		index)

	if _, ok := typ.Type.(prism.AtomicType); ok {
		return env.Block.NewLoad(typ.Type.Realise(), elm)
	}

	return elm
}

// Add similar function which splits instructions into 2 blocks, 1 for the body address calc and the other for elm calculate and load
func (env *Environment) unsafeReadVectorElement(vec prism.Value, index value.Value) value.Value {
	typ := vec.Type.(prism.VectorType)

	elm := env.Block.NewGetElementPtr(
		typ.Type.Realise(),
		env.Block.NewLoad(
			types.NewPointer(typ.Type.Realise()),
			env.Block.NewGetElementPtr(
				vec.Type.Realise(),
				vec.Value,
				i32(0), vectorBodyOffset)),
		index)

	if _, ok := typ.Type.(prism.AtomicType); ok {
		return env.Block.NewLoad(typ.Type.Realise(), elm)
	}

	return elm
}

func calculateVectorSizes(l int) (leng int64, cap int64) {
	leng = int64(l)
	if leng < 4 {
		cap = 8
	} else {
		cap = 2 * leng
	}

	return leng, cap
}

func (env *Environment) validateVectorIndex(vec prism.Value, index value.Value) {
	btrue := env.newBlock(env.CurrentFunction)
	bfalse := env.newBlock(env.CurrentFunction)

	leng := env.readVectorLength(vec)

	env.compilePanic(bfalse, "Panic: index %d out of bounds [%d]\n", index, leng)

	bfalse.NewUnreachable()

	env.Block.NewCondBr(
		env.Block.NewICmp(
			enum.IPredSLE,
			leng,
			index),
		bfalse, btrue)

	env.Block = btrue
}

func (env Environment) vectorFactory(elmType prism.Type, size value.Value) (head *ir.InstAlloca, body *ir.InstBitCast) {
	head = env.Block.NewAlloca(prism.VectorType{Type: elmType}.Realise())
	env.writeLLVectorLength(prism.Value{Value: head, Type: prism.VectorType{Type: elmType}}, size)
	env.writeLLVectorCapacity(prism.Value{Value: head, Type: prism.VectorType{Type: elmType}}, size)
	body = env.buildLLVectorBody(elmType.Realise(), size, elmType.Width())
	env.writeVectorPointer(head, body, prism.VectorType{Type: elmType}.Realise())

	return
}
