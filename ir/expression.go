package ir

import (
	"sundown/sunday/parser"
)

type Expression struct {
	TypeOf      *Type
	Application *Application
	Binary      *Binary
	Atom        *Atom
	Type        *Type
	Block       []*Expression
}

func (e *Expression) String() string {
	if e.Application != nil {
		return e.Application.String()
	} else if e.Binary != nil {
		return e.Binary.String()
	} else if e.Type != nil {
		return e.Type.String()
	} else if e.Atom != nil {
		return e.Atom.String()
	} else if e.Block != nil {
		var str string
		for _, v := range e.Block {
			str += v.String() + ";\n"
		}

		return str
	}

	return ""
}

func AnalyseExpression(expression *parser.Expression) (e *Expression) {
	switch {
	case expression.Type != nil:
		e.Type = AnalyseType(expression.Type)
		e.TypeOf = &Type{Atomic: "Type"} /* type of type type obviously */
	case expression.Primary != nil:
		e = &Expression{Atom: AnalyseAtom(expression.Primary)}
		e.TypeOf = e.Atom.TypeOf
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
