package parse

import (
	"fmt"
	"strconv"
	"sundown/sunday/lex"
)

type Atom struct {
	TypeOf *Type
	Tuple  []*Expression
	Vector []*Expression
	Param  *bool // unused
	Int    *int64
	Nat    *uint64
	Real   *float64
	Bool   *bool
	Char   *int8
	Noun   *Ident
}

func (a *Atom) String() string {
	switch {
	case a.Int != nil:
		return fmt.Sprint(*a.Int)
	case a.Nat != nil:
		return fmt.Sprint(*a.Nat)
	case a.Real != nil:
		return fmt.Sprint(*a.Real)
	case a.Bool != nil:
		if *a.Bool {
			return "True"
		} else {
			return "False"
		}
	case a.Noun != nil:
		return *a.Noun.Namespace + "::" + *a.Noun.Ident
	case a.Char != nil:
		return "'" + string(rune(*a.Char)) + "'"
	case a.Param != nil:
		return "@"
	case a.Vector != nil:
		var str string
		for _, expr := range a.Vector {
			str += ", " + expr.String()
		}

		return "[" + str[2:] + "]"
	case a.Tuple != nil:
		var str string
		for _, expr := range a.Tuple {
			str += ", " + expr.String()
		}

		return "(" + str[2:] + ")"
	}

	return "_"
}

func (t *Type) AsVector() *Type {
	return &Type{Vector: t}
}

func (state *State) AnalyseAtom(primary *lex.Primary) (a *Atom) {
	switch {
	case primary.Tuple != nil:
		var types []*Type
		var strct []*Expression
		for _, expr := range primary.Tuple {
			e := state.AnalyseExpression(expr)
			types = append(types, e.TypeOf)
			strct = append(strct, e)
		}

		a = &Atom{TypeOf: &Type{Tuple: types}, Tuple: strct}
	case primary.Vec != nil:
		var vec []*Expression
		for _, expr := range primary.Vec {
			e := state.AnalyseExpression(expr)
			/* all elements must be of same type */
			// Can't compare types properly yet
			/* if index > 0 && vec[index-1].TypeOf != e.TypeOf {
				panic("parse: Atom: Vector: divergent type at position: " + fmt.Sprint(index) + "\n" + e.TypeOf.String() + " & " + vec[index-1].TypeOf.String())
			} */

			vec = append(vec, e)
		}

		a = &Atom{TypeOf: vec[0].TypeOf.AsVector(), Vector: vec}
	case primary.Int != nil:
		a = &Atom{TypeOf: BaseType("Int"), Int: primary.Int}
	case primary.Real != nil:
		a = &Atom{TypeOf: BaseType("Real"), Real: primary.Real}
	case primary.Char != nil:
		v, err := strconv.Unquote(*primary.Char)
		if err != nil {
			panic(err)
		}

		t := int8(v[0])
		a = &Atom{TypeOf: BaseType("Char"), Char: &t}
	case primary.String != nil:
		arr := []*Expression{}
		for _, char := range *primary.String {
			t := int8(char)
			a := &Atom{TypeOf: BaseType("Char"), Char: &t}
			arr = append(arr, &Expression{TypeOf: BaseType("Char"), Atom: a})
		}

		a = &Atom{TypeOf: BaseType("Char").AsVector(), Vector: arr}
	case primary.Noun != nil:
		var p bool
		if *primary.Noun.Ident == "True" {
			p = true
			a = &Atom{TypeOf: BaseType("Bool"), Bool: &p}
		} else if *primary.Noun.Ident == "False" {
			p = false
			a = &Atom{TypeOf: BaseType("Bool"), Bool: &p}
		} else {
			a = state.GetNoun(IRIdent(primary.Noun))
		}
	case primary.Param != nil:
		b := true
		a = &Atom{
			TypeOf: state.CurrentFunction.Takes,
			Param:  &b,
		}
	default:
		panic("Was a new type added?")
	}

	return a
}
