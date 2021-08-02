package compiler

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

func I64(v int64) constant.Constant {
	return constant.NewInt(types.I64, v)
}

func I32(v int64) constant.Constant {
	return constant.NewInt(types.I32, int64(int32(v)))
}

// Abstract llir's really stupid get implementation
func (state *State) GEP(source *ir.InstAlloca, indices ...value.Value) *ir.InstGetElementPtr {
	return state.Block.NewGetElementPtr(source.Typ.ElemType, source, indices...)
}
