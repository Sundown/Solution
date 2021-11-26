package compiler

import (
	"sundown/solution/temporal"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
)

func (state *State) CompileBlock(body *temporal.Expression) {
	// Block is just an expression[]
	for _, stmt := range body.Block {
		state.CompileExpression(stmt)
	}
}

func (state *State) DeclareFunction(fn *temporal.Function) *ir.Func {
	state.CurrentFunction = state.Module.NewFunc(
		fn.ToLLVMName(),
		ToReturn(fn.Gives),
		ToParam(fn.Takes)...)

	return state.CurrentFunction
}

func (state *State) CompileFunction(fn *temporal.Function) *ir.Func {
	state.CurrentFunction = state.Functions[fn.ToLLVMName()]
	state.CurrentFunctionIR = fn

	state.Block = state.CurrentFunction.NewBlock("")
	state.CompileBlock(fn.Body)

	if fn.Gives.LLType == types.Void {
		state.Block.NewRet(nil)
	}

	return state.CurrentFunction
}

// Complex types decay to pointers, atomic types do not
func ToReturn(t *temporal.Type) (typ types.Type) {
	if t.LLType == types.Void {
		typ = types.Void
	} else if t.Vector != nil || t.Tuple != nil {
		typ = types.NewPointer(t.LLType)
	} else {
		typ = t.LLType
	}

	return typ
}

// Handle void parameters and add pointers to complex types
func ToParam(t *temporal.Type) (typ []*ir.Param) {
	if t.LLType == types.Void {
		typ = []*ir.Param{}
	} else if t.Vector != nil || t.Tuple != nil {
		typ = []*ir.Param{ir.NewParam("@", types.NewPointer(t.LLType))}
	} else {
		typ = []*ir.Param{ir.NewParam("@", t.LLType)}
	}

	return typ
}
