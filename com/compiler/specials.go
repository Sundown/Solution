package compiler

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
)

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
