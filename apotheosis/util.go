package apotheosis

import (
	"github.com/sundown/solution/prism"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

type Value struct {
	Value value.Value
	Type  prism.Type
}

func (env Environment) New(val value.Value) (res value.Value) {
	res = env.Block.NewAlloca(val.Type())
	env.Block.NewStore(val, res)
	return
}
func I64(v int64) constant.Constant {
	return constant.NewInt(types.I64, v)
}

func F64(v float64) constant.Constant {
	return constant.NewFloat(types.Double, v)
}
func I32(v int64) constant.Constant {
	return constant.NewInt(types.I32, int64(int32(v)))
}

// Abstract LLIR's stupid GEP implementation
func (env *Environment) GEP(source *ir.InstAlloca, indices ...value.Value) *ir.InstGetElementPtr {
	return env.Block.NewGetElementPtr(source.Typ.ElemType, source, indices...)
}

// Will work for vectors too once they can be mutated
func (env *Environment) DefaultValue(t prism.Type) value.Value {
	if t.Equals(prism.IntType) {
		return I64(0)
	} else if t.Equals(prism.RealType) {
		return constant.NewFloat(types.Double, 0)
	} else if t.Equals(prism.CharType) {
		return constant.NewInt(types.I8, 0)
	} else if t.Equals(prism.BoolType) {
		return constant.NewBool(false)
	} else {
		prism.Panic("Not yet implemented")
	}
	panic(nil)
}

// Will work for vectors too once they can be mutated
func (env *Environment) Number(t *prism.Type, n float64) value.Value {
	if (*t).Equals(prism.IntType) {
		return I64(int64(n))
	} else if (*t).Equals(prism.RealType) {
		return constant.NewFloat(types.Double, n)
	} else if (*t).Equals(prism.CharType) {
		return constant.NewInt(types.I8, int64(n))
	} else if (*t).Equals(prism.BoolType) {
		return constant.NewBool(false)
	} else {
		prism.Panic("Not yet implemented")
	}
	panic(nil)
}

func (env *Environment) AgnosticAdd(t *prism.Type, x, y value.Value) value.Value {
	if (*t).Equals(prism.IntType) {
		return env.Block.NewAdd(x, y)
	} else if (*t).Equals(prism.RealType) {
		return env.Block.NewFAdd(x, y)
	} else if (*t).Equals(prism.CharType) {
		return env.Block.NewAdd(x, y)
	} else {
		prism.Panic("Not yet implemented")
	}
	panic(nil)
}

func (env *Environment) AgnosticMult(t *prism.Type, x, y value.Value) value.Value {
	if (*t).Equals(prism.IntType) {
		return env.Block.NewMul(x, y)
	} else if (*t).Equals(prism.RealType) {
		return env.Block.NewFMul(x, y)
	} else if (*t).Equals(prism.CharType) {
		return env.Block.NewMul(x, y)
	} else {
		prism.Panic("Not yet implemented")
	}
	panic(nil)
}

func (env *Environment) GetFormatStringln(t *prism.Type) value.Value {
	if (*t).Equals(prism.StringType) {
		return env.Block.NewGetElementPtr(types.NewArray(4, types.I8), env.Module.NewGlobalDef("", constant.NewCharArrayFromString("%s\x0A\x00")), I32(0), I32(0))
	} else if (*t).Equals(prism.IntType) {
		return env.Block.NewGetElementPtr(types.NewArray(4, types.I8), env.Module.NewGlobalDef("", constant.NewCharArrayFromString("%d\x0A\x00")), I32(0), I32(0))
	} else if (*t).Equals(prism.RealType) {
		return env.Block.NewGetElementPtr(types.NewArray(4, types.I8), env.Module.NewGlobalDef("", constant.NewCharArrayFromString("%f\x0A\x00")), I32(0), I32(0))
	} else if (*t).Equals(prism.CharType) {
		return env.Block.NewGetElementPtr(types.NewArray(4, types.I8), env.Module.NewGlobalDef("", constant.NewCharArrayFromString("%c\x0A\x00")), I32(0), I32(0))
	} else if (*t).Equals(prism.BoolType) {
		return env.Block.NewGetElementPtr(types.NewArray(4, types.I8), env.Module.NewGlobalDef("", constant.NewCharArrayFromString("%d\x0A\x00")), I32(0), I32(0))
	} else {
		return env.Block.NewGetElementPtr(types.NewArray(2, types.I8), env.Module.NewGlobalDef("", constant.NewCharArrayFromString("\x0A\x00")), I32(0), I32(0))
	}
}

func (env *Environment) GetFormatString(t *prism.Type) value.Value {
	if (*t).Equals(prism.StringType) {
		return env.Block.NewGetElementPtr(types.NewArray(3, types.I8), env.Module.NewGlobalDef("", constant.NewCharArrayFromString("%s\x00")), I32(0), I32(0))
	} else if (*t).Equals(prism.IntType) {
		return env.Block.NewGetElementPtr(types.NewArray(3, types.I8), env.Module.NewGlobalDef("", constant.NewCharArrayFromString("%d\x00")), I32(0), I32(0))
	} else if (*t).Equals(prism.RealType) {
		return env.Block.NewGetElementPtr(types.NewArray(3, types.I8), env.Module.NewGlobalDef("", constant.NewCharArrayFromString("%f\x00")), I32(0), I32(0))
	} else if (*t).Equals(prism.CharType) {
		return env.Block.NewGetElementPtr(types.NewArray(3, types.I8), env.Module.NewGlobalDef("", constant.NewCharArrayFromString("%c\x00")), I32(0), I32(0))
	} else if (*t).Equals(prism.BoolType) {
		return env.Block.NewGetElementPtr(types.NewArray(3, types.I8), env.Module.NewGlobalDef("", constant.NewCharArrayFromString("%d\x00")), I32(0), I32(0))
	} else {
		return env.Block.NewGetElementPtr(types.NewArray(1, types.I8), env.Module.NewGlobalDef("", constant.NewCharArrayFromString("\x00")), I32(0), I32(0))
	}
}

// Supply the block in which to generate message and exit call, a printf formatter, and variadic params
func (env *Environment) LLVMPanic(block *ir.Block, format string, args ...value.Value) {
	// Certain panic strings are very common, such as bounds checks, this ensured they are not double-allocated.
	fmt_glob := env.PanicStrings[format]
	if fmt_glob == nil {
		fmt_glob = env.Module.NewGlobalDef("", constant.NewCharArrayFromString(format+"\x00"))
		env.PanicStrings[format] = fmt_glob
	}

	block.NewCall(env.GetPrintf(), append([]value.Value{block.NewGetElementPtr(
		types.NewArray(uint64(len(format)+1), types.I8), fmt_glob, I32(0), I32(0))}, args...)...)
	block.NewCall(env.GetExit(), I32(1))
}
