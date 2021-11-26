package parse

import (
	"fmt"
	"strings"
	"sundown/solution/lexer"
	"sundown/solution/oversight"
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

func (state *State) AnalyseDirective(directive *lexer.Directive) {
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
			oversight.Error(
				"Packagename already defined\n" +
					directive.Pos.String()).Exit()
		}

		if d.Instruction.Ident == nil {
			oversight.Error(
				oversight.Yellow("@Package") +
					" requires ident.\n" +
					directive.Pos.String()).Exit()
		}

		if IsReserved(*d.Instruction.Ident) {
			oversight.Error("Identifier \"" +
				oversight.Yellow(*d.Instruction.Ident) +
				"\" is reserved by the compiler.\n" +
				directive.Pos.String()).Exit()
		}
		state.PackageIdent = oversight.Ref(strings.TrimSpace(*d.Instruction.Ident))
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
