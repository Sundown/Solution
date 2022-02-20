package apotheosis

import (
	"sundown/solution/prism"

	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

func (env *Environment) CompileInlineFoldl(fn prism.Expression, vec Value) value.Value {
	vectyp := vec.Type.(prism.VectorType).Type
	lltyp := vectyp.Realise()

	counter := env.Block.NewAlloca(types.I32)
	len := env.ReadVectorLength(vec)
	env.Block.NewStore(env.Block.NewSub(len, I32(3)), counter)

	accum := env.Block.NewAlloca(fn.Type().Realise())

	env.Block.NewStore(env.Apply(fn,
		Value{
			env.UnsafeReadVectorElement(vec, env.Block.NewSub(len, I32(2))),
			vectyp},
		Value{
			env.UnsafeReadVectorElement(vec, env.Block.NewSub(len, I32(1))),
			vectyp}), accum)

	loopblock := env.CurrentFunction.NewBlock("")
	exitblock := env.CurrentFunction.NewBlock("")

	env.Block.NewCondBr(
		env.Block.NewICmp(enum.IPredEQ, len, I32(2)),
		exitblock,
		loopblock)

	env.Block = loopblock

	lcount := loopblock.NewLoad(types.I32, counter)
	loopblock.NewStore(env.Apply(fn,
		Value{
			env.UnsafeReadVectorElement(vec, lcount),
			vectyp},
		Value{
			loopblock.NewLoad(lltyp, accum),
			vectyp}), accum)

	loopblock.NewStore(loopblock.NewSub(lcount, I32(1)), counter)

	loopblock.NewCondBr(
		loopblock.NewICmp(enum.IPredNE, lcount, I32(0)),
		loopblock,
		exitblock)

	env.Block = exitblock
	return loopblock.NewLoad(lltyp, accum)
}
