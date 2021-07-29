package compiler

import (
	"sundown/sunday/parse"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
)

func (state *State) CompileFunction(fn *parse.Function) *ir.Func {
	state.CurrentFunction = state.Functions[fn.ToLLVMName()]
	state.CurrentFunctionIR = fn

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
		ToReturn(fn.Gives),
		ToParam(fn.Takes)...)

	return state.CurrentFunction
}

func ToReturn(t *parse.Type) (typ types.Type) {
	if t.LLType == types.Void {
		typ = types.Void
	} else if t.Vector != nil || t.Tuple != nil {
		typ = types.NewPointer(t.LLType)
	} else {
		typ = t.LLType
	}

	return typ
}

func ToParam(t *parse.Type) (typ []*ir.Param) {
	if t.LLType == types.Void {
		typ = []*ir.Param{}
	} else if t.Vector != nil || t.Tuple != nil {
		typ = []*ir.Param{ir.NewParam("@", types.NewPointer(t.LLType))}
	} else {
		typ = []*ir.Param{ir.NewParam("@", t.LLType)}
	}

	return typ
}