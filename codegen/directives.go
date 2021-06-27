package codegen

import "sundown/sunday/parser"

func (state *State) Direct(instr *parser.Directive) {
	switch *instr.Class {
	case "Entry":
		{
			if val, ok := state.fns[*instr.Instr]; ok {
				state.Entry = val
			} else {
				panic(*instr.Instr + " not found")
			}
		}
	default:
		{
			panic("Unknown instruction: " + *instr.Class)
		}
	}

}
