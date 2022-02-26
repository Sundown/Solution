package apotheosis

import (
	"github.com/sundown/solution/prism"

	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

func (env *Environment) compileInlinePrintln(val Value) value.Value {
	if val.Type.Equals(prism.StringType) {
		return env.Block.NewCall(
			env.GetPrintf(),
			env.GetFormatStringln(&val.Type),
			env.Block.NewLoad(types.I8Ptr, env.Block.NewGetElementPtr(
				val.Type.Realise(),
				val.Value,
				I32(0), vectorBodyOffset)))
	}

	return env.Block.NewCall(
		env.GetPrintf(),
		env.GetFormatStringln(&val.Type),
		val.Value)
}

func (env *Environment) compileInlinePrint(val Value) value.Value {
	if val.Type.Equals(prism.StringType) {
		return env.Block.NewCall(
			env.GetPrintf(),
			env.GetFormatString(&val.Type),
			env.Block.NewLoad(types.I8Ptr, env.Block.NewGetElementPtr(
				val.Type.Realise(),
				val.Value,
				I32(0), vectorBodyOffset)))
	}

	return env.Block.NewCall(
		env.GetPrintf(),
		env.GetFormatString(&val.Type),
		val.Value)
}

func (env *Environment) compileInlineIndex(left, right Value) value.Value {
	return env.readVectorElement(right, left.Value)
}

func (env *Environment) compileInlinePanic(val Value) value.Value {
	env.Block.NewCall(env.GetExit(), env.Block.NewTrunc(val.Value, types.I32))
	env.Block.NewUnreachable()
	return nil
}
