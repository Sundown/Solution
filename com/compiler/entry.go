package compiler

import (
	"io/ioutil"
	"sundown/solution/parse"
	"sundown/solution/util"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
)

type State struct {
	Runtime           *util.Runtime
	IR                *parse.State
	Module            *ir.Module
	Block             *ir.Block
	Functions         map[string]*ir.Func
	Specials          map[string]*ir.Func
	CurrentFunction   *ir.Func
	CurrentFunctionIR *parse.Function
}

func (state *State) Compile(IR *parse.State) {
	state.Specials = make(map[string]*ir.Func)
	state.Functions = make(map[string]*ir.Func)

	// Root reference of IR still useful at some points
	state.IR = IR

	state.Module = ir.NewModule()

	state.Module.SourceFilename = *state.IR.PackageIdent

	state.
		DeclareFunctions().
		CompileFunctions().
		InitMain()

	ioutil.WriteFile(*state.IR.PackageIdent+".ll", []byte(state.Module.String()), 0644)
}

func (state *State) DeclareFunctions() *State {
	for _, fn := range state.IR.Functions {
		if fn.Special {
			// Special form, internally defined
			continue
		}

		state.Functions[fn.ToLLVMName()] = state.DeclareFunction(fn)
	}

	return state
}

func (state *State) CompileFunctions() *State {
	for _, fn := range state.IR.Functions {
		if fn.Special {
			// Special form, internally defined
			continue
		}

		state.Functions[fn.ToLLVMName()] = state.CompileFunction(fn)
	}

	return state
}

func (state *State) InitMain() *State {
	state.CurrentFunction = state.Module.NewFunc("main", types.I32)
	state.Block = state.CurrentFunction.NewBlock("entry")
	state.Block.NewCall(state.Functions[state.IR.EntryFunction.ToLLVMName()])
	state.Block.NewRet(constant.NewInt(types.I32, 0))

	return state
}
