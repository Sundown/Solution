package apotheosis

import (
	"github.com/alecthomas/repr"
	"github.com/sundown/solution/prism"

	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

func (env *Environment) newInlineMap(in prism.Callable, vec prism.Value) value.Value {
	var retType prism.Type

	if fn, ok := in.(prism.MonadicFunction); ok {
		retType = fn.Type()
	} else {
		repr.Println(in.(prism.MonadicCallable).MCallable)
		panic("Unreachable")
	}

	writePred := retType.Kind() != prism.VoidType.ID
	leng := env.readVectorLength(vec)

	var head prism.Value
	if writePred {
		head = env.vectorFactory(retType, leng)
	}

	counterStore := env.new(i32(0))

	loopblock := env.newBlock(env.CurrentFunction)
	env.Block.NewBr(loopblock)
	env.Block = loopblock

	curCounter := loopblock.NewLoad(types.I32, counterStore)

	call := env.apply(in, prism.Value{
		Value: env.unsafeReadVectorElement(vec, curCounter),
		Type:  vec.Type.(prism.VectorType).Type})

	if writePred {
		env.writeElement(head, call, curCounter)
	}

	incr := loopblock.NewAdd(curCounter, i32(1))

	loopblock.NewStore(incr, counterStore)

	env.Block = env.newBlock(env.CurrentFunction)
	loopblock.NewCondBr(loopblock.NewICmp(enum.IPredNE, incr, leng), loopblock, env.Block)

	return head.Value
}
