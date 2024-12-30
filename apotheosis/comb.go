package apotheosis

import (
	"github.com/sundown/solution/prism"

	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

// Often called "zipwith"
func (env *Environment) combineOf(in prism.Callable, a, b prism.Value) value.Value {
	var retType prism.Type

	if fn, ok := in.(prism.DyadicFunction); ok {
		retType = fn.Type()
	} else {
		retType = a.Type.(prism.VectorType).Type
	}

	counterStore := env.new(i32(0))

	aLeng := env.readVectorLength(a)

	resultant := env.vectorFactory(retType, aLeng)

	loopblock := env.newNamedBlock(env.CurrentFunction, "loopblock")

	// Panic if the vectors are not the same length
	panicblock := env.newNamedBlock(env.CurrentFunction, "panicblock")

	env.newPanic(panicblock, "Combination: vector cardinality not equal\x0A\x00")

	panicblock.NewUnreachable()

	bLeng := env.readVectorLength(b)

	env.Block.NewCondBr(
		env.Block.NewICmp(enum.IPredEQ, aLeng, bLeng),
		loopblock, panicblock)
	// End panic block

	env.Block = loopblock

	lcount := env.Block.NewLoad(types.I32, counterStore)
	call := env.apply(in,
		prism.Val(env.unsafeReadVectorElement(a, lcount), a.Type.(prism.VectorType).Type),
		prism.Val(env.unsafeReadVectorElement(b, lcount), b.Type.(prism.VectorType).Type))

	env.writeElement(resultant, call, lcount)

	incr := env.Block.NewAdd(lcount, i32(1))
	env.Block.NewStore(incr, counterStore)

	endBlock := env.newBlock(env.CurrentFunction)

	env.Block.NewCondBr(env.Block.NewICmp(enum.IPredNE, incr, aLeng), env.Block, endBlock)

	env.Block = endBlock
	return resultant.Value
}
