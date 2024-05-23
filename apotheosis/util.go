package apotheosis

import (
	"github.com/sundown/solution/prism"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

func (env *Environment) new(val value.Value) (res value.Value) {
	res = env.Block.NewAlloca(val.Type())
	env.Block.NewStore(val, res)
	return
}

func i64(v int64) constant.Constant {
	return constant.NewInt(types.I64, v)
}

func f64(v float64) constant.Constant {
	return constant.NewFloat(types.Double, v)
}
func i32(v int64) constant.Constant {
	return constant.NewInt(types.I32, int64(int32(v)))
}

func (env *Environment) getFormatString(t prism.Type, end string) value.Value {
	format, size := "", 3
	switch t.Kind() {
	case prism.StringType.Kind():
		format = "%s"
	case prism.IntType.Kind():
		format = "%d"
	case prism.RealType.Kind():
		format = "%f"
	case prism.CharType.Kind():
		format = "%c"
	case prism.BoolType.Kind():
		format = "%d"
	default:
		size = 2
		format = ""
	}

	return env.Block.NewGetElementPtr(
		types.NewArray(uint64(size+len(end)), types.I8),
		env.Module.NewGlobalDef("",
			constant.NewCharArrayFromString(format+end+"\x00")),
		i32(0), i32(0))
}

// Supply the block in which to generate message and exit call, a printf formatter, and variadic params
func (env *Environment) newPanic(block *ir.Block, format string, args ...value.Value) {
	// Certain panic strings are very common, such as bounds checks, this ensured they are not double-allocated.
	fmtGlob := env.PanicStrings[format]
	if fmtGlob == nil {
		fmtGlob = env.Module.NewGlobalDef("", constant.NewCharArrayFromString(format+"\x00"))
		env.PanicStrings[format] = fmtGlob
	}

	block.NewCall(env.getPrintf(), append([]value.Value{block.NewGetElementPtr(
		types.NewArray(uint64(len(format)+1), types.I8), fmtGlob, i32(0), i32(0))}, args...)...)
	block.NewCall(env.getExit(), i32(1))
}
