package compiler

import (
	"sundown/sunday/parse"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
)

func (state *State) CompileFunction(fn *parse.Function) *ir.Func {
	state.CurrentFunction = state.Module.NewFunc(
		fn.ToLLVMName(),
		fn.Gives.LLType,
		ToParam(fn.Takes)...)

	state.Block = state.CurrentFunction.NewBlock("entry")
	state.CompileBlock(fn.Body)

	if fn.Gives.LLType == types.Void {
		state.Block.NewRet(nil)
	}

	return state.CurrentFunction
}

func (state *State) DeclareFunction(fn *parse.Function) *ir.Func {
	state.CurrentFunction = state.Module.NewFunc(
		fn.ToLLVMName(),
		fn.Gives.LLType,
		ToParam(fn.Takes)...)

	return state.CurrentFunction
}

func ToParam(t *parse.Type) []*ir.Param {
	if t.LLType == types.Void {
		return []*ir.Param{}
	}

	return []*ir.Param{ir.NewParam("@", t.LLType)}
}
