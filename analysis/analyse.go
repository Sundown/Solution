package analysis

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
