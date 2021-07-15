package ir

import (
	"sundown/sunday/parser"
)

type Type struct {
	Atomic string
	Vector *Type
	Tuple  []*Type
}

type TypeDef struct {
	Ident *Ident
	Type  *Type
}

func (t *Type) String() string {
	switch {
	case t.Atomic != "":
		return t.Atomic
	case t.Vector != nil:
		return "[" + t.Vector.String() + "]"
	case t.Tuple != nil:
		var str string
		for _, elm := range t.Tuple {
			str += "," + elm.String()
		}

		return "(" + str[2:] + ")"
	}

	return ""
}

func (state *State) AnalyseType(typ *parser.Type) (t *Type) {
	switch {
	case typ.Primative != nil:
		var namespace string
		if typ.Primative.Namespace == nil {
			namespace = "_"
		} else {
			namespace = *typ.Primative.Namespace
		}

		temp := state.TypeDefs[Ident{Namespace: namespace, Ident: *typ.Primative.Ident}]

		if temp == nil {
			temp = state.TypeDefs[Ident{Namespace: *state.PackageIdent, Ident: *typ.Primative.Ident}]
		}

		if temp == nil {
			panic(`Type "` + *typ.Primative.Ident + `" not found `)
		}

		t = temp.Type
	case typ.Vector != nil:
		t = &Type{Vector: state.AnalyseType(typ.Vector)}
	case typ.Tuple != nil:
		t = &Type{}
		for _, temp := range typ.Tuple {
			t.Tuple = append(t.Tuple, state.AnalyseType(temp))
		}
	}

	return t
}

func (state *State) AnalyseTypeDecl(decl *parser.TypeDecl) (t *TypeDef) {
	t = &TypeDef{
		Ident: &Ident{Namespace: *state.PackageIdent, Ident: *decl.Ident},
		Type:  state.AnalyseType(decl.Type),
	}

	return t
}
