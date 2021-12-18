package subtle

import (
	"sundown/solution/palisade"
)

type Expression struct {
	TypeOf   *Type
	Monadic  *MonadicApplication
	Dyadic   *DyadicApplication
	Morpheme *Morpheme
	Block    []*Expression
}

func (e *Expression) String() string {
	if e.Application != nil {
		return e.Application.String()
	} else if e.Morpheme != nil {
		return e.Morpheme.String()
	} else if e.Block != nil {
		var str string
		for _, v := range e.Block {
			str += "  " + v.String() + ";\n"
		}

		return str
	}

	return "//"
}

func (state *State) AnalyseExpression(expr *palisade.Expression) (e *Expression) {
	if expr.Dyadic.Next == nil {
		m := state.AnalyseMonadicApplication(expr.Monadic)
		//e = &Expression{

	}
	/* if expression.Morpheme != nil {
		e = &Expression{Morpheme: state.AnalyseMorpheme(expression.Morpheme)}
		e.TypeOf = e.Morpheme.TypeOf
	} else if expression.Application != nil {
		e = &Expression{Application: state.AnalyseApplication(expression.Application)}
		e.TypeOf = e.Application.Function.Gives
	}

	return e */
	return nil
}
