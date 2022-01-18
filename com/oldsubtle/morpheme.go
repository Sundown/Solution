package subtle

import (
	"strconv"
	"sundown/solution/palisade"
	"sundown/solution/prism"
)

func (t *Type) AsVector() *Type {
	return &Type{Vector: t}
}

func (state *State) AnalyseMorpheme(morpheme *palisade.Morpheme) (a prism.Expression) {
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

	case morpheme.Int != nil:
		a = prism.Int{Value: *morpheme.Int}
	case morpheme.Real != nil:
		a = prism.Real{Value: *morpheme.Real}
	case morpheme.Char != nil:
		v, _ := strconv.Unquote(*morpheme.Char)
		/* if err != nil {
			oversight.Error("Invalid character literal '" +
				oversight.Yellow(*morpheme.Char) +
				"'.\n" + morpheme.Pos.String()).Exit()
		} */

		t := int8(v[0])
		a = prism.Char{Value: string(t)}
	case morpheme.String != nil:
		a = prism.String{Value: *morpheme.String}
	/* borked
	 case morpheme.Ident != nil:
			if *morpheme.Ident.Ident == "True" {
				a = prism.Bool{Value: true}
			} else if *morpheme.Ident.Ident == "False" {
				a = prism.Bool{Value: false}
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
		}*/
	default:
		panic("Was a new type added?")
	}

	return a
}
