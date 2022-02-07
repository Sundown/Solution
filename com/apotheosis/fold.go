package apotheosis

import (
	"sundown/solution/prism"

	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

func (env *Environment) CompileInlineFoldl(fn prism.Expression, vec Value) value.Value {
	lltyp := vec.Type.(prism.VectorType).Type.Realise()

	counter := env.Block.NewAlloca(types.I32)
	env.Block.NewStore(I32(2), counter) // start at third item

	accum := env.Block.NewAlloca(fn.Type().Realise())

	env.Block.NewStore(env.Apply(&fn,
		Value{
			env.UnsafeReadVectorElement(vec, I32(0)),
			vec.Type.(prism.VectorType).Type},
		Value{
			env.UnsafeReadVectorElement(vec, I32(1)),
			vec.Type.(prism.VectorType).Type}), accum)

	loopblock := env.CurrentFunction.NewBlock("")
	env.Block.NewBr(loopblock)
	env.Block = loopblock

	loopblock.NewStore(env.Apply(&fn,
		Value{
			loopblock.NewLoad(lltyp, accum),
			vec.Type.(prism.VectorType).Type},
		Value{
			env.UnsafeReadVectorElement(vec, loopblock.NewLoad(types.I32, counter)),
			vec.Type.(prism.VectorType).Type}), accum)

	cond := loopblock.NewICmp(
		enum.IPredSLT,
		loopblock.NewAdd(loopblock.NewLoad(types.I32, counter), I32(1)),
		env.ReadVectorLength(vec))

	loopblock.NewStore(
		loopblock.NewAdd(loopblock.NewLoad(types.I32, counter), I32(1)),
		counter)

	exitblock := env.CurrentFunction.NewBlock("")

	loopblock.NewCondBr(cond, loopblock, exitblock)

	env.Block = exitblock

	return env.Block.NewLoad(lltyp, accum)
}
