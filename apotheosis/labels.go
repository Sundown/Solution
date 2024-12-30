package apotheosis

import (
	"fmt"

	"github.com/llir/llvm/ir"
)

func (env *Environment) newID() int {
	env.ApotheosisIter++
	return env.ApotheosisIter
}

func (env *Environment) newBlock(fn *ir.Func) *ir.Block {
	return fn.NewBlock(fmt.Sprint(env.newID()))
}

func (env *Environment) newNamedBlock(fn *ir.Func, name string) *ir.Block {
	return fn.NewBlock(name)
}
