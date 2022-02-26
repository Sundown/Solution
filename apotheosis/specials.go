package apotheosis

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
)

func (env *Environment) GetCalloc() *ir.Func {
	if env.Specials["calloc"] == nil {
		env.Specials["calloc"] = env.Module.NewFunc(
			"calloc",
			types.I8Ptr,
			ir.NewParam("size", types.I32),
			ir.NewParam("count", types.I32))
	}

	return env.Specials["calloc"]
}

func (env *Environment) GetMemcpy() *ir.Func {
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

func (env *Environment) GetExit() *ir.Func {
	if env.Specials["exit"] == nil {
		env.Specials["exit"] = env.Module.NewFunc(
			"exit",
			types.Void,
			ir.NewParam("", types.I32))
	}

	return env.Specials["exit"]
}

func (env *Environment) GetMaxDouble() *ir.Func {
	if env.Specials["max.f64"] == nil {
		env.Specials["max.f64"] = env.Module.NewFunc(
			"llvm.maxnum.double",
			types.Double,
			ir.NewParam("", types.Double),
			ir.NewParam("", types.Double))
		env.Specials["max.f64"].Sig.Variadic = true
	}

	return env.Specials["max.f64"]
}

func (env *Environment) GetMinDouble() *ir.Func {
	if env.Specials["min.f64"] == nil {
		env.Specials["min.f64"] = env.Module.NewFunc(
			"llvm.minnum.double",
			types.Double,
			ir.NewParam("", types.Double),
			ir.NewParam("", types.Double))
		env.Specials["min.f64"].Sig.Variadic = true
	}

	return env.Specials["min.f64"]
}

func (env *Environment) GetPrintf() *ir.Func {
	if env.Specials["printf"] == nil {
		env.Specials["printf"] = env.Module.NewFunc(
			"printf",
			types.I32,
			ir.NewParam("", types.I8Ptr))
		env.Specials["printf"].Sig.Variadic = true
	}

	return env.Specials["printf"]
}
