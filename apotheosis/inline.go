package apotheosis

import (
	"github.com/sundown/solution/prism"

	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

func (env *Environment) compileInlinePrintln(val prism.Value) value.Value {
	if val.Type.Equals(prism.StringType) {
		return env.Block.NewCall(
			env.getPrintf(),
			env.getFormatStringln(&val.Type),
			env.Block.NewLoad(types.I8Ptr, env.Block.NewGetElementPtr(
				val.Type.Realise(),
				val.Value,
				i32(0), vectorBodyOffset)))
	}

	return env.Block.NewCall(
		env.getPrintf(),
		env.getFormatStringln(&val.Type),
		val.Value)
}

func (env *Environment) compileInlinePrint(val prism.Value) value.Value {
	if val.Type.Equals(prism.StringType) {
		return env.Block.NewCall(
			env.getPrintf(),
			env.getFormatString(&val.Type),
			env.Block.NewLoad(types.I8Ptr, env.Block.NewGetElementPtr(
				val.Type.Realise(),
				val.Value,
				i32(0), vectorBodyOffset)))
	}

	return env.Block.NewCall(
		env.getPrintf(),
		env.getFormatString(&val.Type),
		val.Value)
}

func (env *Environment) compileInlineIndex(left, right prism.Value) value.Value {
	return env.readVectorElement(right, left.Value)
}

func (env *Environment) compileInlinePanic(val prism.Value) value.Value {
	env.Block.NewCall(env.getExit(), env.Block.NewTrunc(val.Value, types.I32))
	env.Block.NewUnreachable()
	return nil
}
