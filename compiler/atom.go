package compiler

import (
	"sundown/sunday/parse"

	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

func (state *State) CompileAtom(atom *parse.Atom) value.Value {
	switch {
	case atom.Param != nil:
		return state.CurrentFunction.Params[0]
	case atom.Int != nil:
		return constant.NewInt(types.I64, *atom.Int)
	case atom.Real != nil:
		return constant.NewFloat(types.Double, *atom.Real)
	case atom.Char != nil:
		return constant.NewInt(types.I8, int64(*atom.Char))
	case atom.Bool != nil:
		return constant.NewBool(*atom.Bool)
	case atom.Vector != nil:
		return state.CompileVector(atom)
	case atom.Tuple != nil:
		return state.CompileTuple(atom)
	default:
		panic("unreachable")
	}
}
