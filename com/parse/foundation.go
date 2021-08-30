package parse

import (
	"sundown/solution/util"

	"github.com/llir/llvm/ir/types"
)

var Reserved = []string{"Int", "Nat", "Real", "Bool", "Char", "Void", "Bits", "String", "Return", "_", "se", "foundation"}

var BaseTypes = []Type{IntType, NatType, RealType, BoolType, CharType, VoidType, BitsType, StringType}

var IntType = Type{Atomic: util.Ref("Int"), LLType: types.I64, Width: 8}
var NatType = Type{Atomic: util.Ref("Nat"), LLType: types.I64, Width: 8}
var RealType = Type{Atomic: util.Ref("Real"), LLType: types.Double, Width: 8}
var BoolType = Type{Atomic: util.Ref("Bool"), LLType: types.I1, Width: 1}
var CharType = Type{Atomic: util.Ref("Char"), LLType: types.I8, Width: 4}
var VoidType = Type{Atomic: util.Ref("Void"), LLType: types.Void, Width: 0}
var BitsType = Type{
	Atomic: util.Ref("Bits"),
	LLType: types.NewStruct(types.I32, types.I32Ptr), Width: 8, // Pointer
}

var StringType = Type{
	Vector: &Type{Atomic: util.Ref("String"), LLType: types.I8},
	LLType: types.NewStruct(types.I32, types.I32, types.I8Ptr), Width: 8, // Pointers
}

func (state *State) PopulateTypes(tarr []Type) {
	id := Ident{Namespace: util.Ref("_"), Ident: util.Ref("Int")}
	state.TypeDefs[id.AsKey()] = &IntType

	id = Ident{Namespace: util.Ref("_"), Ident: util.Ref("Nat")}
	state.TypeDefs[id.AsKey()] = &NatType

	id = Ident{Namespace: util.Ref("_"), Ident: util.Ref("Real")}
	state.TypeDefs[id.AsKey()] = &RealType

	id = Ident{Namespace: util.Ref("_"), Ident: util.Ref("Bool")}
	state.TypeDefs[id.AsKey()] = &BoolType

	id = Ident{Namespace: util.Ref("_"), Ident: util.Ref("Char")}
	state.TypeDefs[id.AsKey()] = &CharType

	id = Ident{Namespace: util.Ref("_"), Ident: util.Ref("Void")}
	state.TypeDefs[id.AsKey()] = &VoidType

	id = Ident{Namespace: util.Ref("_"), Ident: util.Ref("Bits")}
	state.TypeDefs[id.AsKey()] = &BitsType

	id = Ident{Namespace: util.Ref("_"), Ident: util.Ref("String")}
	state.TypeDefs[id.AsKey()] = &StringType
}

func IsReserved(s string) bool {
	for _, r := range Reserved {
		if r == s {
			return true
		}
	}

	return false
}

func (t *Type) In(arr []Type) bool {
	if t.Atomic != nil {
		for _, typ := range arr {
			if *typ.Atomic == *t.Atomic {
				return true
			}
		}
	}

	return false
}

func BaseType(s string) *Type {
	for _, typ := range BaseTypes {
		if *typ.Atomic == s {
			return &typ
		}
	}

	return nil
}
