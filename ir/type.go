package ir

import "sundown/sunday/parser"

type Type struct {
	Atomic string
	Vector *Type
	Tuple  []*Type
}

type TypeDef struct {
	Ident *string
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

func AnalyseType(typ *parser.Type) (t *Type) {
	switch {
	case typ.Primative != nil:
		/* TODO: actually make this generate proper
		 * type sigs instead of just strings by
		 * looking at typedefs/builtins */
		t = &Type{Atomic: *typ.Primative.Type}
	case typ.Vector != nil:
		t = &Type{Vector: AnalyseType(typ.Vector)}
	case typ.Tuple != nil:
		for _, temp := range typ.Tuple {
			t.Tuple = append(t.Tuple, AnalyseType(temp))
		}
	}

	return t
}
