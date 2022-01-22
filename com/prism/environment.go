package prism

import (
	"github.com/llir/llvm/ir/types"
)

func NewEnvironment() *Environment {
	var env Environment
	env.MonadicFunctions = make(map[Ident]*MonadicFunction)
	env.DyadicFunctions = make(map[Ident]*DyadicFunction)
	env.Types = make(map[Ident]Type)

	env.Types[ParseIdent("Int")] = AtomicType{
		ID:           TypeInt,
		WidthInBytes: 8,
		Name:         ParseIdent("Int"),
		Actual:       types.I64,
	}
	env.Types[ParseIdent("Real")] = AtomicType{
		ID:           TypeReal,
		WidthInBytes: 8,
		Name:         ParseIdent("Real"),
		Actual:       types.Double,
	}
	env.Types[ParseIdent("Char")] = AtomicType{
		ID:           TypeChar,
		WidthInBytes: 1,
		Name:         ParseIdent("Char"),
		Actual:       types.I8,
	}
	env.Types[ParseIdent("String")] = AtomicType{
		ID:           TypeString,
		WidthInBytes: 12, // TODO
		Name:         ParseIdent("String"),
		Actual:       types.I8Ptr,
	}
	env.Types[ParseIdent("Bool")] = AtomicType{
		ID:           TypeBool,
		WidthInBytes: 1,
		Name:         ParseIdent("Bool"),
		Actual:       types.I1,
	}
	env.Types[ParseIdent("Void")] = AtomicType{
		ID:           TypeVoid,
		WidthInBytes: 0,
		Name:         ParseIdent("Void"),
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
