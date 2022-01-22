package apotheosis

import (
	"sundown/solution/prism"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
)

func (env *Environment) CompileBlock(body *[]prism.Expression) {
	// Block is just an expression[]
	for _, stmt := range *body {
		env.CompileExpression(&stmt)
	}
}

func (env *Environment) DeclareDyadicFunction(fn prism.DyadicFunction) *ir.Func {
	return env.Module.NewFunc(
		fn.LLVMise(),
		ToReturn(fn.Type()),
		ToParam(fn.AlphaType), ToParam(fn.OmegaType))
}

func (env *Environment) DeclareMonadicFunction(fn prism.MonadicFunction) *ir.Func {
	return env.Module.NewFunc(
		fn.LLVMise(),
		ToReturn(fn.Type()),
		ToParam(fn.OmegaType))
}

func (env *Environment) CompileDyadicFunction(fn prism.DyadicFunction) *ir.Func {
	env.CurrentFunction = env.LLDyadicFunctions[fn.LLVMise()]
	env.CurrentFunctionIR = fn

	env.Block = env.CurrentFunction.NewBlock("")
	env.CompileBlock(&fn.Body)

	if fn.Returns.Kind() == prism.VoidType.ID {
		env.Block.NewRet(nil)
	}

	return env.CurrentFunction
}

func (env *Environment) CompileMonadicFunction(fn prism.MonadicFunction) *ir.Func {
	env.CurrentFunction = env.LLMonadicFunctions[fn.LLVMise()]
	env.CurrentFunctionIR = fn

	env.Block = env.CurrentFunction.NewBlock("")
	env.CompileBlock(&fn.Body)

	if fn.Returns.Kind() == prism.VoidType.ID {
		env.Block.NewRet(nil)
	}

	return env.CurrentFunction
}

// Complex types decay to pointers, atomic types do not
func ToReturn(t prism.Type) (typ types.Type) {
	if t.Kind() == prism.VoidType.ID {
		typ = types.Void
	} else if _, ok := t.(prism.AtomicType); !ok {
		typ = types.NewPointer(t.Realise())
	} else {
		typ = t.Realise()
	}

	return typ
}

// Handle void parameters and add pointers to complex types
func ToParam(t prism.Type) (typ *ir.Param) {
	if t.Kind() == prism.VoidType.ID {
		typ = nil
	} else if _, ok := t.(prism.AtomicType); !ok {
		typ = ir.NewParam("", types.NewPointer(t.Realise()))
	} else {
		typ = ir.NewParam("", t.Realise())
	}

	return typ
}
