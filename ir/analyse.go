package ir

import (
	"sundown/sunday/parser"
)

type State struct {
	Package    *string
	Directives []*Directive
	Functions  []*Function
	TypeDefs   []*TypeDef
}

func (p *State) String() string {
	var str string
	for _, directive := range p.Directives {
		str += directive.String() + "\n"
	}

	str += "\n"

	for _, statement := range p.Functions {
		str += statement.String()
	}

	return str
}

func (state *State) Analyse(program *parser.Program) {
	for _, statement := range program.Statements {
		if statement.Directive != nil {
			state.Directives = append(
				state.Directives,
				AnalyseDirective(statement.Directive))
		} else {
			state.Functions = append(
				state.Functions,
				AnalyseStatement(statement.FnDecl))
		}
	}
}
