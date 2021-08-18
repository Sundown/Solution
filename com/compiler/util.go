package compiler

import (
	"sundown/solution/parse"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

func I64(v int64) constant.Constant {
	return constant.NewInt(types.I64, v)
}

func I32(v int64) constant.Constant {
	return constant.NewInt(types.I32, int64(int32(v)))
}

// Abstract LLIR's stupid GEP implementation
func (state *State) GEP(source *ir.InstAlloca, indices ...value.Value) *ir.InstGetElementPtr {
	return state.Block.NewGetElementPtr(source.Typ.ElemType, source, indices...)
}

// Will work for vectors too once they can be mutated
func (state *State) DefaultValue(t *parse.Type) value.Value {
	if t.Equals(&parse.IntType) {
		return I64(0)
	} else if t.Equals(&parse.RealType) {
		return constant.NewFloat(types.Double, 0)
	} else if t.Equals(&parse.CharType) {
		return constant.NewInt(types.I8, 0)
	} else if t.Equals(&parse.BoolType) {
		return constant.NewBool(false)
	} else {
		panic("Not yet implemented")
	}
}

func (state *State) AgnosticAdd(t *parse.Type, x, y value.Value) value.Value {
	if t.Equals(&parse.IntType) {
		return state.Block.NewAdd(x, y)
	} else if t.Equals(&parse.RealType) {
		return state.Block.NewFAdd(x, y)
	} else if t.Equals(&parse.CharType) {
		return state.Block.NewAdd(x, y)
	} else {
		panic("Not yet implemented")
	}
}

func (state *State) GetFormatString(t *parse.Type) value.Value {
	if t.Equals(&parse.StringType) {
		return state.Block.NewGetElementPtr(types.NewArray(4, types.I8), state.Module.NewGlobalDef("", constant.NewCharArrayFromString("%s\x0A\x00")), I32(0), I32(0))
	} else if t.Equals(&parse.IntType) {
		return state.Block.NewGetElementPtr(types.NewArray(4, types.I8), state.Module.NewGlobalDef("", constant.NewCharArrayFromString("%d\x0A\x00")), I32(0), I32(0))
	} else if t.Equals(&parse.RealType) {
		return state.Block.NewGetElementPtr(types.NewArray(4, types.I8), state.Module.NewGlobalDef("", constant.NewCharArrayFromString("%f\x0A\x00")), I32(0), I32(0))
	} else if t.Equals(&parse.CharType) {
		return state.Block.NewGetElementPtr(types.NewArray(4, types.I8), state.Module.NewGlobalDef("", constant.NewCharArrayFromString("%c\x0A\x00")), I32(0), I32(0))
	} else {
		return state.Block.NewGetElementPtr(types.NewArray(2, types.I8), state.Module.NewGlobalDef("", constant.NewCharArrayFromString("\x0A\x00")), I32(0), I32(0))
	}
}

// Supply the block in which to generate message and exit call, a printf formatter, and variadic params
func (state *State) LLVMPanic(block *ir.Block, format string, args ...value.Value) {
	var fmt value.Value = block.NewGetElementPtr(
		types.NewArray(uint64(len(format)+1), types.I8),
		state.Module.NewGlobalDef("", constant.NewCharArrayFromString(format+"\x00")), I32(0), I32(0))
	block.NewCall(state.GetPrintf(), append([]value.Value{fmt}, args...)...)
	block.NewCall(state.GetExit(), I32(1))
}
