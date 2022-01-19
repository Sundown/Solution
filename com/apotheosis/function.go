package apotheosis

import (
	"sundown/solution/prism"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
)

func (state *State) CompileBlock(body *prism.Expression) {
	// Block is just an expression[]
	for _, stmt := range body.Block {
		state.CompileExpression(stmt)
	}
}

func (state *State) DeclareFunction(fn *prism.Function) *ir.Func {
	args := []*ir.Param{ToParam(fn.TakesAlpha), ToParam(fn.TakesOmega)}
	if args[1] == nil {
		if args[0] != nil {
			args = []*ir.Param{ToParam(fn.TakesAlpha)}
		} else {
			args = []*ir.Param{}
		}
	} // I hope this is me being stupid, otherwise Go is broken

	state.CurrentFunction = state.Module.NewFunc(
		fn.ToLLVMName(),
		ToReturn(fn.Gives),
		args...)

	return state.CurrentFunction
}

func (state *State) CompileFunction(fn *prism.Function) *ir.Func {
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
func ToReturn(t *prism.Type) (typ types.Type) {
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
func ToParam(t *prism.Type) (typ *ir.Param) {
	if t.LLType == types.Void {
		typ = nil
	} else if t.Vector != nil || t.Tuple != nil {
		typ = ir.NewParam("@", types.NewPointer(t.LLType))
	} else {
		typ = ir.NewParam("@", t.LLType)
	}

	return typ
}
