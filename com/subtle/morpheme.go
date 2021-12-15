package subtle

import (
	"fmt"
	"strconv"
	"sundown/solution/oversight"
	"sundown/solution/palisade"
)

type Morpheme struct {
	TypeOf   *Type
	Tuple    []*Expression
	Vector   []*Expression
	Param    *bool
	Function *Function
	Int      *int64
	Nat      *uint64
	Real     *float64
	Bool     *bool
	Char     *int8
	Noun     *Ident
}

func (a *Morpheme) String() string {
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

func (state *State) AnalyseMorpheme(morpheme *palisade.Morpheme) (a *Morpheme) {
	switch {
	/* case morpheme.Tuple != nil:
	var types []*Type
	var strct []*Expression
	for _, expr := range morpheme.Tuple {
		e := state.AnalyseExpression(expr)
		types = append(types, e.TypeOf)
		strct = append(strct, e)
	}

	a = &Morpheme{TypeOf: &Type{Tuple: types}, Tuple: strct}*/
	case morpheme.Vec != nil:
		var vec []*Expression
		for index, expr := range morpheme.Vec {
			e := state.AnalyseExpression(expr)
			// all elements must be of same type

			if index > 0 && !vec[index-1].TypeOf.AsLLType().Equal(e.TypeOf.AsLLType()) {
				oversight.Error(
					"Element of type " +
						oversight.Yellow(e.TypeOf.String()) +
						" diverges from vector type of " +
						oversight.Yellow(vec[index-1].TypeOf.String()) +
						"\n" + expr.Pos.String()).Exit()
			}

			vec = append(vec, e)
		}

		a = &Morpheme{TypeOf: vec[0].TypeOf.AsVector(), Vector: vec}
	case morpheme.Int != nil:
		a = &Morpheme{TypeOf: &IntType, Int: morpheme.Int}
	case morpheme.Real != nil:
		a = &Morpheme{TypeOf: &RealType, Real: morpheme.Real}
	case morpheme.Char != nil:
		v, err := strconv.Unquote(*morpheme.Char)
		if err != nil {
			oversight.Error("Invalid character literal '" +
				oversight.Yellow(*morpheme.Char) +
				"'.\n" + morpheme.Pos.String()).Exit()
		}

		t := int8(v[0])
		a = &Morpheme{TypeOf: &CharType, Char: &t}
	case morpheme.String != nil:
		arr := []*Expression{}
		for _, char := range *morpheme.String {
			t := int8(char)
			a := &Morpheme{TypeOf: &CharType, Char: &t}
			arr = append(arr, &Expression{TypeOf: &CharType, Morpheme: a})
		}

		a = &Morpheme{TypeOf: (&CharType).AsVector(), Vector: arr}
	case morpheme.Noun != nil:
		var p bool
		if *morpheme.Noun.Ident == "True" {
			p = true
			a = &Morpheme{TypeOf: &BoolType, Bool: &p}
		} else if *morpheme.Noun.Ident == "False" {
			p = false
			a = &Morpheme{TypeOf: &BoolType, Bool: &p}
		} else {
			a = state.GetNoun(morpheme.Noun)
		}
	case morpheme.ParamAlpha != nil:
		b := true
		a = &Morpheme{
			TypeOf: state.CurrentFunction.TakesAlpha,
			Param:  &b,
		}
	case morpheme.ParamOmega != nil:
		b := true
		a = &Morpheme{
			TypeOf: state.CurrentFunction.TakesOmega,
			Param:  &b,
		}
	case morpheme.Nil != nil:
		b := int64(0)
		a = &Morpheme{
			Int:    &b,
			TypeOf: &VoidType,
		}
	default:
		panic("Was a new type added?")
	}

	return a
}
