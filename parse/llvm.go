package parse

import "github.com/llir/llvm/ir/types"

func (t *Type) AsLLType() types.Type {
	if t.Atomic != nil {
		// Type already calculated
		return t.LLType
	} else if t.Vector != nil {
		// Recurse until atomic type(s) found
		// Vectors are always of the form <length | capacity | *data>
		return types.NewStruct(
			types.I64,                             // length
			types.I64,                             // capacity
			types.NewPointer(t.Vector.AsLLType())) // *data
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

// Used for calloc'ing vectors
func (t *Type) WidthInBytes() int64 {
	if t.Atomic != nil {
		return t.Width
	} else {
		return 8
	}
}

func (e *Expression) Type() *Type {
	if e.Atom != nil {
		return e.Atom.TypeOf
	} else if a := e.Application; a != nil {
		// Implement T -> T transfer
		// ([T], Int) -> T i.e. ([Char], Int) -> Char
		if a.Argument.TypeOf.Atomic != nil && *a.Argument.TypeOf.Atomic == "T" {
			// Poly
		} else {

		}
	}

	return nil
}
