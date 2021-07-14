package analysis

import (
	"fmt"
	"sundown/sunday/parser"
)

type Atom struct {
	TypeOf *Type
	Struct []*Expression
	Vector []*Expression
	Int    int64
	Nat    uint64
	Real   float64
	Bool   bool
	String string
	Noun   string
	Param  uint
}

func AnalyseAtom(primary *parser.Primary) (a *Atom) {
	switch {
	case primary.Tuple != nil:
		var types []*Type
		var strct []*Expression
		for _, expr := range primary.Tuple {
			e := AnalyseExpression(expr)
			types = append(types, e.TypeOf)
			strct = append(strct, e)
		}

		a = &Atom{TypeOf: &Type{Struct: types}, Struct: strct}
	case primary.Vec != nil:
		var vec []*Expression
		for index, expr := range primary.Vec {
			e := AnalyseExpression(expr)
			/* all elements must be of same type */
			if index > 0 && vec[index-1].TypeOf != e.TypeOf {
				panic("Analysis: Atom: Vector: divergent type at position: " + fmt.Sprint(index))
			}

			vec = append(vec, e)
		}

		a = &Atom{TypeOf: vec[0].TypeOf, Vector: vec}
	case primary.Int != nil:
		a = &Atom{TypeOf: &Type{Atomic: "Int"}, Int: *primary.Int}
	case primary.Real != nil:
		a = &Atom{TypeOf: &Type{Atomic: "Real"}, Real: *primary.Real}
	case primary.Bool != nil:
		/* TODO: add a third bool state "Maybe", maybe */
		var b bool
		if *primary.Bool == "True" {
			b = true
		} else {
			b = false
		}

		a = &Atom{TypeOf: &Type{Atomic: "Bool"}, Bool: b}
	case primary.String != nil:
		/* TODO: strings might need their "" cut off each end because parser sometimes leaves them */
		a = &Atom{TypeOf: &Type{Atomic: "String"}, String: *primary.String}
	case primary.Noun != nil:
		/* TODO: add lookup because noun is actually referencing something, has a type etc */
		a = &Atom{TypeOf: &Type{Atomic: "Noun"}, Noun: *primary.Noun}
	case primary.Param != nil:
		/* TODO: add param index if it exists, needs parser modification too */
		a = &Atom{TypeOf: &Type{Atomic: "Param"}, Param: 0}
	}

	return a
}
