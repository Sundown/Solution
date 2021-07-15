package foundation

import "github.com/llir/llvm/ir/types"

var (
	// Numerical
	Int  = types.I64
	Nat  = types.I64
	Real = types.Double
	Bool = types.I1

	// Complex
	Str = types.NewStruct(types.I32, types.I8Ptr)

	// Aux
	Void = types.Void
)
