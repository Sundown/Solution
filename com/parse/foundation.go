package parse

import (
	"sundown/solution/util"

	"github.com/llir/llvm/ir/types"
)

var und = "_"

var IntType = Type{Atomic: util.Ref("Int"), LLType: types.I64, Width: 8}
var NatType = Type{Atomic: util.Ref("Nat"), LLType: types.I64, Width: 8}
var RealType = Type{Atomic: util.Ref("Real"), LLType: types.Double, Width: 8}
var BoolType = Type{Atomic: util.Ref("Bool"), LLType: types.I1, Width: 1}
var CharType = Type{Atomic: util.Ref("Char"), LLType: types.I8, Width: 4}
var VoidType = Type{Atomic: util.Ref("Void"), LLType: types.Void, Width: 0}
var StringType = Type{
	Vector: &Type{Atomic: util.Ref("Char"), LLType: types.I8},
	LLType: types.NewStruct(types.I32, types.I32, types.I8Ptr), Width: 8,
}

func (state *State) PopulateTypes() {
	state.gulpType(IntType)
	state.gulpType(NatType)
	state.gulpType(RealType)
	state.gulpType(BoolType)
	state.gulpType(CharType)
	state.gulpType(VoidType)
	state.gulpType(StringType)
}

func (state *State) gulpType(t Type) {
	id := Ident{Namespace: &und, Ident: t.Atomic}
	state.TypeDefs[(&id).AsKey()] = &t
}
