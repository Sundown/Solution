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
	counter := env.New(env.Block.NewSub(len, I32(3)))

	accum := env.New(env.Apply(fn,
		prism.Value{Value: env.UnsafereadVectorElement(vec, env.Block.NewSub(len, I32(2))), Type: vectyp},
		prism.Value{Value: env.UnsafereadVectorElement(vec, env.Block.NewSub(len, I32(1))), Type: vectyp}))

	loopblock := env.CurrentFunction.NewBlock("")
	exitblock := env.CurrentFunction.NewBlock("")

	env.Block.NewCondBr(env.Block.NewICmp(enum.IPredEQ, len, I32(2)), exitblock, loopblock)

	env.Block = loopblock

	lcount := loopblock.NewLoad(types.I32, counter)
	lcount.SetName("counter_load")

	loopblock.NewStore(
		env.Apply(fn,
			prism.Value{Value: env.UnsafereadVectorElement(vec, lcount), Type: vectyp},
			prism.Value{Value: loopblock.NewLoad(vectyp.Realise(), accum), Type: vectyp}),
		accum)

	loopblock.NewStore(loopblock.NewSub(lcount, I32(1)), counter)

	loopblock.NewCondBr(loopblock.NewICmp(enum.IPredNE, lcount, I32(0)), loopblock, exitblock)

	env.Block = exitblock
	return env.Block.NewLoad(vectyp.Realise(), accum)
}
