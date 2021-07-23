package compiler

import (
	"sundown/sunday/parse"

	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

func (state *State) CompileAtom(atom *parse.Atom) value.Value {
	if atom.Param != nil {
		return state.CurrentFunction.Params[0]
	} else if atom.Int != nil {
		return constant.NewInt(types.I64, *atom.Int)
	} else if atom.Real != nil {
		return constant.NewFloat(types.Double, *atom.Real)
	} else if atom.Char != nil {
		return constant.NewInt(types.I8, int64(*atom.Char))
	} else if atom.Bool != nil {
		return constant.NewBool(*atom.Bool)
	} else if atom.Vector != nil {
		return state.CompileVector(atom)
	} else if atom.Tuple != nil {
		return state.CompileTuple(atom)
	} else {
		panic("unreachable")
	}
}
