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
	env.Types[Ident{"_", "T"}] = Universal{}

	env.InternBuiltins()

	return &env
}

func (env Environment) Intern(f Function) {
	if fn, ok := f.(MonadicFunction); ok {
		if _, ok := env.MonadicFunctions[fn.Name]; ok {
			panic("Monadic function " + fn.Name.String() + " already exists")
		}

		env.MonadicFunctions[fn.Name] = &fn
	} else if fn, ok := f.(DyadicFunction); ok {
		if _, ok := env.DyadicFunctions[fn.Name]; ok {
			panic("Dyadic function " + fn.Name.String() + " already exists")
		}

		env.DyadicFunctions[fn.Name] = &fn
	}
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

	Panic("Monadic verb " + *v.Ident + " not found")
	panic(nil)
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
