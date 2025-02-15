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

	head := env.vectorFactory(vector.Type().(prism.VectorType).Type, i32(leng))

	bodyptr := env.Block.NewGetElementPtr(head.Type.Realise(), head.Value, i32(0), vectorBodyOffset)

	if len(*vector.Body) > 0 {
		// Are all elements const?
		if lo.ContainsBy(*vector.Body, prism.IsConstant) {
			constant := env.Module.NewGlobalDef(
				fmt.Sprint(env.newID()),
				constant.NewArray(
					types.NewArray(uint64(len(*vector.Body)), elmType),
					lo.Map(*vector.Body, func(e prism.Expression, _ int) constant.Constant {
						return env.newExpression(&e).(constant.Constant)
					})...))

			env.Block.NewStore(env.Block.NewBitCast(constant, types.NewPointer(elmType)), bodyptr)
		} else {
			body := env.buildVectorBody(i32(cap), vector.Type().(prism.VectorType).Width())
			for index, element := range *vector.Body {
				v := env.newExpression(&element)

				if _, ok := (*vector.Body)[0].Type().(prism.AtomicType); !ok {
					v = env.Block.NewLoad(elmType, v)
				}

				env.Block.NewStore(v, env.Block.NewGetElementPtr(elmType, body, i32(int64(index))))
			}

			env.Block.NewStore(body, bodyptr)
		}
	}

	return head.Value
}

func (env *Environment) writeVectorPointer(head prism.Value, body value.Value) value.Value {
	typ := head.Type.Realise()
	dest := env.Block.NewGetElementPtr(typ, head.Value, i32(0), vectorBodyOffset)
	src := env.Block.NewBitCast(body, types.NewPointer(head.Type.(prism.VectorType).Type.Realise()))
	env.Block.NewStore(src, dest)
	return head.Value
}

func (env *Environment) buildVectorBody(cap value.Value, width int64) value.Value {
	return env.Block.NewCall(
		env.getCalloc(),
		i32(width),
		cap)
}

func (env *Environment) writeVectorLength(vec prism.Value, len value.Value) {
	typ := vec.Type.Realise()
	dest := env.Block.NewGetElementPtr(typ, vec.Value, i32(0), vectorLenOffset)
	env.Block.NewStore(len, dest)
}

func (env *Environment) writeVectorCapacity(vec prism.Value, cap value.Value) {
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

// /1⊃
func (env *Environment) writeMatrixElement(vec prism.Value, elm value.Value, index value.Value) {
	//offset := env.Block.NewMul(i32(32), index)
	bodyPtr := env.getBodyPointer(vec)
	//dest := env.Block.NewAdd(bodyPtr, offset)

	width := constant.NewInt(types.I64, vec.Type.(prism.VectorType).Type.Width()) // usually 32
	env.Block.NewCall(
		env.getMemcpy(),
		bodyPtr,
		env.Block.NewBitCast(elm, types.I8Ptr),
		width,
		constant.NewInt(types.I1, 0))
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

	head.Value = env.Block.NewCall(env.getCalloc(), i32(8), i32(4))
	head.Value = env.Block.NewBitCast(head.Value, types.NewPointer(prism.Vec(elmType).Realise()))
	head.Type = prism.Vec(elmType)
	env.writeVectorLength(head, size)
	env.writeVectorCapacity(head, capacity)
	env.writeLLVectorWidth(head, width)
	// Body pointer guaranteed to be null due to use of calloc

	return head
}

func (env *Environment) vectorFactory(elmType prism.Type, size value.Value) prism.Value {
	newHead := env.newVectorHeader(elmType, size)

	head := env.Block.NewBitCast(newHead.Value, types.NewPointer(prism.Vec(elmType).Realise()))

	env.writeVectorPointer(
		prism.Val(head, prism.Vec(elmType)),
		env.buildVectorBody(size, elmType.Width()))

	return prism.Val(head, prism.Vec(elmType))
}
