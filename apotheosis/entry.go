package apotheosis

import (
	"github.com/sundown/solution/prism"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
)

type Environment struct {
	*prism.Environment
}

func Compile(penv *prism.Environment) *prism.Environment {
	prism.Verbose("Init compiler")

	env := &Environment{penv}
	env.Specials = make(map[string]*ir.Func)
	env.LLDyadicFunctions = make(map[string]*ir.Func)
	env.LLMonadicFunctions = make(map[string]*ir.Func)
	env.PanicStrings = make(map[string]*ir.Global)

	env.Module = ir.NewModule()

	env.Module.SourceFilename = env.Output + ".ll"

	env.Module.AttrGroupDefs = append(
		ir.NewModule().AttrGroupDefs,
		&ir.AttrGroupDef{ID: 0, FuncAttrs: []ir.FuncAttribute{enum.FuncAttrAlwaysInline}})

	env.
		DeclareFunctions().
		CompileFunctions().
		InitMain()

	return env.Environment
}

func (env *Environment) DeclareFunctions() *Environment {
	for _, fn := range env.DyadicFunctions {
		if fn.Special {
			continue
		}

		env.LLDyadicFunctions[fn.LLVMise()] = env.DeclareDyadicFunction(*fn)
	}
	for _, fn := range env.MonadicFunctions {
		if fn.Special {
			continue
		}

		env.LLMonadicFunctions[fn.LLVMise()] = env.DeclareMonadicFunction(*fn)
	}

	return env
}

func (env *Environment) CompileFunctions() *Environment {
	for _, fn := range env.DyadicFunctions {
		if fn.Special {
			continue
		}

		env.LLDyadicFunctions[fn.LLVMise()] = env.CompileDyadicFunction(*fn)
	}
	for _, fn := range env.MonadicFunctions {
		if fn.Special {
			continue
		}

		env.LLMonadicFunctions[fn.LLVMise()] = env.CompileMonadicFunction(*fn)
	}

	return env
}

func (env *Environment) InitMain() *Environment {
	env.CurrentFunction = env.Module.NewFunc("main", types.I32)
	env.Block = env.CurrentFunction.NewBlock("entry")
	env.Block.NewCall(env.LLMonadicFunctions[env.EntryFunction.LLVMise()], I64(0))
	env.Block.NewRet(constant.NewInt(types.I32, 0))

	return env
}
