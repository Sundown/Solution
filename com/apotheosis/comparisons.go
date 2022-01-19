package apotheosis

import (
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/value"
)

func (env *Environment) CompileInlineEqual(left Value, right Value) value.Value {
	return env.Block.NewICmp(enum.IPredEQ, left.Value, right.Value)
}
