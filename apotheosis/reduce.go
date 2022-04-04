package apotheosis

import (
	"github.com/sundown/solution/prism"

	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

func (env *Environment) compileInlineReduce(fn prism.DyadicFunction, vec prism.Value) value.Value {
	vectyp := vec.Type.(prism.VectorType).Type

	len := env.readVectorLength(vec)
	counter := env.new(env.Block.NewSub(len, i32(3)))

	accum := env.new(env.apply(fn,
		prism.Value{Value: env.unsafeReadVectorElement(vec, env.Block.NewSub(len, i32(2))), Type: vectyp},
		prism.Value{Value: env.unsafeReadVectorElement(vec, env.Block.NewSub(len, i32(1))), Type: vectyp}))

	loopblock := env.newBlock(env.CurrentFunction)
	exitblock := env.newBlock(env.CurrentFunction)

	env.Block.NewCondBr(env.Block.NewICmp(enum.IPredEQ, len, i32(2)), exitblock, loopblock)

	env.Block = loopblock

	lcount := loopblock.NewLoad(types.I32, counter)
	lcount.SetName("counter_load")

	loopblock.NewStore(
		env.apply(fn,
			prism.Value{Value: env.unsafeReadVectorElement(vec, lcount), Type: vectyp},
			prism.Value{Value: loopblock.NewLoad(vectyp.Realise(), accum), Type: vectyp}),
		accum)

	loopblock.NewStore(loopblock.NewSub(lcount, i32(1)), counter)

	loopblock.NewCondBr(loopblock.NewICmp(enum.IPredNE, lcount, i32(0)), loopblock, exitblock)

	env.Block = exitblock
	return env.Block.NewLoad(vectyp.Realise(), accum)
}
