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

func (state *State) GetCalloc() *ir.Func {
	if state.Specials["calloc"] == nil {
		state.Specials["calloc"] = state.Module.NewFunc(
			"calloc",
			types.I8Ptr,
			ir.NewParam("size", types.I32),
			ir.NewParam("count", types.I32))
	}

	return state.Specials["calloc"]
}

func (state *State) GetExit() *ir.Func {
	if state.Specials["exit"] == nil {
		state.Specials["exit"] = state.Module.NewFunc(
			"exit",
			types.Void,
			ir.NewParam("", types.I32))
	}

	return state.Specials["exit"]
}

func (state *State) GetPrintf() *ir.Func {
	if state.Specials["printf"] == nil {
		state.Specials["printf"] = state.Module.NewFunc(
			"printf",
			types.I32,
			ir.NewParam("", types.I8Ptr))
		state.Specials["printf"].Sig.Variadic = true
	}

	return state.Specials["printf"]
}
