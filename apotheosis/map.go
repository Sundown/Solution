package apotheosis

import (
	"github.com/sundown/solution/prism"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
)

func (env *Environment) compileInlineMap(fn prism.MonadicFunction, vec prism.Value) (head *ir.InstAlloca) {
	writePred := fn.Type().Kind() != prism.VoidType.ID
	leng := env.readVectorLength(vec)
	var body *ir.InstBitCast

	if writePred {
		head, body = env.vectorFactory(fn.Type(), leng)
	}

	counterStore := env.new(i32(0))

	loopblock := env.newBlock(env.CurrentFunction)
	env.Block.NewBr(loopblock)
	env.Block = loopblock

	curCounter := loopblock.NewLoad(types.I32, counterStore)

	call := env.apply(fn, prism.Value{
		Value: env.unsafeReadVectorElement(vec, curCounter),
		Type:  vec.Type.(prism.VectorType).Type})

	if writePred {
		loopblock.NewStore(call, loopblock.NewGetElementPtr(fn.Type().Realise(), body, curCounter))
	}

	incr := loopblock.NewAdd(curCounter, i32(1))

	loopblock.NewStore(incr, counterStore)

	env.Block = env.newBlock(env.CurrentFunction)
	loopblock.NewCondBr(loopblock.NewICmp(enum.IPredNE, incr, leng), loopblock, env.Block)

	return
}
