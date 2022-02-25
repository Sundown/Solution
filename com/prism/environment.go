package prism

import (
	"sundown/solution/palisade"

	"github.com/llir/llvm/ir"
)

type Environment struct {
	CurrentlyInlining bool
	Iter              uint
	LexResult         *palisade.PalisadeResult
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
	EntryFunction      MonadicFunction
	Module             *ir.Module
	Block              *ir.Block
	LLDyadicFunctions  map[string]*ir.Func
	LLMonadicFunctions map[string]*ir.Func
	Specials           map[string]*ir.Func
	CurrentFunction    *ir.Func
	CurrentFunctionIR  Expression
	PanicStrings       map[string]*ir.Global
}

func NewEnvironment() *Environment {
	var env Environment
	env.CurrentlyInlining = false
	env.Iter = 0
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
	env.DyadicFunctions[GEPSpecial.Name] = &GEPSpecial
	env.DyadicFunctions[AppendSpecial.Name] = &AppendSpecial
	env.DyadicFunctions[AddSpecial.Name] = &AddSpecial
	env.DyadicFunctions[SubSpecial.Name] = &SubSpecial
	env.DyadicFunctions[MulSpecial.Name] = &MulSpecial
	env.DyadicFunctions[DivSpecial.Name] = &DivSpecial
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
}

func (env Environment) FetchMVerb(v *palisade.Ident) MonadicFunction {
	if found, ok := env.MonadicFunctions[Intern(*v)]; ok {
		return *found
	}

	panic("Monadic verb " + *v.Ident + " not found")
}

func (env Environment) FetchVerb(v *palisade.Ident) Expression {
	if found, ok := env.MonadicFunctions[Intern(*v)]; ok {
		return *found
	} else if found, ok := env.DyadicFunctions[Intern(*v)]; ok {
		return *found
	}

	panic("Verb " + *v.Ident + " not found")
}
