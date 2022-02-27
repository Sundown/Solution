package apotheosis

import (
	"github.com/sundown/solution/prism"

	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

func (env *Environment) compileAtom(morpheme *prism.Morpheme) value.Value {
	switch v := (*morpheme).(type) {
	case prism.Int:
		return constant.NewInt(types.I64, v.Value)
	case prism.Real:
		return constant.NewFloat(types.Double, v.Value)
	case prism.Char:
		return constant.NewInt(types.I8, int64(v.Value[0]))
	case prism.Bool:
		return constant.NewBool(v.Value)
	case prism.Vector:
		return env.compileVector(v)
	/*case prism.Tuple:
		return env.compileTuple(morpheme)
	case prism.Function:
		fmt.Println("fn", morpheme)
		return env.Functions[morpheme.Function.ToLLVMName()] */
	default:
		prism.Panic("unreachable")
	}
	panic(nil)
}
