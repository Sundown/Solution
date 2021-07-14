package ir

import (
	"sundown/sunday/parser"
)

type Program struct {
	Statements []*Function
	Directives []*Directive
}

type Binary struct {
	TypeOf *Type
	Op     string
	Left   *Expression
	Right  *Expression
}

func (p *Program) String() string {
	var str string
	for _, directive := range p.Directives {
		str += directive.String() + ";\n"
	}

	str += "\n"

	for _, statement := range p.Statements {
		str += statement.String()
	}

	return str
}

func (b *Binary) String() string {
	return b.Left.String() + b.Op + b.Right.String()
}

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

func AnalyseBinary(binary *parser.Expression) (b *Binary) {
	/* TODO: return an *Application at some point because Binary object
	is somewhat of a hack */
	b.Op = *binary.Op
	return b
}
