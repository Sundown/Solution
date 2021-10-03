package parse

import (
	"fmt"
	"strconv"
	"sundown/solution/lexer"
	"sundown/solution/util"
)

type Atom struct {
	TypeOf   *Type
	Tuple    []*Expression
	Vector   []*Expression
	Param    *bool // unused
	Function *Function
	Int      *int64
	Nat      *uint64
	Real     *float64
	Bool     *bool
	Char     *int8
	Noun     *Ident
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
	case a.Function != nil:
		return a.Function.SigString()
	}

	return "_"
}

func (t *Type) AsVector() *Type {
	return &Type{Vector: t}
}

func (state *State) AnalyseAtom(primary *lexer.Primary) (a *Atom) {
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
		for index, expr := range primary.Vec {
			e := state.AnalyseExpression(expr)
			// all elements must be of same type

			if index > 0 && !vec[index-1].TypeOf.AsLLType().Equal(e.TypeOf.AsLLType()) {
				util.Error(
					"Element of type " +
						util.Yellow(e.TypeOf.String()) +
						" diverges from vector type of " +
						util.Yellow(vec[index-1].TypeOf.String()) +
						"\n" + expr.Pos.String()).Exit()
			}

			vec = append(vec, e)
		}

		a = &Atom{TypeOf: vec[0].TypeOf.AsVector(), Vector: vec}
	case primary.Int != nil:
		a = &Atom{TypeOf: &IntType, Int: primary.Int}
	case primary.Real != nil:
		a = &Atom{TypeOf: &RealType, Real: primary.Real}
	case primary.Char != nil:
		v, err := strconv.Unquote(*primary.Char)
		if err != nil {
			util.Error("Invalid character literal '" +
				util.Yellow(*primary.Char) +
				"'.\n" + primary.Pos.String()).Exit()
		}

		t := int8(v[0])
		a = &Atom{TypeOf: &CharType, Char: &t}
	case primary.String != nil:
		arr := []*Expression{}
		for _, char := range *primary.String {
			t := int8(char)
			a := &Atom{TypeOf: &CharType, Char: &t}
			arr = append(arr, &Expression{TypeOf: &CharType, Atom: a})
		}

		a = &Atom{TypeOf: (&CharType).AsVector(), Vector: arr}
	case primary.Noun != nil:
		var p bool
		if *primary.Noun.Ident == "True" {
			p = true
			a = &Atom{TypeOf: &BoolType, Bool: &p}
		} else if *primary.Noun.Ident == "False" {
			p = false
			a = &Atom{TypeOf: &BoolType, Bool: &p}
		} else {
			a = state.GetNoun(primary.Noun)
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
