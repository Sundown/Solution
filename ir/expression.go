package ir

import (
	"sundown/sunday/parser"
)

type Expression struct {
	TypeOf      *Type
	Application *Application
	Atom        *Atom
	Type        *Type
	Block       []*Expression
}

func (e *Expression) String() string {
	if e.Application != nil {
		return e.Application.String()
	} else if e.Type != nil {
		return e.Type.String()
	} else if e.Atom != nil {
		return e.Atom.String()
	} else if e.Block != nil {
		var str string
		for _, v := range e.Block {
			str += "  " + v.String() + ";\n"
		}

		return str
	}

	return ""
}

func (state *State) AnalyseExpression(expression *parser.Expression) (e *Expression) {
	switch {
	case expression.Type != nil:
		e = &Expression{
			TypeOf: &Type{Atomic: "Type"},
			Type:   AnalyseType(expression.Type),
		}
	case expression.Primary != nil:
		e = &Expression{Atom: state.AnalyseAtom(expression.Primary)}
		e.TypeOf = e.Atom.TypeOf
	case expression.Application != nil:
		e = &Expression{Application: state.AnalyseApplication(expression.Application)}
		e.TypeOf = e.Application.Function.Gives
	}

	return e
}
