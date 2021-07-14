package analysis

import "sundown/sunday/parser"

type Directive struct {
	Class       string
	Instruction struct {
		String, Ident string
		Number        float64
	}
}

func AnalyseDirective(directive *parser.Directive) (d *Directive) {
	d = &Directive{Class: *directive.Class}

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
