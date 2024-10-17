package apotheosis

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
)

var voidVector = types.NewStruct(types.I32, types.I32, types.I32, types.I64Ptr)

func (env *Environment) getCalloc() *ir.Func {
	if env.Specials["calloc"] == nil {
		env.Specials["calloc"] = env.Module.NewFunc(
			"calloc",
			types.I8Ptr,
			ir.NewParam("size", types.I32),
			ir.NewParam("count", types.I32))
	}

	return env.Specials["calloc"]
}

// TODO APO needs to be i64
func (env *Environment) getPowInt() *ir.Func {
	if env.Specials["powi"] == nil {
		env.Specials["powi"] = env.Module.NewFunc(
			"llvm.powi.f64.i32",
			types.Double,
			ir.NewParam("b", types.Double),
			ir.NewParam("e", types.I32))
	}

	return env.Specials["powi"]
}

func (env *Environment) getPowReal() *ir.Func {
	if env.Specials["powf"] == nil {
		env.Specials["powf"] = env.Module.NewFunc(
			"llvm.powi.f64",
			types.Double,
			ir.NewParam("b", types.Double),
			ir.NewParam("e", types.Double))
	}

	return env.Specials["powf"]
}

func (env *Environment) getPutchar() *ir.Func {
	if env.Specials["putchar"] == nil {
		env.Specials["putchar"] = env.Module.NewFunc(
			"putchar",
			types.I32,
			ir.NewParam("c", types.I32))
	}

	return env.Specials["putchar"]
}

// memcpy(i8* dest, i8* src, i64 len, i1 is_volatile)
func (env *Environment) getMemcpy() *ir.Func {
	if env.Specials["memcpy"] == nil {
		env.Specials["memcpy"] = env.Module.NewFunc(
			"llvm.memcpy",
			types.Void,
			ir.NewParam("dest", types.I8Ptr),
			ir.NewParam("src", types.I8Ptr),
			ir.NewParam("len", types.I64),
			ir.NewParam("isvolatile", types.I1))
	}

	return env.Specials["memcpy"]
}

// realloc(i8* dest, i64 len)
func (env *Environment) getRealloc() *ir.Func {
	if env.Specials["realloc"] == nil {
		env.Specials["realloc"] = env.Module.NewFunc(
			"realloc",
			types.Void,
			ir.NewParam("mem", types.I8Ptr),
			ir.NewParam("len", types.I64))
	}

	return env.Specials["realloc"]
}

func (env *Environment) getExit() *ir.Func {
	if env.Specials["exit"] == nil {
		env.Specials["exit"] = env.Module.NewFunc(
			"exit",
			types.Void,
			ir.NewParam("exitcode", types.I32))
	}

	return env.Specials["exit"]
}

func (env *Environment) getMaxDouble() *ir.Func {
	if env.Specials["max.f64"] == nil {
		env.Specials["max.f64"] = env.Module.NewFunc(
			"llvm.maxnum.double",
			types.Double,
			ir.NewParam("", types.Double),
			ir.NewParam("", types.Double))

	}

	return env.Specials["max.f64"]
}

func (env *Environment) getMinDouble() *ir.Func {
	if env.Specials["min.f64"] == nil {
		env.Specials["min.f64"] = env.Module.NewFunc(
			"llvm.minnum.double",
			types.Double,
			ir.NewParam("", types.Double),
			ir.NewParam("", types.Double))

	}

	return env.Specials["min.f64"]
}

func (env *Environment) getStrlen() *ir.Func {
	if env.Specials["strlen"] == nil {
		env.Specials["strlen"] = env.Module.NewFunc(
			"strlen",
			types.I64,
			ir.NewParam("s", types.I8Ptr))
		env.Specials["strlen"].Sig.Variadic = true
	}

	return env.Specials["strlen"]
}

func (env *Environment) getPrintf() *ir.Func {
	if env.Specials["printf"] == nil {
		env.Specials["printf"] = env.Module.NewFunc(
			"printf",
			types.I32,
			ir.NewParam("fmt", types.I8Ptr))
		env.Specials["printf"].Sig.Variadic = true
	}

	return env.Specials["printf"]
}

func (env *Environment) getWriteVectorLength() *ir.Func {
	if env.Specials["writeVectorLength"] == nil {
		env.Specials["writeVectorLength"] = env.Module.NewFunc(
			"writeVectorLength",
			types.Void,
			ir.NewParam("v", types.NewPointer(voidVector)),
			ir.NewParam("length", types.I32))

	}

	return env.Specials["writeVectorLength"]
}

func (env *Environment) getWriteVectorCapacity() *ir.Func {
	if env.Specials["writeVectorCapacity"] == nil {
		env.Specials["writeVectorCapacity"] = env.Module.NewFunc(
			"writeVectorCapacity",
			types.Void,
			ir.NewParam("v", types.NewPointer(voidVector)),
			ir.NewParam("length", types.I32))

	}

	return env.Specials["writeVectorCapacity"]
}

func (env *Environment) getWriteVectorWidth() *ir.Func {
	if env.Specials["writeVectorWidth"] == nil {
		env.Specials["writeVectorWidth"] = env.Module.NewFunc(
			"writeVectorWidth",
			types.Void,
			ir.NewParam("v", types.NewPointer(voidVector)),
			ir.NewParam("length", types.I32))

	}

	return env.Specials["writeVectorWidth"]
}

func (env *Environment) getWriteVectorPointer() *ir.Func {
	if env.Specials["writeVectorPointer"] == nil {
		env.Specials["writeVectorPointer"] = env.Module.NewFunc(
			"writeVectorPointer",
			types.Void,
			ir.NewParam("v", types.NewPointer(voidVector)),
			ir.NewParam("data", types.I8Ptr))

	}

	return env.Specials["writeVectorPointer"]
}

func (env *Environment) getReadVectorLength() *ir.Func {
	if env.Specials["readVectorLength"] == nil {
		env.Specials["readVectorLength"] = env.Module.NewFunc(
			"readVectorLength",
			types.I32,
			ir.NewParam("v", types.NewPointer(voidVector)))

	}

	return env.Specials["readVectorLength"]
}

func (env *Environment) getReadVectorCapacity() *ir.Func {
	if env.Specials["readVectorCapacity"] == nil {
		env.Specials["readVectorCapacity"] = env.Module.NewFunc(
			"readVectorCapacity",
			types.I32,
			ir.NewParam("v", types.NewPointer(voidVector)))

	}

	return env.Specials["readVectorCapacity"]
}

func (env *Environment) getReadVectorWidth() *ir.Func {
	if env.Specials["readVectorWidth"] == nil {
		env.Specials["readVectorWidth"] = env.Module.NewFunc(
			"readVectorWidth",
			types.I32,
			ir.NewParam("v", types.NewPointer(voidVector)))

	}

	return env.Specials["readVectorWidth"]
}

func (env *Environment) getCreateVectorHeader() *ir.Func {
	if env.Specials["createVectorHeader"] == nil {
		env.Specials["createVectorHeader"] = env.Module.NewFunc(
			"createVectorHeader",
			types.NewPointer(voidVector),
			ir.NewParam("length", types.I32),
			ir.NewParam("capacity", types.I32),
			ir.NewParam("width", types.I32))
	}

	return env.Specials["createVectorHeader"]
}
