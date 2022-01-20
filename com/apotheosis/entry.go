package apotheosis

import (
	"sundown/solution/prism"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
)

type Environment struct {
	*prism.Environment
}

func Compile(penv *prism.Environment) *prism.Environment {
	prism.Verbose("Init compiler")

	env := &Environment{penv}
	env.Specials = make(map[string]*ir.Func)
	env.LLDFunctions = make(map[string]*ir.Func)
	env.LLMFunctions = make(map[string]*ir.Func)
	env.PanicStrings = make(map[string]*ir.Global)

	env.Module = ir.NewModule()

	env.Module.SourceFilename = env.Output + ".ll"

	env.
		DeclareFunctions().
		CompileFunctions().
		InitMain()

	return env.Environment
}

func (env *Environment) DeclareFunctions() *Environment {
	for _, fn := range env.DFunctions {
		if fn.Special {
			continue
		}

		env.LLDFunctions[fn.LLVMise()] = env.DeclareDFunction(*fn)
	}
	for _, fn := range env.MFunctions {
		if fn.Special {
			continue
		}

		env.LLMFunctions[fn.LLVMise()] = env.DeclareMFunction(*fn)
	}

	return env
}

func (env *Environment) CompileFunctions() *Environment {
	for _, fn := range env.DFunctions {
		if fn.Special {
			continue
		}

		env.LLDFunctions[fn.LLVMise()] = env.CompileDFunction(*fn)
	}
	for _, fn := range env.MFunctions {
		if fn.Special {
			continue
		}

		env.LLMFunctions[fn.LLVMise()] = env.CompileMFunction(*fn)
	}

	return env
}

func (env *Environment) InitMain() *Environment {
	env.CurrentFunction = env.Module.NewFunc("main", types.I32)
	env.Block = env.CurrentFunction.NewBlock("entry")
	env.Block.NewCall(env.LLDFunctions["_::add_[Int],[Int]->Int"], I64(0), I64(1))
	env.Block.NewRet(constant.NewInt(types.I32, 0))

	return env
}
