package parse

import (
	"fmt"
	"sundown/sunday/lex"
)

type State struct {
	PackageIdent  *string
	EntryIdent    *string
	EntryFunction *Function
	Functions     map[IdentKey]*Function
	NounDefs      map[IdentKey]*Atom
	TypeDefs      map[IdentKey]*Type
}

func (p *State) String() string {
	str := "Package: " + *p.PackageIdent + "\nEntry: " + *p.EntryIdent + "\n\n"

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

func (state *State) Parse(program *lex.State) *State {
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
			state.AnalyseDirective(statement.Directive)
		}
	}

	// Add types, nouns, and function DECLARATIONS to the state before
	// parsing function bodies to allow referencing before declaration
	for _, statement := range program.Statements {
		if statement.TypeDecl != nil {
			state.AnalyseTypeDecl(statement.TypeDecl)
		} else if statement.NounDecl != nil {
			state.AnalyseNounDecl(statement.NounDecl)
		} else if statement.FnDecl != nil {
			state.AnalyseFnDecl(statement.FnDecl)
		}
	}

	for _, statement := range program.Statements {
		if statement.FnDecl != nil {
			state.AnalyseFnDef(statement.FnDecl)
		}
	}

	entry := state.GetFunction(&Ident{Namespace: state.PackageIdent, Ident: state.EntryIdent})

	if entry == nil {
		panic("Entry function not defined")
	} else {
		state.EntryFunction = entry

	}

	return state
}
