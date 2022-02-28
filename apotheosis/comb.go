package apotheosis

import (
	"github.com/sundown/solution/prism"

	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

// https://hackage.haskell.org/package/base-4.16.0.0/docs/Prelude.html#v:zipWith
func (env Environment) CombineOf(in prism.Callable, a, b prism.Value) value.Value {
	var ret_typ prism.Type
	if fn, ok := in.(prism.DyadicFunction); ok {
		ret_typ = fn.Type()
	} else {
		ret_typ = a.Type.(prism.VectorType).Type
	}

	loopblock := env.CurrentFunction.NewBlock("")
	panicblock := env.CurrentFunction.NewBlock("")

	env.LLVMPanic(panicblock, "Combination: vector range mismatch\x0A\x00") // "...\n\0"
	panicblock.NewUnreachable()

	counter := env.Block.NewAlloca(types.I32)
	env.Block.NewStore(I32(0), counter)

	len := env.readVectorLength(a)

	newvec, body := env.LLVectorFactory(ret_typ, len)

	env.Block.NewCondBr(
		env.Block.NewICmp(enum.IPredEQ, len, env.readVectorLength(b)),
		loopblock,
		panicblock)

	env.Block = loopblock

	lcount := loopblock.NewLoad(types.I32, counter)
	call := env.Apply(in,
		prism.Value{env.UnsafereadVectorElement(a, lcount), a.Type.(prism.VectorType).Type},
		prism.Value{env.UnsafereadVectorElement(b, lcount), b.Type.(prism.VectorType).Type})

	loopblock.NewStore(call, loopblock.NewGetElementPtr(ret_typ.Realise(), body, lcount))

	loopblock.NewStore(loopblock.NewAdd(lcount, I32(1)), counter)

	env.Block = env.CurrentFunction.NewBlock("")
	loopblock.NewCondBr(loopblock.NewICmp(enum.IPredNE, lcount, len), loopblock, env.Block)

	return env.writeVectorPointer(newvec, body, prism.VectorType{Type: ret_typ}.Realise())
}
