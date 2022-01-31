package prism

import (
	"github.com/llir/llvm/ir/types"
)

func NewEnvironment() *Environment {
	var env Environment
	env.MonadicFunctions = make(map[Ident]*MonadicFunction)
	env.DyadicFunctions = make(map[Ident]*DyadicFunction)
	env.Types = make(map[Ident]Type)

	env.Types[Ident{"_", "Int"}] = AtomicType{
		ID:           TypeInt,
		WidthInBytes: 8,
		Name:         Ident{"_", "Int"},
		Actual:       types.I64,
	}
	env.Types[Ident{"_", "Real"}] = AtomicType{
		ID:           TypeReal,
		WidthInBytes: 8,
		Name:         Ident{"_", "Real"},
		Actual:       types.Double,
	}
	env.Types[Ident{"_", "Char"}] = AtomicType{
		ID:           TypeChar,
		WidthInBytes: 1,
		Name:         Ident{"_", "Char"},
		Actual:       types.I8,
	}
	env.Types[Ident{"_", "String"}] = AtomicType{
		ID:           TypeString,
		WidthInBytes: 12, // TODO
		Name:         Ident{"_", "String"},
		Actual:       types.I8Ptr,
	}
	env.Types[Ident{"_", "Bool"}] = AtomicType{
		ID:           TypeBool,
		WidthInBytes: 1,
		Name:         Ident{"_", "Bool"},
		Actual:       types.I1,
	}
	env.Types[Ident{"_", "Void"}] = AtomicType{
		ID:           TypeVoid,
		WidthInBytes: 0,
		Name:         Ident{"_", "Void"},
		Actual:       types.Void,
	}

	env.MonadicFunctions[ReturnSpecial.Name] = &ReturnSpecial
	env.MonadicFunctions[PrintlnSpecial.Name] = &PrintlnSpecial
	env.MonadicFunctions[PrintSpecial.Name] = &PrintSpecial
	env.MonadicFunctions[LenSpecial.Name] = &LenSpecial
	env.MonadicFunctions[CapSpecial.Name] = &CapSpecial
	env.MonadicFunctions[SumSpecial.Name] = &SumSpecial
	env.MonadicFunctions[ProductSpecial.Name] = &ProductSpecial
	env.DyadicFunctions[GEPSpecial.Name] = &GEPSpecial
	env.DyadicFunctions[AppendSpecial.Name] = &AppendSpecial
	env.DyadicFunctions[AddSpecial.Name] = &AddSpecial
	env.DyadicFunctions[SubSpecial.Name] = &SubSpecial
	env.DyadicFunctions[MulSpecial.Name] = &MulSpecial
	return &env
}

func (e Environment) String() (s string) {
	for _, f := range e.DyadicFunctions {
		s += f.String()
	}
	for _, f := range e.MonadicFunctions {
		s += f.String()
	}

	return
}
