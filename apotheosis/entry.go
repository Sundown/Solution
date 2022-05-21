package apotheosis

// TODO use official LLVM bindings
import (
	"github.com/sundown/solution/prism"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
)

// Environment inherited from prism.Environment, enforced pointer
type Environment struct {
	*prism.Environment
}

// Compile parsed AST to LLVM IR, including datalayouts for NVVM
func Compile(penv *prism.Environment) *prism.Environment {
	prism.Verbose("Init compiler")

	env := &Environment{penv}
	env.ApotheosisIter = 0
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

	// For CUDA
	//env.Module.TargetTriple = "nvptx64-unknown-cuda"
	//env.Module.DataLayout = "e-p:64:64:64-i1:8:8-i8:8:8-i16:16:16-i32:32:32-i64:64:64-i128:128:128-f32:32:32-f64:64:64-v16:16:16-v32:32:32-v64:64:64-v128:128:128-n16:32:64"

	env.
		declareFunctions().
		buildFunctions().
		initMain()

	return env.Environment
}

func (env *Environment) declareFunctions() *Environment {
	for _, fn := range env.DyadicFunctions {
		if fn.Attrs().Special {
			continue
		}

		if _, ok := (*fn).OmegaType.(prism.Universal); ok {
			continue
		}

		if _, ok := (*fn).AlphaType.(prism.Universal); ok {
			continue
		}

		env.LLDyadicFunctions[fn.LLVMise()] = env.declareDyadicFunction(*fn)
	}
	for _, fn := range env.MonadicFunctions {
		if fn.Attrs().Special {
			continue
		}

		if _, ok := (*fn).OmegaType.(prism.Universal); ok {
			continue
		}

		env.LLMonadicFunctions[fn.LLVMise()] = env.declareMonadicFunction(*fn)
	}

	return env
}

func (env *Environment) buildFunctions() *Environment {
	for _, fn := range env.DyadicFunctions {
		if fn.Attrs().Special {
			continue
		}

		if _, ok := (*fn).OmegaType.(prism.Universal); ok {
			continue
		}

		if _, ok := (*fn).AlphaType.(prism.Universal); ok {
			continue
		}

		env.LLDyadicFunctions[fn.LLVMise()] = env.newDyadicFunction(*fn)
	}

	for _, fn := range env.MonadicFunctions {
		if fn.Attrs().Special {
			continue
		}

		if _, ok := (*fn).OmegaType.(prism.Universal); ok {
			continue
		}

		env.LLMonadicFunctions[fn.LLVMise()] = env.newMonadicFunction(*fn)

	}

	return env
}

func (env *Environment) initMain() *Environment {
	env.CurrentFunction = env.Module.NewFunc("main", types.I32, ir.NewParam("argc", types.I32), ir.NewParam("argv", types.NewPointer(types.I8Ptr)))
	env.Block = env.CurrentFunction.NewBlock("entry")
	/*
		curString := env.Block.NewLoad(
			types.I8Ptr,
			env.Block.NewGetElementPtr(
				types.I8Ptr,
				env.CurrentFunction.Params[1],
				i32(1)))

		strlen := env.Block.NewCall(env.getStrlen(), curString)

		head, body := env.vectorFactory(
			prism.CharType,
			env.Block.NewTrunc(strlen, types.I32))

		env.Block.NewCall(
			env.getMemcpy(),
			body, curString, strlen, constant.NewInt(types.I1, 0)) */

	env.Block.NewCall(env.LLMonadicFunctions[env.EntryFunction.LLVMise()], i64(1))

	env.Block.NewRet(i32(0))

	return env
}

func (env *Environment) insertCallables() {
	env.LLDyadicCallables[","] = prism.MakeDC(env.newInlineAppend, true)
	env.LLDyadicCallables["+"] = prism.MakeDC(env.newInlineDAdd, false)
	env.LLDyadicCallables["-"] = prism.MakeDC(env.newInlineDSub, false)
	env.LLDyadicCallables["×"] = prism.MakeDC(env.newInlineMul, false)
	env.LLDyadicCallables["÷"] = prism.MakeDC(env.newInlineDiv, false)
	env.LLDyadicCallables["*"] = prism.MakeDC(env.newInlinePow, false)
	env.LLDyadicCallables["="] = prism.MakeDC(env.newInlineEqual, false)
	env.LLDyadicCallables["⌈"] = prism.MakeDC(env.newInlineMax, false)
	env.LLDyadicCallables["⌊"] = prism.MakeDC(env.newInlineMin, false)
	env.LLDyadicCallables["∧"] = prism.MakeDC(env.newInlineAnd, false)
	env.LLDyadicCallables["∨"] = prism.MakeDC(env.newInlineOr, false)
	env.LLDyadicCallables["⊃"] = prism.MakeDC(env.newInlineIndex, true)
	env.LLDyadicCallables["⊢"] = prism.MakeDC(env.newInlineRightTacD, false)
	env.LLDyadicCallables["⊣"] = prism.MakeDC(env.newInlineLeftTacD, false)

	env.LLMonadicCallables["*"] = prism.MakeMC(env.newInlineExp, false)
	env.LLMonadicCallables["-"] = prism.MakeMC(env.newInlineMSub, false)
	env.LLMonadicCallables["~"] = prism.MakeMC(env.newInlineNot, false)
	env.LLMonadicCallables["⍳"] = prism.MakeMC(env.newInlineIota, true)
	env.LLMonadicCallables["⊂"] = prism.MakeMC(env.newInlineEnclose, true)
	env.LLMonadicCallables["⊢"] = prism.MakeMC(env.newInlineRightTacM, false)
	env.LLMonadicCallables["Println"] = prism.MakeMC(env.newInlinePrintln, true)
	env.LLMonadicCallables["Print"] = prism.MakeMC(env.newInlinePrint, true)
	env.LLMonadicCallables["Panic"] = prism.MakeMC(env.newInlinePanic, false)
	env.LLMonadicCallables["≢"] = prism.MakeMC(env.newInlineTally, true)
	env.LLMonadicCallables["__Cap"] = prism.MakeMC(env.newInlineCapacity, false)
	env.LLMonadicCallables["⌈"] = prism.MakeMC(env.newInlineCeil, false)
	env.LLMonadicCallables["⌊"] = prism.MakeMC(env.newInlineFloor, false)
}
