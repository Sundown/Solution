package analysis

import (
	"fmt"
	"sundown/sunday/parser"
)

func Analyse(program *parser.Program) (output Program) {
	for _, statement := range program.Statements {
		if statement.Directive != nil {
			output.Directives = append(
				output.Directives,
				AnalyseDirective(statement.Directive))
		} else {
			output.Statements = append(
				output.Statements,
				AnalyseStatement(statement.FnDecl))
		}
	}

	return output
}

func AnalyseDirective(directive *parser.Directive) (d *Directive) {
	d.Class = *directive.Class

	/* professional gopher moment */
	if directive.Instr.Ident != nil {
		d.Instruction.Ident = *directive.Instr.Ident
	} else if directive.Instr.String != nil {
		d.Instruction.String = *directive.Instr.String
	} else if directive.Instr.Number != nil {
		d.Instruction.Number = *directive.Instr.Number
	}

	return d
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

func AnalyseStatement(statement *parser.FnDecl) (s *Statement) {
	e := Expression{TypeOf: AnalyseType(statement.Type.Type)}
	for _, expr := range statement.Expressions {
		e.Block = append(e.Block, AnalyseExpression(expr))
	}

	return &Statement{Ident: *statement.Ident, Value: &e}
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

func AnalyseExpression(expression *parser.Expression) (e *Expression) {
	switch {
	case expression.Type != nil:
		e.Type = AnalyseType(expression.Type)
		e.TypeOf = &Type{Atomic: "Type"} /* type of type type obviously */
	case expression.Primary != nil:
		e.Atom = AnalyseAtom(expression.Primary)
	}

	return e
}
