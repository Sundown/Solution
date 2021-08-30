package parse

import (
	"fmt"
	"sundown/solution/lex"
	"sundown/solution/util"
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

func (state *State) AnalyseDirective(directive *lex.Directive) {
	d := &Directive{Class: directive.Class}

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
		if state.PackageIdent != nil {
			util.Error("Packagename already defined\n" + directive.Pos.String()).Exit()
		}

		if d.Instruction.Ident == nil {
			panic("Package defined with wrong type")
		}

		if IsReserved(*d.Instruction.Ident) {
			panic(`"` + *d.Instruction.Ident + `" is a reserved package name`)
		}

		state.PackageIdent = d.Instruction.Ident
	case "Entry":
		if state.EntryFunction != nil {
			panic("Entry already defined")
		}

		if d.Instruction.Ident == nil {
			panic("Entry defined with wrong type")
		}

		state.EntryIdent = d.Instruction.Ident
	default:
		panic("Unknown directive")
	}
}
