package apotheosis

import (
	"github.com/sundown/solution/prism"

	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

func (env *Environment) compileInlineIota(val prism.Value) value.Value {
	head, body := env.vectorFactory(val.Type, env.Block.NewTrunc(env.castInt(val), types.I32))

	counterStore := env.new(i64(0))

	loopblock := env.newBlock(env.CurrentFunction)
	env.Block.NewBr(loopblock)
	env.Block = loopblock

	curCounter := loopblock.NewLoad(types.I64, counterStore)

	ptr := loopblock.NewGetElementPtr(val.Type.Realise(), body, curCounter)

	incr := loopblock.NewAdd(curCounter, i64(1))

	loopblock.NewStore(incr, ptr)

	loopblock.NewStore(incr, counterStore)

	env.Block = env.newBlock(env.CurrentFunction)
	loopblock.NewCondBr(loopblock.NewICmp(enum.IPredNE, incr, val.Value), loopblock, env.Block)

	return head
}

func (env *Environment) compileInlineEnclose(val prism.Value) value.Value {
	head, body := env.vectorFactory(val.Type, i32(1))
	env.Block.NewStore(val.Value, env.Block.NewGetElementPtr(val.Type.Realise(), body, i32(0)))
	return head
}

func (env *Environment) invokePrint(val prism.Value, end string) value.Value {
	if val.Type.Equals(prism.StringType) {
		return env.Block.NewCall(
			env.getPrintf(),
			env.getFormatString(val.Type, end),
			env.Block.NewLoad(types.I8Ptr, env.Block.NewGetElementPtr(
				val.Type.Realise(),
				val.Value,
				i32(0), vectorBodyOffset)))
	}

	// TODO extend this once matrices work so it is recursive
	if prism.IsVector(val.Type) {
		env.compileInlineMap(prism.MakeMC(env.compileInlinePrintSpace, true), val)

		if end == "\x0A" {
			return env.Block.NewCall(env.getPutchar(), i32(0x0A)) // newline
		} else if end == "" {
			return env.Block.NewCall(env.getPutchar(), i32(0x20)) // space
		} else {
			return env.Block.NewCall(env.getPutchar(), i32(0)) // no
		}
	}

	return env.Block.NewCall(
		env.getPrintf(),
		env.getFormatString(val.Type, end),
		val.Value)
}

func (env *Environment) compileInlinePrintSpace(val prism.Value) value.Value {
	return env.invokePrint(val, "\x20")
}

func (env *Environment) compileInlinePrintln(val prism.Value) value.Value {
	return env.invokePrint(val, "\x0A")
}

func (env *Environment) compileInlinePrint(val prism.Value) value.Value {
	return env.invokePrint(val, "")
}

func (env *Environment) compileInlineIndex(left, right prism.Value) value.Value {
	return env.readVectorElement(right, left.Value)
}

func (env *Environment) compileInlinePanic(val prism.Value) value.Value {
	env.Block.NewCall(env.getExit(), env.Block.NewTrunc(val.Value, types.I32))
	env.Block.NewUnreachable()
	return nil
}
