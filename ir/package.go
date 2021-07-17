package ir

import (
	"fmt"
	"sundown/sunday/parser"
)

type State struct {
	PackageIdent  *string
	EntryIdent    *string
	EntryFunction *Function
	Directives    map[*string]*Directive
	Functions     map[IdentKey]*Function
	NounDefs      map[IdentKey]*Atom
	TypeDefs      map[IdentKey]*Type
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

func (state *State) PrintFunctions() {
	var str string
	for _, function := range state.Functions {
		str += function.SigString() + "\n"
	}

	fmt.Println(str)
}

func (state *State) Analyse(program *parser.Program) {
	state.Directives = make(map[*string]*Directive)
	state.Functions = make(map[IdentKey]*Function)
	state.TypeDefs = make(map[IdentKey]*Type)
	state.NounDefs = make(map[IdentKey]*Atom)
	state.PopulateTypes(BaseTypes)

	// Temporary
	und := "_"
	ret := "Return"
	retid := Ident{Namespace: &und, Ident: &ret}
	state.Functions[retid.AsKey()] = &Function{Ident: &retid, Takes: AtomicType("T"), Gives: AtomicType("T"), Body: nil}

	for _, statement := range program.Statements {
		if statement.Directive != nil {
			dir := state.AnalyseDirective(statement.Directive)

			if state.Directives[dir.Class] == nil {
				state.Directives[dir.Class] = dir
			} else {
				panic("Directive already present")
			}

		}
	}

	for _, statement := range program.Statements {
		if statement.TypeDecl != nil {
			def := state.AnalyseType(statement.TypeDecl.Type)
			key := IdentKey{Namespace: *state.PackageIdent, Ident: *statement.TypeDecl.Ident}
			if state.TypeDefs[key] == nil {
				state.TypeDefs[key] = def
			} else {
				panic("Type already defined in package")
			}
		} else if statement.NounDecl != nil {
			state.AnalyseNounDecl(statement.NounDecl)
		}
	}

	for _, statement := range program.Statements {
		if statement.FnDecl != nil {
			// TODO: do a first pass, ingesting only the declarationss
			def := state.AnalyseStatement(statement.FnDecl)
			if state.Functions[def.Ident.AsKey()] == nil {
				state.Functions[def.Ident.AsKey()] = def
			} else {
				panic("Function already defined in package")
			}
		}

	}

	found := state.GetFunction(&Ident{Namespace: state.PackageIdent, Ident: state.EntryIdent})

	if found == nil {
		panic("Entry function not defined")
	} else {
		state.EntryFunction = found

	}

	state.PrintFunctions()
}
