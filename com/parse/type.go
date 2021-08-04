package parse

import (
	"sundown/solution/lex"

	"github.com/llir/llvm/ir/types"
)

type Type struct {
	Atomic *string
	Vector *Type
	Tuple  []*Type
	LLType types.Type
	Width  int64
}

type TypeDef struct {
	Ident *Ident
	Type  *Type
}

func (t *Type) String() string {
	switch {
	case t.Atomic != nil:
		return *t.Atomic
	case t.Vector != nil:
		// String is [Char] but should be printed "String" for convenience
		if t.Vector.Atomic != nil && *t.Vector.Atomic == "String" {
			return "String"
		}

		return "[" + t.Vector.String() + "]"
	case t.Tuple != nil:
		var str string
		for _, elm := range t.Tuple {
			str += ", " + elm.String()
		}

		return "(" + str[2:] + ")"
	}

	return "T"
}

func AtomicType(s string) *Type {
	return &Type{Atomic: &s}
}

func (a *Type) Equals(b *Type) bool {
	if a.Atomic != nil && *a.Atomic == "T" || b.Atomic != nil && *b.Atomic == "T" {
		return true /* ;) */
	}

	if a.Atomic != nil && b.Atomic != nil {
		return *a.Atomic == *b.Atomic
	} else if a.Vector != nil && b.Vector != nil {
		return a.Vector.Equals(b.Vector)
	} else if a.Tuple != nil && b.Tuple != nil {
		for i := range a.Tuple {
			if a.Tuple[i] == nil ||
				b.Tuple[i] == nil ||
				!a.Tuple[i].Equals(b.Tuple[i]) {
				return false
			}
		}

		return true
	} else {
		return false
	}
}

func (state *State) AnalyseTypeDecl(typ *lex.TypeDecl) {
	if IsReserved(*typ.Ident) {
		panic("Trying to assign type to reserved name")
	}

	key := IdentKey{Namespace: *state.PackageIdent, Ident: *typ.Ident}
	if state.TypeDefs[key] == nil {
		state.TypeDefs[key] = state.AnalyseType(typ.Type)
	} else {
		panic("Type already defined")
	}
}

func (state *State) AnalyseType(typ *lex.Type) (t *Type) {
	switch {
	case typ.Primative != nil:
		var namespace string
		if typ.Primative.Namespace == nil {
			namespace = "_"
		} else {
			namespace = *typ.Primative.Namespace
		}

		temp := state.TypeDefs[IdentKey{Namespace: namespace, Ident: *typ.Primative.Ident}]

		if temp == nil {
			temp = state.TypeDefs[IdentKey{Namespace: *state.PackageIdent, Ident: *typ.Primative.Ident}]
		}

		if temp == nil {
			panic(`Type "` + *typ.Primative.Ident + `" not found `)
		}

		t = temp
	case typ.Vector != nil:
		t = &Type{Vector: state.AnalyseType(typ.Vector)}
	case typ.Tuple != nil:
		t = &Type{}
		for _, temp := range typ.Tuple {
			t.Tuple = append(t.Tuple, state.AnalyseType(temp))
		}
	}

	t.LLType = t.AsLLType()
	t.Width = t.WidthInBytes()
	return t
}
