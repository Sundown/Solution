package apotheosis

import (
	"fmt"

	"github.com/samber/lo"
	"github.com/sundown/solution/prism"

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

	//head := env.Block.NewCall(env.getCreateVectorHeader(), constant.NewInt(types.I32, leng), constant.NewInt(types.I32, cap), constant.NewInt(types.I32, vector.Type().(prism.VectorType).Type.Width()))
	head := env.vectorFactory(vector.Type().(prism.VectorType).Type, i32(leng))

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

	env.writeVectorPointer(prism.Val(head.Value, vector.Type()), body)

	return head.Value
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
	typ := head.Type.Realise()
	dest := env.Block.NewGetElementPtr(typ, head.Value, i32(0), vectorBodyOffset)
	src := env.Block.NewBitCast(body, types.NewPointer(head.Type.(prism.VectorType).Type.Realise()))
	env.Block.NewStore(src, dest)
	return head.Value
}

func (env *Environment) buildVectorBody(typ types.Type, cap int64, width int64) value.Value {
	return env.Block.NewCall(
		env.getCalloc(),
		i32(width),
		i32(cap))
}

func (env *Environment) buildLLVectorBody(typ prism.Type, cap value.Value, width int64) value.Value {
	return env.Block.NewCall(
		env.getCalloc(),
		i32(width),
		cap)
}

func (env *Environment) writeLLVectorLength(vec prism.Value, len value.Value) {
	typ := vec.Type.Realise()
	dest := env.Block.NewGetElementPtr(typ, vec.Value, i32(0), vectorLenOffset)
	env.Block.NewStore(len, dest)
}

func (env *Environment) writeLLVectorCapacity(vec prism.Value, cap value.Value) {
	typ := vec.Type.Realise()
	dest := env.Block.NewGetElementPtr(typ, vec.Value, i32(0), vectorCapOffset)
	env.Block.NewStore(cap, dest)
}

func (env *Environment) writeLLVectorWidth(vec prism.Value, width value.Value) {
	typ := vec.Type.Realise()
	dest := env.Block.NewGetElementPtr(typ, vec.Value, i32(0), vectorWidthOffset)
	env.Block.NewStore(width, dest)
}

func (env *Environment) readVectorSizes(vec prism.Value) (length value.Value, capacity value.Value) {
	return env.readVectorLength(vec), env.readVectorCapacity(vec)
}

func (env *Environment) readVectorLength(vec prism.Value) value.Value {
	src := env.Block.NewGetElementPtr(vec.Type.Realise(), vec.Value, vectorLenOffset)
	return env.Block.NewLoad(types.I32, src)
}

func (env *Environment) readVectorCapacity(vec prism.Value) value.Value {
	src := env.Block.NewGetElementPtr(vec.Type.Realise(), vec.Value, vectorCapOffset)
	return env.Block.NewLoad(types.I32, src)
}

func (env *Environment) readVectorWidth(vec prism.Value) value.Value {
	src := env.Block.NewGetElementPtr(vec.Type.Realise(), vec.Value, vectorWidthOffset)
	return env.Block.NewLoad(types.I32, src)
}

func (env *Environment) writeElement(vec prism.Value, elm value.Value, index value.Value) {
	if _, ok := vec.Type.(prism.VectorType).Type.(prism.AtomicType); ok {
		env.writeVectorElement(vec, elm, index)
	} else {
		env.writeMatrixElement(vec, elm, index)
	}
}

func (env *Environment) writeVectorElement(vec prism.Value, elm value.Value, index value.Value) {
	env.Block.NewStore(
		elm,
		env.Block.NewGetElementPtr(
			vec.Type.(prism.VectorType).Type.Realise(),
			env.getBodyPointer(vec),
			index))
}

func (env *Environment) writeMatrixElement(vec prism.Value, elm value.Value, index value.Value) {

	offset := env.Block.NewMul(i32(vec.Type.(prism.VectorType).Type.Width()), index)
	bodyPtr := env.Block.NewPtrToInt(env.getBodyPointer(vec), types.I32)
	dest := env.Block.NewAdd(bodyPtr, offset)

	width := constant.NewInt(types.I64, vec.Type.(prism.VectorType).Type.Width())
	env.Block.NewCall(env.getMemcpy(), env.Block.NewIntToPtr(dest, types.NewPointer(types.I8)), env.Block.NewBitCast(elm, types.I8Ptr), width, constant.NewInt(types.I1, 0))
}

func (env *Environment) getBodyPointer(vec prism.Value) value.Value {
	return env.Block.NewLoad(types.NewPointer(vec.Type.(prism.VectorType).Type.Realise()), env.Block.NewGetElementPtr(vec.Type.Realise(), vec.Value, i32(0), vectorBodyOffset))
}

func (env *Environment) readVectorElement(vec prism.Value, index value.Value) value.Value {
	typ := vec.Type.(prism.VectorType)

	env.validateVectorIndex(vec, index)

	elm := env.Block.NewGetElementPtr(
		typ.Type.Realise(),
		env.getBodyPointer(vec),
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

func (env *Environment) calculateCapacity(length value.Value) value.Value {
	return env.Block.NewCall(env.getCalculateCapacity(), length)
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

	env.newPanic(bfalse, "Index %d out of bounds [%d]\n", index, leng)

	bfalse.NewUnreachable()

	env.Block.NewCondBr(
		env.Block.NewICmp(
			enum.IPredSLE,
			leng,
			index),
		bfalse, btrue)

	env.Block = btrue
}

func (env *Environment) newVectorHeader(elmType prism.Type, size value.Value) (head prism.Value) {
	capacity := env.calculateCapacity(size)
	width := constant.NewInt(types.I32, elmType.Width())

	//head := env.Block.NewCall(env.getCreateVectorHeader(), size, capacity, width)

	head.Value = env.Block.NewCall(env.getCalloc(), i32(4), i32(1))
	head.Value = env.Block.NewBitCast(head.Value, types.NewPointer(prism.Vec(elmType).Realise()))
	head.Type = prism.Vec(elmType)
	env.writeLLVectorLength(head, size)
	env.writeLLVectorCapacity(head, capacity)
	env.writeLLVectorWidth(head, width)
	// Body pointer guaranteed to be null due to use of calloc

	return head
}

func (env *Environment) vectorFactory(elmType prism.Type, size value.Value) prism.Value {
	newHead := env.newVectorHeader(elmType, size)

	head := env.Block.NewBitCast(newHead.Value, types.NewPointer(prism.Vec(elmType).Realise()))

	env.writeVectorPointer(
		prism.Val(head, prism.Vec(elmType)),
		env.buildLLVectorBody(elmType, size, elmType.Width()))

	return prism.Val(head, prism.Vec(elmType))
}
