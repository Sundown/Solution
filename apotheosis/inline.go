package apotheosis

import (
	"github.com/sundown/solution/prism"

	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

func (env *Environment) compileInlineIota(val prism.Value) value.Value {
	head := env.vectorFactory(val.Type, env.Block.NewTrunc(env.castInt(val), types.I32))

	counterStore := env.new(i64(0))

	loopblock := env.newBlock(env.CurrentFunction)
	env.Block.NewBr(loopblock)
	env.Block = loopblock

	curCounter := loopblock.NewLoad(types.I64, counterStore)

	incr := loopblock.NewAdd(curCounter, i64(1))

	env.writeElement(head, incr, curCounter)

	loopblock.NewStore(incr, counterStore)

	env.Block = env.newBlock(env.CurrentFunction)
	loopblock.NewCondBr(loopblock.NewICmp(enum.IPredNE, incr, val.Value), loopblock, env.Block)

	return head.Value
}

func (env *Environment) compileInlineEnclose(val prism.Value) value.Value {
	head := env.vectorFactory(val.Type, i32(1))
	env.writeElement(head, val.Value, i32(0))
	return head.Value
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

func (env *Environment) compileInlineTally(val prism.Value) value.Value {
	return env.Block.NewSExt(env.readVectorLength(val), types.I64)
}

func (env *Environment) compileInlineCapacity(val prism.Value) value.Value {
	return env.Block.NewSExt(env.readVectorCapacity(val), types.I64)
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
	return env.readVectorElement(right, env.Block.NewTrunc(env.Block.NewSub(left.Value, i64(1)), types.I32))
}

func (env *Environment) compileInlinePanic(val prism.Value) value.Value {
	env.Block.NewCall(env.getExit(), env.Block.NewTrunc(val.Value, types.I32))
	env.Block.NewUnreachable()
	return nil
}
