package apotheosis

import (
	"fmt"

	"github.com/llir/llvm/ir"
)

func (env *Environment) NewID() int {
	env.ApotheosisIter++
	return env.ApotheosisIter
}

func (env *Environment) NewBlock(fn *ir.Func) *ir.Block {
	return fn.NewBlock(fmt.Sprint(env.NewID()))
}
