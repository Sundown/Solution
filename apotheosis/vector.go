package apotheosis

import (
	"fmt"

	"github.com/sundown/solution/prism"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
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

	head := env.Block.NewBitCast(env.Block.NewCall(env.getCalloc(), i32(3), i32(8)), types.NewPointer(headType))

	// Store vector length
	env.writeVectorLength(head, leng, headType)

	// Store vector capacity
	env.writeVectorCapacity(head, cap, headType)

	isConstant := true
	for _, elm := range *vector.Body {
		if !prism.IsConstant(elm) {
			isConstant = false
			break
		}
	}

	body := env.buildVectorBody(elmType, cap, vector.Type().(prism.VectorType).Width())

	if len(*vector.Body) > 0 {
		if isConstant {
			env.Block.NewCall(
				env.getMemcpy(),
				env.Block.NewBitCast(body, types.I8Ptr),
				env.Block.NewBitCast(env.populateConstBody(elmType, *vector.Body), types.I8Ptr),
				env.Block.NewMul(i64(int64(len(*vector.Body))),
					i64(vector.Type().(prism.VectorType).Type.Width())),
				constant.NewInt(types.I1, 0))
		} else {
			env.populateBody(body, elmType, *vector.Body)
		}
	}

	env.writeVectorPointer(prism.Value{head, vector.Type()}, body)

	return head
}

func (env *Environment) populateConstBody(elementType types.Type, exprVec []prism.Expression) value.Value {
	accum := make([]constant.Constant, len(exprVec))
	for i, expr := range exprVec {
		accum[i] = env.compileExpression(&expr).(constant.Constant)
	}

	return env.Module.NewGlobalDef(fmt.Sprint(env.newID()), constant.NewArray(types.NewArray(uint64(len(exprVec)), elementType), accum...))
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

func (env *Environment) writeVectorPointer(head prism.Value, body value.Value) value.Value {
	env.Block.NewStore(body, env.Block.NewGetElementPtr(head.Type.Realise(), head.Value, i32(0), vectorBodyOffset))

	return head.Value
}

func (env *Environment) buildVectorBody(typ types.Type, cap int64, width int64) *ir.InstBitCast {
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

func (env *Environment) writeVectorLength(vectorStruct value.Value, len int64, typ types.Type) {
	env.Block.NewStore(
		i32(len),
		env.Block.NewGetElementPtr(
			typ,
			vectorStruct,
			i32(0), vectorLenOffset))
}

func (env *Environment) writeVectorCapacity(vectorStruct value.Value, cap int64, typ types.Type) {
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

func (env *Environment) readVectorSizes(vec prism.Value) (length value.Value, capacity value.Value) {
	return env.readVectorLength(vec), env.readVectorCapacity(vec)
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

func (env *Environment) writeElement(vec prism.Value, elm value.Value, index value.Value) {
	env.Block.NewStore(
		elm,
		env.Block.NewGetElementPtr(
			vec.Type.(prism.VectorType).Type.Realise(),
			env.Block.NewGetElementPtr(vec.Type.Realise(), vec.Value, i32(0), vectorBodyOffset),
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

func (env Environment) vectorFactory(elmType prism.Type, size value.Value) prism.Value {
	head := env.Block.NewBitCast(
		env.Block.NewCall(env.getCalloc(), i32(3), i32(8)),
		types.NewPointer(prism.Vec(elmType).Realise()))

	env.writeLLVectorLength(prism.Val(head, prism.Vec(elmType)), size)
	env.writeLLVectorCapacity(prism.Val(head, prism.Vec(elmType)), size)

	env.writeVectorPointer(
		prism.Val(head, prism.Vec(elmType)),
		env.buildLLVectorBody(elmType, size, elmType.Width()))

	return prism.Val(head, prism.Vec(elmType))
}

func (env Environment) dualVectorFactory(elmType prism.Type, size value.Value) (value.Value, value.Value) {
	head := env.Block.NewBitCast(
		env.Block.NewCall(env.getCalloc(), i32(3), i32(8)),
		types.NewPointer(prism.VectorType{Type: elmType}.Realise()))

	env.writeLLVectorLength(prism.Val(head, prism.Vec(elmType)), size)
	env.writeLLVectorCapacity(prism.Val(head, prism.Vec(elmType)), size)

	body := env.buildLLVectorBody(elmType, size, elmType.Width())
	env.writeVectorPointer(prism.Val(head, prism.Vec(elmType)), body)

	return head, body
}
