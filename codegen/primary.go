package codegen

import (
	"fmt"
	"sundown/sunday/parser"

	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

func (state *State) MakePrimary(primary *parser.Primary) (value.Value, types.Type) {
	if primary.Param != nil {
		// Funtions only have one parameter
		return state.function.Params[0], state.function.Params[0].Type()
	} else if primary.Int != nil {
		return constant.NewInt(&int_bt, *primary.Int), int_t
	} else if primary.Real != nil {
		return constant.NewFloat(types.Double, *primary.Real), real_t
	} else if primary.Bool != nil {
		// All builtin keywords are Titlecase
		if *primary.Bool == "True" {
			return constant.NewBool(true), bool_t
		} else if *primary.Bool == "False" {
			return constant.NewBool(false), bool_t
		}
	} else if primary.Vec != nil {
		// Yet to be fully implemented due to complexity of some types
		// I.E. vector of tuple of (T_a, T_b)
		return state.MakeVector(primary.Vec)
	} else if primary.String != nil {
		// Trim the " from head and tail of string (left there by parser)
		// and append null terminator
		s := (*primary.String)[1:len(*primary.String)-1] + "\x00"

		// Return GEP of global def'd string,
		// could change this to a GEP of an alloca
		return constant.NewGetElementPtr(
			types.NewArray(uint64(len(s)), types.I8),
			state.module.NewGlobalDef("", constant.NewCharArrayFromString(s)),
			constant.NewInt(types.I32, 0),
			constant.NewInt(types.I32, 0)), types.I8Ptr
	} else {
		fmt.Println("other")
	}

	return nil, nil
}
