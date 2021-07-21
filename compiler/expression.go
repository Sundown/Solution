package compiler

import (
	"sundown/sunday/parse"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

func (state *State) CompileExpression(expr *parse.Expression) value.Value {
	if expr.Atom != nil {
		return state.CompileAtom(expr.Atom)
	} else if expr.Application != nil {
		return state.CompileApplication(expr.Application)
	} else {
		panic("unreachable")
	}
}

func (state *State) CompileAtom(atom *parse.Atom) value.Value {
	if atom.Int != nil {
		return constant.NewInt(types.I64, *atom.Int)
	} else if atom.Real != nil {
		return constant.NewFloat(types.Double, *atom.Real)
	} else if atom.Vector != nil {
		return state.CompileVector(atom)
	} else {
		panic("unreachable")
	}
}

func (state *State) CompileApplication(app *parse.Application) value.Value {
	switch *app.Function.Ident.Ident {
	case "Return":
		state.Block.NewRet(state.CompileExpression(app.Argument))
		return nil
	default:
		return state.Block.NewCall(
			state.Functions[app.Function.ToLLVMName()],
			state.CompileExpression(app.Argument))
	}
}

func (state *State) Calloc() *ir.Func {
	return state.Module.NewFunc("calloc", types.I8Ptr,
		ir.NewParam("size", types.I32), ir.NewParam("count", types.I32))
}
