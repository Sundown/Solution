package codegen

import "sundown/sunday/parser"

var packagename string

func (state *State) Direct(instr *parser.Directive) {
	switch *instr.Class {
	case "Entry":
		{
			if instr.Instr.Ident == nil {
				panic("Directive " + *instr.Class + " requires ident")
			}

			if val, ok := state.fns[*instr.Instr.Ident]; ok {
				state.Entry = val
			} else {
				panic(*instr.Instr.Ident + " not found")
			}
		}
	case "Package":
		{
			if instr.Instr.String != nil {
				packagename = *instr.Instr.String
				packagename = packagename[1 : len(packagename)-1]
			} else if instr.Instr.Ident != nil {
				packagename = *instr.Instr.Ident
			} else {
				panic("Directive " + *instr.Class + " requires ident or string")
			}
		}
	default:
		{
			panic("Unknown instruction: " + *instr.Class)
		}
	}

}
