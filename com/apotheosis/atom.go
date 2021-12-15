package apotheosis

import (
	"fmt"
	"sundown/solution/subtle"

	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

func (state *State) CompileAtom(morpheme *subtle.Morpheme) value.Value {
	switch {
	case morpheme.Param != nil:
		return state.CurrentFunction.Params[0]
	case morpheme.Int != nil:
		return constant.NewInt(types.I64, *morpheme.Int)
	case morpheme.Real != nil:
		return constant.NewFloat(types.Double, *morpheme.Real)
	case morpheme.Char != nil:
		return constant.NewInt(types.I8, int64(*morpheme.Char))
	case morpheme.Bool != nil:
		return constant.NewBool(*morpheme.Bool)
	case morpheme.Vector != nil:
		return state.CompileVector(morpheme)
	case morpheme.Tuple != nil:
		return state.CompileTuple(morpheme)
	case morpheme.Function != nil:
		fmt.Println("fn", morpheme)
		return state.Functions[morpheme.Function.ToLLVMName()]
	default:
		panic("unreachable")
	}
}
