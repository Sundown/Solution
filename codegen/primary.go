package codegen

import (
	"fmt"
	"sundown/girl/parser"

	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

func (state *State) MakePrimary(primary *parser.Primary) value.Value {
	if primary.Param != nil {
		return state.function.Params[0]
	} else if primary.Int != nil {
		return constant.NewInt(&int_bt, *primary.Int)
	} else if primary.Real != nil {
		return constant.NewFloat(types.Double, *primary.Real)
	} else if primary.Bool != nil {
		if *primary.Bool == "True" {
			return constant.NewBool(true)
		} else if *primary.Bool == "False" {
			return constant.NewBool(false)
		}
	} else if primary.Vec != nil {
		vec, _ := state.compile_vector(primary.Vec)
		return vec
	} else if primary.String != nil {
		s := (*primary.String)[1:len(*primary.String)-1] + "\x00"
		slen := int64(len(s))
		i := constant.NewCharArrayFromString(s)
		str := state.module.NewGlobalDef("", i)
		ptr := constant.NewGetElementPtr(
			types.NewArray(uint64(slen), types.I8),
			str,
			constant.NewInt(types.I32, 0),
			constant.NewInt(types.I32, 0))

		return ptr
	} else {
		fmt.Println("other")
	}

	return nil
}
