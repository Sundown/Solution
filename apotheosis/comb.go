package apotheosis

import (
	"github.com/sundown/solution/prism"

	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

// https://hackage.haskell.org/package/base-4.16.0.0/docs/Prelude.html#v:zipWith
func (env *Environment) combineOf(in prism.Callable, a, b prism.Value) value.Value {
	var retType prism.Type
	if fn, ok := in.(prism.DyadicFunction); ok {
		retType = fn.Type()
	} else {
		retType = a.Type.(prism.VectorType).Type
	}

	loopblock := env.newBlock(env.CurrentFunction)
	panicblock := env.newBlock(env.CurrentFunction)

	env.newPanic(panicblock, "Combination: vector cardinality not equal\x0A\x00")
	panicblock.NewUnreachable()

	counter := env.Block.NewAlloca(types.I32)
	env.Block.NewStore(i32(0), counter)

	len := env.readVectorLength(a)

	newvec := env.vectorFactory(retType, len)

	blen := env.readVectorLength(b)

	env.Block.NewCondBr(
		env.Block.NewICmp(enum.IPredEQ, len, blen),
		loopblock, panicblock)

	env.Block = loopblock

	lcount := loopblock.NewLoad(types.I32, counter)
	call := env.apply(in,
		prism.Val(env.unsafeReadVectorElement(a, lcount), a.Type.(prism.VectorType).Type),
		prism.Val(env.unsafeReadVectorElement(b, lcount), b.Type.(prism.VectorType).Type))

	env.writeElement(newvec, call, lcount)

	loopblock.NewStore(loopblock.NewAdd(lcount, i32(1)), counter)

	env.Block = env.newBlock(env.CurrentFunction)
	loopblock.NewCondBr(loopblock.NewICmp(enum.IPredNE, lcount, len), loopblock, env.Block)

	return newvec.Value
}
