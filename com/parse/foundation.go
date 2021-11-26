package parse

import (
	"sundown/solution/oversight"

	"github.com/llir/llvm/ir/types"
)

var und = "_"

var IntType = Type{Atomic: oversight.Ref("Int"), LLType: types.I64, Width: 8}
var NatType = Type{Atomic: oversight.Ref("Nat"), LLType: types.I64, Width: 8}
var RealType = Type{Atomic: oversight.Ref("Real"), LLType: types.Double, Width: 8}
var BoolType = Type{Atomic: oversight.Ref("Bool"), LLType: types.I1, Width: 1}
var CharType = Type{Atomic: oversight.Ref("Char"), LLType: types.I8, Width: 4}
var VoidType = Type{Atomic: oversight.Ref("Void"), LLType: types.Void, Width: 0}
var StringType = Type{
	Vector: &Type{Atomic: oversight.Ref("Char"), LLType: types.I8},
	LLType: types.NewStruct(types.I32, types.I32, types.I8Ptr), Width: 8,
}

func (state *State) PopulateTypes() {
	state.gulpType(IntType)
	state.gulpType(NatType)
	state.gulpType(RealType)
	state.gulpType(BoolType)
	state.gulpType(CharType)
	state.gulpType(VoidType)

	id := Ident{Namespace: &und, Ident: oversight.Ref("String")}
	state.TypeDefs[(&id).AsKey()] = &StringType
}

func (state *State) gulpType(t Type) {
	id := Ident{Namespace: &und, Ident: t.Atomic}
	state.TypeDefs[(&id).AsKey()] = &t
}
