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

func (state *State) DeclareDFunction(fn prism.DFunction) *ir.Func {
	return state.Module.NewFunc(
		fn.LLVMise(),
		ToReturn(fn.Type()),
		ToParam(fn.AlphaType), ToParam(fn.OmegaType))
}

func (state *State) DeclareMFunction(fn prism.MFunction) *ir.Func {
	return state.Module.NewFunc(
		fn.LLVMise(),
		ToReturn(fn.Type()),
		ToParam(fn.OmegaType))
}

func (state *State) CompileDFunction(fn prism.DFunction) *ir.Func {
	state.CurrentFunction = state.DFunctions[fn.LLVMise()]
	state.CurrentFunctionIR = fn

	state.Block = state.CurrentFunction.NewBlock("")
	state.CompileBlock(&fn.Body)

	if fn.Returns.Kind() == prism.VoidType.ID {
		state.Block.NewRet(nil)
	}

	// TODO remove this
	state.Block.NewRet(state.DefaultValue(fn.Returns))

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
