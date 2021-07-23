package compiler

import (
	"sundown/sunday/parse"

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

func (state *State) CompileApplication(app *parse.Application) value.Value {
	switch *app.Function.Ident.Ident {
	case "Return":
		state.Block.NewRet(state.CompileExpression(app.Argument))
		return nil
	case "GEP":
		a := app.Argument.Atom

		if a == nil || a.Tuple == nil {
			panic("GEP requires tuple: (<structure>, uint)")
		}

		tuple := state.CompileExpression(app.Argument.Atom.Tuple[0])
		typ := app.Argument.Atom.Tuple[0].TypeOf

		ll_typ := typ.AsLLType()

		i := app.Argument.Atom.Tuple[1].Atom.Int

		gep := state.Block.NewGetElementPtr(ll_typ, tuple, I32(0), I32(*i))

		return gep
	default:
		return state.Block.NewCall(
			state.Functions[app.Function.ToLLVMName()],
			state.CompileExpression(app.Argument))
	}
}
