package ir

import (
	"fmt"
	"sundown/sunday/parser"
)

type Directive struct {
	Class       *string
	Instruction struct {
		String, Ident *string
		Number        *float64
	}
}

func (d *Directive) String() string {
	var instr string
	if d.Instruction.String != nil {
		instr = *d.Instruction.String
	} else if d.Instruction.Number != nil {
		instr = fmt.Sprint(d.Instruction.Number)
	} else if d.Instruction.Ident != nil {
		instr = *d.Instruction.Ident
	}

	return "@" + *d.Class + " " + instr + ";"
}

func (state *State) AnalyseDirective(directive *parser.Directive) (d *Directive) {
	d = &Directive{Class: directive.Class}

	/* professional gopher moment */
	if directive.Instr.Ident != nil {
		d.Instruction.Ident = directive.Instr.Ident
	} else if directive.Instr.String != nil {
		d.Instruction.String = directive.Instr.String
	} else if directive.Instr.Number != nil {
		d.Instruction.Number = directive.Instr.Number
	}

	switch *d.Class {
	case "Package":
		if d.Instruction.Ident == nil {
			panic("Package defined with wrong type")
		}

		if d.IsFoundational() {
			panic(`"` + *d.Instruction.Ident + `" is a reserved package name`)
		}

		state.PackageIdent = d.Instruction.Ident
	case "Entry":
		if d.Instruction.Ident == nil {
			panic("Entry defined with wrong type")
		}
		state.EntryIdent = d.Instruction.Ident
	default:
		panic("Unknown directive")
	}

	return d
}

func (d *Directive) IsFoundational() bool {
	return *d.Instruction.Ident == "_" || *d.Instruction.Ident == "foundation" || *d.Instruction.Ident == "se"
}
