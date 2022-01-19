package apotheosis

import (
	"sundown/solution/oversight"
	"sundown/solution/prism"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
)

type State struct {
	Runtime           *oversight.Runtime
	Env               *prism.Environment
	Module            *ir.Module
	Block             *ir.Block
	DFunctions        map[string]*ir.Func
	MFunctions        map[string]*ir.Func
	Specials          map[string]*ir.Func
	CurrentFunction   *ir.Func
	CurrentFunctionIR prism.Expression
	PanicStrings      map[string]*ir.Global
}

func (state *State) Compile(*prism.Environment) *ir.Module {
	oversight.Verbose("Init compiler")
	state.Specials = make(map[string]*ir.Func)
	state.DFunctions = make(map[string]*ir.Func)
	state.MFunctions = make(map[string]*ir.Func)
	state.PanicStrings = make(map[string]*ir.Global)

	// Root reference of IR still useful at some points
	//state.IR = IR

	state.Module = ir.NewModule()

	state.Module.SourceFilename = "out_file"

	state.
		DeclareFunctions().
		CompileFunctions().
		InitMain()

	if state.Runtime.Output == "" {
		state.Runtime.Output = "out_file"
	}

	return state.Module
}

func (state *State) DeclareFunctions() *State {
	for _, fn := range state.Env.DFunctions {
		state.DFunctions[fn.LLVMise()] = state.DeclareFunction(fn)
	}
	for _, fn := range state.Env.MFunctions {
		state.MFunctions[fn.LLVMise()] = state.DeclareFunction(fn)
	}

	return state
}

func (state *State) CompileFunctions() *State {
	for _, fn := range state.Env.DFunctions {
		state.DFunctions[fn.LLVMise()] = state.CompileDFunction(*fn)
	}
	for _, fn := range state.Env.MFunctions {
		state.MFunctions[fn.LLVMise()] = state.CompileMFunction(*fn)
	}

	return state
}

func (state *State) InitMain() *State {
	state.CurrentFunction = state.Module.NewFunc("main", types.I32)
	state.Block = state.CurrentFunction.NewBlock("entry")
	// TODO state.Block.NewCall(state.Functions[state.IR.EntryFunction.ToLLVMName()])
	state.Block.NewRet(constant.NewInt(types.I32, 0))

	return state
}
