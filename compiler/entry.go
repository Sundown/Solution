package compiler

import (
	"io/ioutil"
	"sundown/sunday/parse"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
)

type State struct {
	IR              *parse.State
	Module          *ir.Module
	Block           *ir.Block
	Functions       map[string]*ir.Func
	Specials        map[string]*ir.Func
	CurrentFunction *ir.Func
}

func (state *State) Compile(IR *parse.State) {
	state.Specials = make(map[string]*ir.Func)
	state.Functions = make(map[string]*ir.Func)

	state.IR = IR

	state.Module = ir.NewModule()
	state.Module.SourceFilename = *state.IR.PackageIdent

	state.InitCalloc()

	// This doesn't seem to be necessary because parser already substitutes
	/* for key, def := range state.IR.NounDefs {
		state.Module.NewGlobalDef(key.Ident, state.CompileAtom(def))
	} */

	for _, fn := range state.IR.Functions {
		if *fn.Ident.Ident == "Return" {
			continue
		}

		state.Functions[fn.ToLLVMName()] = state.DeclareFunction(fn)
	}

	for _, fn := range state.IR.Functions {
		if *fn.Ident.Ident == "Return" {
			continue
		}

		state.Functions[fn.ToLLVMName()] = state.CompileFunction(fn)
	}

	state.CurrentFunction = state.Module.NewFunc("main", types.I32)
	state.Block = state.CurrentFunction.NewBlock("entry")
	state.Block.NewCall(state.Functions[state.IR.EntryFunction.ToLLVMName()])
	state.Block.NewRet(constant.NewInt(types.I32, 0))

	ioutil.WriteFile("out.ll", []byte(state.Module.String()), 0644)
}
