package apotheosis

import (
	"github.com/sundown/solution/prism"

	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

// TODO:
// Consider binary/ternary system for reducing, only checking n every so many iterations.
// Create different blocks to handle 1, 2, or 3 operations at a time, branching probably
// makes this slow even with Clang vectorisation

/* Code roughly represents
T Reduce(F, A) {
	T accum = A[-1] // last element
	if len(A) == 1 {
		return accum
	}

	for i := len(A) - 2; i != 0; i-- { // step through backwards
		accum = F(accum, A[i])
 	}

	return accum
}*/

func (env *Environment) compileInlineReduce(fn prism.DyadicFunction, vec prism.Value) value.Value {
	vectyp := vec.Type.(prism.VectorType).Type

	len := env.readVectorLength(vec)
	storeCounter := env.new(env.Block.NewSub(len, i32(2)))

	// Alloc memory for accum, bitcast to ptr of vector element type
	accum := env.Block.NewBitCast(
		value.Value(env.Block.NewCall(env.getCalloc(), i32(1), i32(vectyp.Width()))),
		types.NewPointer(vectyp.Realise()))

	// Load the last element of vector
	e := env.unsafeReadVectorElement(vec, env.Block.NewSub(len, i32(1)))
	if prism.IsVector(vectyp) {
		// Memcpy if subtype is vector
		e = env.Block.NewBitCast(e, types.I8Ptr)
		env.Block.NewCall(
			env.getMemcpy(),
			env.Block.NewBitCast(accum, types.I8Ptr),
			env.Block.NewBitCast(e, types.I8Ptr),
			i64(vectyp.Width()),
			constant.NewBool(false))
	} else {
		// Store if not
		env.Block.NewStore(e, env.Block.NewBitCast(accum, types.NewPointer(vectyp.Realise())))
	}

	loopblock := env.newBlock(env.CurrentFunction)
	exitblock := env.newBlock(env.CurrentFunction)

	// Return immediately if only one element in vector
	env.Block.NewCondBr(env.Block.NewICmp(enum.IPredEQ, len, i32(1)), exitblock, loopblock)

	env.Block = loopblock

	curCounter := loopblock.NewLoad(types.I32, storeCounter)

	// Same idea as accum setup, memcpy if subtype is vector, otherwise store
	if prism.IsVector(vectyp) {
		env.Block.NewCall(
			env.getMemcpy(),
			env.Block.NewBitCast(accum, types.I8Ptr),
			env.Block.NewBitCast(env.apply(fn,
				prism.Val(env.unsafeReadVectorElement(vec, curCounter), vectyp),
				prism.Val(accum, vectyp)), types.I8Ptr),
			i64(vectyp.Width()), constant.NewBool(false))
	} else {
		env.Block.NewStore(
			env.apply(fn,
				prism.Val(env.unsafeReadVectorElement(vec, curCounter), vectyp),
				prism.Val(env.Block.NewLoad(vectyp.Realise(), accum), vectyp)),
			env.Block.NewBitCast(accum, types.NewPointer(vectyp.Realise())))
	}

	// i--
	env.Block.NewStore(env.Block.NewSub(curCounter, i32(1)), storeCounter)

	// if i != 0 { goto loop } else { goto exit }
	env.Block.NewCondBr(env.Block.NewICmp(enum.IPredNE, curCounter, i32(0)), loopblock, exitblock)

	env.Block = exitblock

	// Vectors are always returned as pointers
	if prism.IsVector(vectyp) {
		return accum
	}

	return env.Block.NewLoad(vectyp.Realise(), accum)
}
