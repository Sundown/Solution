package parse

import "github.com/llir/llvm/ir/types"

func (t *Type) AsLLType() types.Type {
	if t.Atomic != nil {
		// Type already calculated
		return t.LLType
	} else if t.Vector != nil {
		// Recurse until atomic type(s) found
		// Vectors are always of the form <length | capacity | *data>
		return types.NewStruct(types.I32, types.I32, types.NewPointer(t.Vector.AsLLType()))
	} else if t.Tuple != nil {
		// Recurse each item in tuple
		var lltypes []types.Type
		for _, t := range t.Tuple {
			lltypes = append(lltypes, t.AsLLType())
		}

		return types.NewStruct(lltypes...)
	} else {
		panic("Type is empty")
	}
}

func (t *Type) WidthInBytes() int64 {
	if t.Atomic != nil {
		return t.Width
	} else {
		return 8
	}
}
