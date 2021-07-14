package analysis

import "sundown/sunday/parser"

type Type struct {
	Atomic string
	Vector *Type
	Struct []*Type
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
	case typ.Struct != nil:
		for _, temp := range typ.Struct {
			t.Struct = append(t.Struct, AnalyseType(temp))
		}
	}

	return t
}
