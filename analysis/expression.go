package analysis

import "sundown/sunday/parser"

type Expression struct {
	TypeOf      *Type
	Application *Application
	Binary      *Binary
	Atom        *Atom
	Type        *Type
	Block       []*Expression
}

func AnalyseExpression(expression *parser.Expression) (e *Expression) {
	switch {
	case expression.Type != nil:
		e.Type = AnalyseType(expression.Type)
		e.TypeOf = &Type{Atomic: "Type"} /* type of type type obviously */
	case expression.Primary != nil:
		e = &Expression{Atom: AnalyseAtom(expression.Primary)}
		e.Type = e.Atom.TypeOf
	case expression.Application != nil:
		e = &Expression{Application: AnalyseApplication(expression.Application)}
		e.Type = e.Application.Function.Gives
	}

	/* 	if expression.Binary != nil {
		fmt.Println("-----Binary-----")
		e = &Expression{
			TypeOf: e.Type,
			Binary: AnalyseBinary(expression.Binary)}
	} */

	return e
}
