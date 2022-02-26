package apotheosis

import (
	"github.com/sundown/solution/prism"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
)

func (env *Environment) CompileInlineMap(fn prism.Expression, vec Value) (head *ir.InstAlloca) {
	write_pred := fn.Type().Kind() != prism.VoidType.ID
	leng := env.ReadVectorLength(vec)
	var body *ir.InstBitCast

	if write_pred {
		head, body = env.LLVectorFactory(fn.Type(), leng)
	}

	counter_store := env.New(I32(0))

	loopblock := env.CurrentFunction.NewBlock("")
	env.Block.NewBr(loopblock)
	env.Block = loopblock

	cur_counter := loopblock.NewLoad(types.I32, counter_store)

	call := env.Apply(fn, Value{
		env.UnsafeReadVectorElement(vec, cur_counter),
		vec.Type.(prism.VectorType).Type})

	if write_pred {
		loopblock.NewStore(call, loopblock.NewGetElementPtr(fn.Type().Realise(), body, cur_counter))
	}

	incr := loopblock.NewAdd(cur_counter, I32(1))

	loopblock.NewStore(incr, counter_store)

	env.Block = env.CurrentFunction.NewBlock("")
	loopblock.NewCondBr(loopblock.NewICmp(enum.IPredNE, incr, leng), loopblock, env.Block)

	return
}
