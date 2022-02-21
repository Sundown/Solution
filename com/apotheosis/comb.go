package apotheosis

import (
	"sundown/solution/prism"

	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

// TODO define interface type for callable
// https://hackage.haskell.org/package/base-4.16.0.0/docs/Prelude.html#v:zipWith
func (env Environment) CombineOf(in interface{}, a, b Value) value.Value {
	fn, ok := in.(prism.DyadicFunction)
	var ret_typ prism.Type
	if ok {
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

	// New Vector of type [fn.Returns]
	vectyp := prism.VectorType{Type: ret_typ}
	newvec := env.Block.NewAlloca(vectyp.Realise())

	len := env.ReadVectorLength(a)

	env.WriteLLVectorLength(Value{newvec, vectyp}, len)
	env.WriteLLVectorCapacity(Value{newvec, vectyp}, env.ReadVectorCapacity(a))

	body := env.Block.NewBitCast(
		env.Block.NewCall(env.GetCalloc(),
			I32(ret_typ.Width()),
			len),
		types.NewPointer(ret_typ.Realise()))

	env.Block.NewCondBr(
		env.Block.NewICmp(enum.IPredEQ, len, env.ReadVectorLength(b)),
		loopblock,
		panicblock)

	env.Block = loopblock

	lcount := loopblock.NewLoad(types.I32, counter)
	call := env.Apply(in,
		Value{
			env.UnsafeReadVectorElement(a, lcount),
			a.Type.(prism.VectorType).Type},
		Value{
			env.UnsafeReadVectorElement(b, lcount),
			b.Type.(prism.VectorType).Type})

	loopblock.NewStore(call,
		loopblock.NewGetElementPtr(ret_typ.Realise(), body, lcount))

	loopblock.NewStore(loopblock.NewAdd(lcount, I32(1)), counter)

	env.Block = env.CurrentFunction.NewBlock("")
	loopblock.NewCondBr(
		loopblock.NewICmp(enum.IPredNE, lcount, len),
		loopblock,
		env.Block)

	env.Block.NewStore(body,
		env.Block.NewGetElementPtr(vectyp.Realise(), newvec, I32(0), vectorBodyOffset))

	return newvec
}
