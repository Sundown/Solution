package apotheosis

import (
	"fmt"

	"github.com/sundown/solution/prism"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

func (env *Environment) compileFunction(f *prism.Function) value.Value {
	if mfn, ok := env.LLMonadicFunctions[(*f).LLVMise()]; ok {
		return mfn
	} else if dfn, ok := env.LLDyadicFunctions[(*f).LLVMise()]; ok {
		return dfn
	}

	prism.Panic("Not found")
	panic(nil)
}

func (env *Environment) compileBlock(body *[]prism.Expression) {
	// Block is just an expression[]
	for _, stmt := range *body {
		env.compileExpression(&stmt)
	}
}

func (env *Environment) declareDyadicFunction(fn prism.DyadicFunction) *ir.Func {
	return env.Module.NewFunc(
		fn.LLVMise(),
		toReturn(fn.Type()),
		env.toParam(fn.AlphaType), env.toParam(fn.OmegaType))
}

func (env *Environment) declareMonadicFunction(fn prism.MonadicFunction) *ir.Func {
	return env.Module.NewFunc(
		fn.LLVMise(),
		toReturn(fn.Type()),
		env.toParam(fn.OmegaType))
}

func (env *Environment) compileDyadicFunction(fn prism.DyadicFunction) *ir.Func {
	env.CurrentFunction = env.LLDyadicFunctions[fn.LLVMise()]
	env.CurrentFunctionIR = fn

	env.Block = env.newBlock(env.CurrentFunction)
	env.compileBlock(&fn.Body)

	if fn.Returns.Kind() == prism.VoidType.ID {
		env.Block.NewRet(nil)
	}

	if fn.ShouldInline() {
		env.CurrentFunction.FuncAttrs = []ir.FuncAttribute{enum.FuncAttrAlwaysInline}
	}

	return env.CurrentFunction
}

func (env *Environment) compileMonadicFunction(fn prism.MonadicFunction) *ir.Func {
	env.CurrentFunction = env.LLMonadicFunctions[fn.LLVMise()]
	env.CurrentFunctionIR = fn

	env.Block = env.newBlock(env.CurrentFunction)
	env.compileBlock(&fn.Body)

	if fn.Returns.Kind() == prism.VoidType.ID {
		env.Block.NewRet(nil)
	}

	return env.CurrentFunction
}

// Complex types decay to pointers, atomic types do not
func toReturn(t prism.Type) (typ types.Type) {
	typ = t.Realise()

	if prism.IsVector(t) {
		typ = types.NewPointer(t.Realise())
	}

	return typ
}

// Handle void parameters and add pointers to complex types
func (env *Environment) toParam(t prism.Type) (typ *ir.Param) {
	if t.Kind() == prism.VoidType.ID {
		typ = nil
	} else if _, ok := t.(prism.AtomicType); !ok {
		typ = ir.NewParam("", types.NewPointer(t.Realise()))
	} else {
		typ = ir.NewParam("", t.Realise())
	}

	typ.SetName(fmt.Sprint(env.newID()))

	return typ
}
