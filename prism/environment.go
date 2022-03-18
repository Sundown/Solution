package prism

import (
	"github.com/sundown/solution/palisade"

	"github.com/llir/llvm/ir"
)

func (e Environment) String() (s string) {
	for _, f := range e.DyadicFunctions {
		s += f.String()
	}
	for _, f := range e.MonadicFunctions {
		s += f.String()
	}

	return
}

type Environment struct {
	Iter       uint
	IsPilotRun bool
	LexResult  *palisade.PalisadeResult
	//
	MonadicFunctions map[Ident]*MonadicFunction
	DyadicFunctions  map[Ident]*DyadicFunction
	Types            map[Ident]Type
	//
	EmitFormat   string
	Output       string
	Verbose      *bool
	Optimisation *int64
	File         string
	//
	ApotheosisIter     int
	EntryFunction      MonadicFunction
	Module             *ir.Module
	Block              *ir.Block
	LLDyadicFunctions  map[string]*ir.Func
	LLMonadicFunctions map[string]*ir.Func
	Specials           map[string]*ir.Func

	LLMonadicCallables map[string]Callable
	LLDyadicCallables  map[string]Callable

	CurrentFunction   *ir.Func
	CurrentFunctionIR Expression
	PanicStrings      map[string]*ir.Global
}

func NewEnvironment() *Environment {
	var env Environment
	env.Iter = 0
	env.IsPilotRun = false
	env.MonadicFunctions = make(map[Ident]*MonadicFunction)
	env.DyadicFunctions = make(map[Ident]*DyadicFunction)
	env.Types = make(map[Ident]Type)

	env.Types[Ident{"_", "Int"}] = IntType
	env.Types[Ident{"_", "Real"}] = RealType
	env.Types[Ident{"_", "Char"}] = CharType
	env.Types[Ident{"_", "String"}] = StringType
	env.Types[Ident{"_", "Bool"}] = BoolType
	env.Types[Ident{"_", "Void"}] = VoidType

	env.MonadicFunctions[ReturnSpecial.Name] = &ReturnSpecial
	env.MonadicFunctions[PrintlnSpecial.Name] = &PrintlnSpecial
	env.MonadicFunctions[PrintSpecial.Name] = &PrintSpecial
	env.MonadicFunctions[LenSpecial.Name] = &LenSpecial
	env.MonadicFunctions[CapSpecial.Name] = &CapSpecial
	env.DyadicFunctions[PickSpecial.Name] = &PickSpecial
	env.DyadicFunctions[AppendSpecial.Name] = &AppendSpecial
	env.DyadicFunctions[AddSpecial.Name] = &AddSpecial
	env.DyadicFunctions[SubSpecial.Name] = &SubSpecial
	env.DyadicFunctions[MulSpecial.Name] = &MulSpecial
	env.DyadicFunctions[DivSpecial.Name] = &DivSpecial
	env.DyadicFunctions[RightHook.Name] = &RightHook
	env.DyadicFunctions[MaxSpecial.Name] = &MaxSpecial
	env.DyadicFunctions[MinSpecial.Name] = &MinSpecial
	env.DyadicFunctions[EqSpecial.Name] = &EqSpecial
	env.DyadicFunctions[AndSpecial.Name] = &AndSpecial
	env.DyadicFunctions[OrSpecial.Name] = &OrSpecial
	env.MonadicFunctions[CeilSpecial.Name] = &CeilSpecial
	env.MonadicFunctions[FloorSpecial.Name] = &FloorSpecial
	return &env
}

func (e *Environment) Iterate() int {
	e.Iter++
	return int(e.Iter)
}

func (env Environment) FetchDVerb(v *palisade.Ident) DyadicFunction {
	if found, ok := env.DyadicFunctions[Intern(*v)]; ok {
		return *found
	}

	panic("Dyadic verb " + *v.Ident + " not found")
	//panic(nil)
}

func (env Environment) FetchMVerb(v *palisade.Ident) MonadicFunction {
	if found, ok := env.MonadicFunctions[Intern(*v)]; ok {
		return *found
	}

	panic("Monadic verb " + *v.Ident + " not found")
	//panic(nil)
}

func (env Environment) FetchVerb(v *palisade.Ident) Expression {
	if found, ok := env.MonadicFunctions[Intern(*v)]; ok {
		return *found
	} else if found, ok := env.DyadicFunctions[Intern(*v)]; ok {
		return *found
	}

	Panic("Verb " + *v.Ident + " not found")
	panic(nil)
}
