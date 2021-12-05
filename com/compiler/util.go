package compiler

import (
	"sundown/solution/temporal"

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
func (state *State) DefaultValue(t *temporal.Type) value.Value {
	if t.Equals(&temporal.IntType) {
		return I64(0)
	} else if t.Equals(&temporal.RealType) {
		return constant.NewFloat(types.Double, 0)
	} else if t.Equals(&temporal.CharType) {
		return constant.NewInt(types.I8, 0)
	} else if t.Equals(&temporal.BoolType) {
		return constant.NewBool(false)
	} else {
		panic("Not yet implemented")
	}
}

// Will work for vectors too once they can be mutated
func (state *State) Number(t *temporal.Type, n float64) value.Value {
	if t.Equals(&temporal.IntType) {
		return I64(int64(n))
	} else if t.Equals(&temporal.RealType) {
		return constant.NewFloat(types.Double, n)
	} else if t.Equals(&temporal.CharType) {
		return constant.NewInt(types.I8, int64(n))
	} else if t.Equals(&temporal.BoolType) {
		return constant.NewBool(false)
	} else {
		panic("Not yet implemented")
	}
}

func (state *State) AgnosticAdd(t *temporal.Type, x, y value.Value) value.Value {
	if t.Equals(&temporal.IntType) {
		return state.Block.NewAdd(x, y)
	} else if t.Equals(&temporal.RealType) {
		return state.Block.NewFAdd(x, y)
	} else if t.Equals(&temporal.CharType) {
		return state.Block.NewAdd(x, y)
	} else {
		panic("Not yet implemented")
	}
}

func (state *State) AgnosticMult(t *temporal.Type, x, y value.Value) value.Value {
	if t.Equals(&temporal.IntType) {
		return state.Block.NewMul(x, y)
	} else if t.Equals(&temporal.RealType) {
		return state.Block.NewFMul(x, y)
	} else if t.Equals(&temporal.CharType) {
		return state.Block.NewMul(x, y)
	} else {
		panic("Not yet implemented")
	}
}

func (state *State) GetFormatStringln(t *temporal.Type) value.Value {
	if t.Equals(&temporal.StringType) {
		return state.Block.NewGetElementPtr(types.NewArray(4, types.I8), state.Module.NewGlobalDef("", constant.NewCharArrayFromString("%s\x0A\x00")), I32(0), I32(0))
	} else if t.Equals(&temporal.IntType) {
		return state.Block.NewGetElementPtr(types.NewArray(4, types.I8), state.Module.NewGlobalDef("", constant.NewCharArrayFromString("%d\x0A\x00")), I32(0), I32(0))
	} else if t.Equals(&temporal.RealType) {
		return state.Block.NewGetElementPtr(types.NewArray(4, types.I8), state.Module.NewGlobalDef("", constant.NewCharArrayFromString("%f\x0A\x00")), I32(0), I32(0))
	} else if t.Equals(&temporal.CharType) {
		return state.Block.NewGetElementPtr(types.NewArray(4, types.I8), state.Module.NewGlobalDef("", constant.NewCharArrayFromString("%c\x0A\x00")), I32(0), I32(0))
	} else if t.Equals(&temporal.BoolType) {
		return state.Block.NewGetElementPtr(types.NewArray(4, types.I8), state.Module.NewGlobalDef("", constant.NewCharArrayFromString("%d\x0A\x00")), I32(0), I32(0))
	} else {
		return state.Block.NewGetElementPtr(types.NewArray(2, types.I8), state.Module.NewGlobalDef("", constant.NewCharArrayFromString("\x0A\x00")), I32(0), I32(0))
	}
}

func (state *State) GetFormatString(t *temporal.Type) value.Value {
	if t.Equals(&temporal.StringType) {
		return state.Block.NewGetElementPtr(types.NewArray(3, types.I8), state.Module.NewGlobalDef("", constant.NewCharArrayFromString("%s\x00")), I32(0), I32(0))
	} else if t.Equals(&temporal.IntType) {
		return state.Block.NewGetElementPtr(types.NewArray(3, types.I8), state.Module.NewGlobalDef("", constant.NewCharArrayFromString("%d\x00")), I32(0), I32(0))
	} else if t.Equals(&temporal.RealType) {
		return state.Block.NewGetElementPtr(types.NewArray(3, types.I8), state.Module.NewGlobalDef("", constant.NewCharArrayFromString("%f\x00")), I32(0), I32(0))
	} else if t.Equals(&temporal.CharType) {
		return state.Block.NewGetElementPtr(types.NewArray(3, types.I8), state.Module.NewGlobalDef("", constant.NewCharArrayFromString("%c\x00")), I32(0), I32(0))
	} else if t.Equals(&temporal.BoolType) {
		return state.Block.NewGetElementPtr(types.NewArray(3, types.I8), state.Module.NewGlobalDef("", constant.NewCharArrayFromString("%d\x00")), I32(0), I32(0))
	} else {
		return state.Block.NewGetElementPtr(types.NewArray(1, types.I8), state.Module.NewGlobalDef("", constant.NewCharArrayFromString("\x00")), I32(0), I32(0))
	}
}

// Supply the block in which to generate message and exit call, a printf formatter, and variadic params
func (state *State) LLVMPanic(block *ir.Block, format string, args ...value.Value) {
	// Certain panic strings are very common, such as bounds checks, this ensured they are not double-allocated.
	fmt_glob := state.PanicStrings[format]
	if fmt_glob == nil {
		fmt_glob = state.Module.NewGlobalDef("", constant.NewCharArrayFromString(format+"\x00"))
		state.PanicStrings[format] = fmt_glob
	}

	block.NewCall(state.GetPrintf(), append([]value.Value{block.NewGetElementPtr(
		types.NewArray(uint64(len(format)+1), types.I8), fmt_glob, I32(0), I32(0))}, args...)...)
	block.NewCall(state.GetExit(), I32(1))
}
