package apotheosis

import (
	"github.com/sundown/solution/prism"

	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

func (env *Environment) compileInlineReduce(fn prism.DyadicFunction, vec prism.Value) value.Value {
	vectyp := vec.Type.(prism.VectorType).Type

	len := env.readVectorLength(vec)
	counter := env.new(env.Block.NewSub(len, i32(2)))

	accum := env.Block.NewBitCast(value.Value(env.Block.NewCall(env.getCalloc(), i32(1), i32(vectyp.Width()))), types.NewPointer(vectyp.Realise()))

	e := env.unsafeReadVectorElement(vec, env.Block.NewSub(len, i32(1)))
	if prism.IsVector(vectyp) {
		e = env.Block.NewBitCast(e, types.I8Ptr)
		env.Block.NewCall(env.getMemcpy(), env.Block.NewBitCast(accum, types.I8Ptr), env.Block.NewBitCast(e, types.I8Ptr), i64(vectyp.Width()), constant.NewBool(false))
	} else {
		env.Block.NewStore(e, env.Block.NewBitCast(accum, types.NewPointer(vectyp.Realise())))
	}

	// TODO urgent: add check for short vectors

	loopblock := env.newBlock(env.CurrentFunction)
	exitblock := env.newBlock(env.CurrentFunction)
	env.Block.NewBr(loopblock)
	env.Block = loopblock

	lcount := loopblock.NewLoad(types.I32, counter)

	if prism.IsVector(vectyp) {
		env.Block.NewCall(
			env.getMemcpy(),
			env.Block.NewBitCast(accum, types.I8Ptr),
			env.Block.NewBitCast(env.apply(fn,
				prism.Value{Value: env.unsafeReadVectorElement(vec, lcount), Type: vectyp},
				prism.Val(accum, vectyp)), types.I8Ptr),
			i64(vectyp.Width()), constant.NewBool(false))
	} else {
		env.Block.NewStore(
			env.apply(fn,
				prism.Value{Value: env.unsafeReadVectorElement(vec, lcount), Type: vectyp},
				prism.Val(env.Block.NewLoad(vectyp.Realise(), accum), vectyp)),
			env.Block.NewBitCast(accum, types.NewPointer(vectyp.Realise())))
	}

	env.Block.NewStore(env.Block.NewSub(lcount, i32(1)), counter)

	env.Block.NewCondBr(env.Block.NewICmp(enum.IPredNE, lcount, i32(0)), loopblock, exitblock)

	env.Block = exitblock

	if prism.IsVector(vectyp) {
		return accum
	}

	return env.Block.NewLoad(vectyp.Realise(), accum)
}
