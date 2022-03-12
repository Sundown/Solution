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

// Entry point to Apotheosis codegen, pass prism.Environment with parsed AST
func Compile(penv *prism.Environment) *prism.Environment {
	prism.Verbose("Init compiler")

	env := &Environment{penv}
	env.Specials = make(map[string]*ir.Func)
	env.LLMonadicFunctions = make(map[string]*ir.Func)
	env.LLDyadicFunctions = make(map[string]*ir.Func)
	env.LLMonadicCallables = make(map[string]prism.Callable)
	env.LLDyadicCallables = make(map[string]prism.Callable)
	env.PanicStrings = make(map[string]*ir.Global)

	env.insertCallables()

	env.Module = ir.NewModule()

	env.Module.SourceFilename = env.Output + ".ll"

	env.Module.AttrGroupDefs = append(
		ir.NewModule().AttrGroupDefs,
		&ir.AttrGroupDef{ID: 0, FuncAttrs: []ir.FuncAttribute{enum.FuncAttrAlwaysInline}})

	env.
		DeclareFunctions().
		compileFunctions().
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

		env.LLMonadicFunctions[fn.LLVMise()] = env.declareMonadicFunction(*fn)
	}

	return env
}

func (env *Environment) compileFunctions() *Environment {
	for _, fn := range env.DyadicFunctions {
		if fn.Special {
			continue
		}

		env.LLDyadicFunctions[fn.LLVMise()] = env.compileDyadicFunction(*fn)
	}
	for _, fn := range env.MonadicFunctions {
		if fn.Special {
			continue
		}

		env.LLMonadicFunctions[fn.LLVMise()] = env.compileMonadicFunction(*fn)
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

func (env *Environment) insertCallables() {
	env.LLDyadicCallables[","] = prism.MakeDC(env.compileInlineAppend, true)
	env.LLDyadicCallables["+"] = prism.MakeDC(env.compileInlineAdd, false)
	env.LLDyadicCallables["-"] = prism.MakeDC(env.compileInlineSub, false)
	env.LLDyadicCallables["*"] = prism.MakeDC(env.compileInlineMul, false)
	env.LLDyadicCallables["÷"] = prism.MakeDC(env.compileInlineDiv, false)
	env.LLDyadicCallables["="] = prism.MakeDC(env.compileInlineEqual, false)
	env.LLDyadicCallables["Max"] = prism.MakeDC(env.compileInlineMax, false)
	env.LLDyadicCallables["Min"] = prism.MakeDC(env.compileInlineMin, false)
	env.LLDyadicCallables["&"] = prism.MakeDC(env.compileInlineAnd, false)
	env.LLDyadicCallables["|"] = prism.MakeDC(env.compileInlineAnd, false)
	env.LLDyadicCallables["GEP"] = prism.MakeDC(env.compileInlineIndex, true)
	env.LLDyadicCallables["⊢"] = prism.MakeDC(env.compileInlineRightHook, false)

	env.LLMonadicCallables["Println"] = prism.MakeMC(env.compileInlinePrintln, false)
	env.LLMonadicCallables["Print"] = prism.MakeMC(env.compileInlinePrint, false)
	env.LLMonadicCallables["Panic"] = prism.MakeMC(env.compileInlinePanic, false)
	env.LLMonadicCallables["Len"] = prism.MakeMC(env.readVectorLength, false)
	env.LLMonadicCallables["Cap"] = prism.MakeMC(env.readVectorCapacity, false)
	env.LLMonadicCallables["Max"] = prism.MakeMC(env.compileInlineCeil, false)
	env.LLMonadicCallables["Min"] = prism.MakeMC(env.compileInlineFloor, false)
}
