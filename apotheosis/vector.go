package apotheosis

import (
	"fmt"

	"github.com/samber/lo"
	"github.com/sundown/solution/prism"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

var (
	vectorLenOffset   = i32(0)
	vectorCapOffset   = i32(1)
	vectorWidthOffset = i32(2)
	vectorBodyOffset  = i32(3)
)

// newVector maps from prism.Vector to Value contaning LLVM vector.
func (env *Environment) newVector(vector prism.Vector) value.Value {
	leng, cap := calculateVectorSizes(len(*vector.Body))
	elmType := vector.Type().(prism.VectorType).Type.Realise()

	head := env.Block.NewCall(env.getCreateVectorHeader(), constant.NewInt(types.I32, leng), constant.NewInt(types.I32, cap), constant.NewInt(types.I32, vector.Type().(prism.VectorType).Type.Width()))
	//head = env.Block.NewBitCast(head, headType)
	// Perform calloc for body and place width and capacity and let LLVM know the type.
	//
	body := env.buildVectorBody(elmType, cap, vector.Type().(prism.VectorType).Width())

	if len(*vector.Body) > 0 {
		// Are all elements const?
		if lo.ContainsBy(*vector.Body, prism.IsConstant) {

			allocSize := env.Block.NewMul(i64(int64(len(*vector.Body))),
				i64(vector.Type().(prism.VectorType).Type.Width()))

			env.Block.NewCall(
				env.getMemcpy(),
				body, //env.Block.NewBitCast(body, types.I8Ptr), // TODO maybe should just be body, try later
				env.populateConstBody(elmType, *vector.Body),
				allocSize,
				constant.NewInt(types.I1, 0))
		} else {
			env.populateBody(body, elmType, *vector.Body)
		}
	}

	env.writeVectorPointer(prism.Val(head, vector.Type()), body)

	return head
}

func (env *Environment) populateConstBody(elementType types.Type, exprVec []prism.Expression) value.Value {
	accum := make([]constant.Constant, len(exprVec))
	for i, expr := range exprVec {
		accum[i] = env.newExpression(&expr).(constant.Constant)
	}

	return env.Module.NewGlobalDef(
		fmt.Sprint(env.newID()),
		constant.NewArray(types.NewArray(uint64(len(exprVec)), elementType), accum...))
}

// Maps from expression[] to vector in LLVM
func (env *Environment) populateBody(body value.Value, elmType types.Type, exprVec []prism.Expression) {
	for index, element := range exprVec {
		v := env.newExpression(&element)

		if _, ok := exprVec[0].Type().(prism.AtomicType); !ok {
			v = env.Block.NewLoad(elmType, v)
		}

		env.Block.NewStore(v, env.Block.NewGetElementPtr(elmType, body, i32(int64(index))))
	}
}

func (env *Environment) writeVectorPointer(head prism.Value, body value.Value) value.Value {
	env.Block.NewCall(env.getWriteVectorPointer(), head.Value, body)
	return head.Value
}

func (env *Environment) buildVectorBody(typ types.Type, cap int64, width int64) value.Value {
	return env.Block.NewBitCast(env.Block.NewCall(
		env.getCalloc(),
		i32(width),
		i32(cap)),
		types.NewPointer(typ))
}

func (env *Environment) buildLLVectorBody(typ prism.Type, cap value.Value, width int64) *ir.InstBitCast {
	return env.Block.NewBitCast(env.Block.NewCall(
		env.getCalloc(),
		i32(width),
		cap),
		types.NewPointer(typ.Realise()))
}

func (env *Environment) writeVectorLength(vectorStruct value.Value, len int64) {
	env.Block.NewCall(env.getWriteVectorLength(), vectorStruct, i32(len))
}

func (env *Environment) writeVectorCapacity(vectorStruct value.Value, cap int64) {
	env.Block.NewCall(env.getWriteVectorCapacity(), vectorStruct, i32(cap))
}

func (env *Environment) writeLLVectorLength(vec prism.Value, len value.Value) {
	env.Block.NewCall(env.getWriteVectorLength(), vec.Value, len)
}

func (env *Environment) writeLLVectorCapacity(vec prism.Value, cap value.Value) {
	env.Block.NewCall(env.getWriteVectorCapacity(), vec.Value, cap)
}

func (env *Environment) readVectorSizes(vec prism.Value) (length value.Value, capacity value.Value) {
	return env.readVectorLength(vec), env.readVectorCapacity(vec)
}

func (env *Environment) readVectorLength(vec prism.Value) value.Value {
	return env.Block.NewCall(env.getReadVectorLength(), vec.Value)
}

func (env *Environment) readVectorCapacity(vec prism.Value) value.Value {
	return env.Block.NewCall(env.getReadVectorCapacity(), vec.Value)
}

func (env *Environment) writeElement(vec prism.Value, elm value.Value, index value.Value) {
	env.Block.NewStore(
		elm,
		env.Block.NewGetElementPtr(
			vec.Type.(prism.VectorType).Type.Realise(),
			env.Block.NewLoad(types.NewPointer(vec.Type.(prism.VectorType).Type.Realise()), env.Block.NewGetElementPtr(vec.Type.Realise(), vec.Value, i32(0), vectorBodyOffset)),
			index))
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

func (env *Environment) unsafeReadVectorElement(vec prism.Value, index value.Value) value.Value {
	typ := vec.Type.(prism.VectorType).Type

	elm := env.Block.NewGetElementPtr(
		typ.Realise(),
		env.Block.NewLoad(
			types.NewPointer(typ.Realise()),
			env.Block.NewGetElementPtr(
				vec.Type.Realise(),
				vec.Value,
				i32(0), vectorBodyOffset)),
		index)

	if _, ok := typ.(prism.AtomicType); ok {
		return env.Block.NewLoad(typ.Realise(), elm)
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

	env.newPanic(bfalse, "RUBICON: index %d out of bounds [%d]\n", index, leng)

	bfalse.NewUnreachable()

	env.Block.NewCondBr(
		env.Block.NewICmp(
			enum.IPredSLE,
			leng,
			index),
		bfalse, btrue)

	env.Block = btrue
}

func (env *Environment) vectorFactory(elmType prism.Type, size value.Value) prism.Value {
	head := env.Block.NewCall(env.getCreateVectorHeader(), constant.NewInt(types.I32, 0), size, constant.NewInt(types.I32, elmType.Width()))

	env.writeVectorPointer(
		prism.Val(head, prism.Vec(elmType)),
		env.buildLLVectorBody(elmType, size, elmType.Width()))

	return prism.Val(head, prism.Vec(elmType))
}

func (env *Environment) dualVectorFactory(elmType prism.Type, size value.Value) (value.Value, value.Value) {
	head := env.Block.NewBitCast(
		env.Block.NewCall(env.getCalloc(), i32(3), i32(8)),
		types.NewPointer(prism.VectorType{Type: elmType}.Realise()))

	env.writeLLVectorLength(prism.Val(head, prism.Vec(elmType)), size)
	env.writeLLVectorCapacity(prism.Val(head, prism.Vec(elmType)), size)

	body := env.buildLLVectorBody(elmType, size, elmType.Width())
	env.writeVectorPointer(prism.Val(head, prism.Vec(elmType)), body)

	return head, body
}
