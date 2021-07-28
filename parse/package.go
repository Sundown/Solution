package parse

import (
	"fmt"
	"sundown/sunday/lex"
	"sundown/sunday/util"
)

type State struct {
	PackageIdent    *string
	EntryIdent      *string
	EntryFunction   *Function
	CurrentFunction *Function
	Functions       map[IdentKey]*Function
	NounDefs        map[IdentKey]*Atom
	TypeDefs        map[IdentKey]*Type
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

func (state *State) AddSpecialForm(ident string, takes *Type, gives *Type) *State {
	k := Ident{Namespace: util.Ref("_"), Ident: &ident}
	state.Functions[k.AsKey()] = &Function{
		Ident:   &k,
		Takes:   takes,
		Gives:   gives,
		Body:    nil,
		Special: true,
	}

	return state
}

func (state *State) Parse(program *lex.State) *State {
	state.BuildParserEnv()

	state.
		AddSpecialForm("Return", AtomicType("T"), AtomicType("T")).
		AddSpecialForm("GEP", AtomicType("T"), AtomicType("T")).
		AddSpecialForm("Print", AtomicType("T"), AtomicType("T")).
		AddSpecialForm("Sum", AtomicType("T"), AtomicType("T")).
		AddSpecialForm("Len", AtomicType("[T]"), AtomicType("Int"))

	state.
		CollectDirectives(program).
		ForkStatements(program).
		CollectFunctions(program)

	entry := state.GetFunction(&Ident{Namespace: state.PackageIdent, Ident: state.EntryIdent})

	if entry == nil {
		panic("Entry function not defined")
	} else {
		state.EntryFunction = entry

	}

	return state
}

func (state *State) BuildParserEnv() *State {
	state.Functions = make(map[IdentKey]*Function)
	state.TypeDefs = make(map[IdentKey]*Type)
	state.NounDefs = make(map[IdentKey]*Atom)
	state.PopulateTypes(BaseTypes)

	return state
}

func (state *State) CollectDirectives(p *lex.State) *State {
	for _, statement := range p.Statements {
		if statement.Directive != nil {
			state.AnalyseDirective(statement.Directive)
		}
	}

	return state
}

func (state *State) ForkStatements(p *lex.State) *State {
	// Add types, nouns, and function DECLARATIONS to the state before
	// parsing function bodies to allow referencing before declaration
	for _, statement := range p.Statements {
		if statement.TypeDecl != nil {
			state.AnalyseTypeDecl(statement.TypeDecl)
		} else if statement.NounDecl != nil {
			state.AnalyseNounDecl(statement.NounDecl)
		} else if statement.FnDecl != nil {
			state.AnalyseFnDecl(statement.FnDecl)
		}
	}

	return state
}

func (state *State) CollectFunctions(p *lex.State) *State {
	for _, statement := range p.Statements {
		if statement.FnDecl != nil {
			state.AnalyseFnDef(statement.FnDecl)
		}
	}

	return state
}
