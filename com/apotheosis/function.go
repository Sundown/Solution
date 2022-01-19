package apotheosis

import (
	"sundown/solution/prism"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
)

func (state *State) CompileBlock(body *[]prism.Expression) {
	// Block is just an expression[]
	for _, stmt := range *body {
		state.CompileExpression(&stmt)
	}
}

func (state *State) DeclareFunction(fn prism.Function) *ir.Func {
	args := []*ir.Param{}
	switch f := fn.(type) {
	case prism.DFunction:
		args = []*ir.Param{ToParam(f.AlphaType), ToParam(f.OmegaType)}
	case prism.MFunction:
		args = []*ir.Param{ToParam(f.OmegaType)}
	}

	state.CurrentFunction = state.Module.NewFunc(
		fn.LLVMise(),
		ToReturn(fn.Type()),
		args...)

	return state.CurrentFunction
}

func (state *State) CompileDFunction(fn prism.DFunction) *ir.Func {
	state.CurrentFunction = state.DFunctions[fn.LLVMise()]
	state.CurrentFunctionIR = fn

	state.Block = state.CurrentFunction.NewBlock("")
	state.CompileBlock(&fn.Body)

	if fn.Returns.Kind() == prism.VoidType.ID {
		state.Block.NewRet(nil)
	}

	return state.CurrentFunction
}

func (state *State) CompileMFunction(fn prism.MFunction) *ir.Func {
	state.CurrentFunction = state.MFunctions[fn.LLVMise()]
	state.CurrentFunctionIR = fn

	state.Block = state.CurrentFunction.NewBlock("")
	state.CompileBlock(&fn.Body)

	if fn.Returns.Kind() == prism.VoidType.ID {
		state.Block.NewRet(nil)
	}

	return state.CurrentFunction
}

// Complex types decay to pointers, atomic types do not
func ToReturn(t prism.Type) (typ types.Type) {
	if t.Kind() == prism.VoidType.ID {
		typ = types.Void
	} else if _, ok := t.(prism.AtomicType); !ok {
		typ = types.NewPointer(t.Realise())
	} else {
		typ = t.Realise()
	}

	return typ
}

// Handle void parameters and add pointers to complex types
func ToParam(t prism.Type) (typ *ir.Param) {
	if t.Kind() == prism.VoidType.ID {
		typ = nil
	} else if _, ok := t.(prism.AtomicType); !ok {
		typ = ir.NewParam("", types.NewPointer(t.Realise()))
	} else {
		typ = ir.NewParam("", t.Realise())
	}

	return typ
}
