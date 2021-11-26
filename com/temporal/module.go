package temporal

import (
	"fmt"
	"sundown/solution/lexer"
	"sundown/solution/oversight"
)

type State struct {
	PackageIdent    *string
	EntryIdent      *string
	EntryFunction   *Function
	CurrentFunction *Function
	Functions       map[IdentKey]*Function
	NounDefs        map[IdentKey]*Atom
	TypeDefs        map[IdentKey]*Type
	Imports         []string
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
	k := Ident{Namespace: oversight.Ref("_"), Ident: &ident}
	state.Functions[k.AsKey()] = &Function{
		Ident:   &k,
		Takes:   takes,
		Gives:   gives,
		Body:    nil,
		Special: true,
	}

	return state
}

func (state *State) Parse(program *lexer.State) *State {
	oversight.Verbose("Init parser")
	entry := state.
		BuildParserEnv().
		AddSpecialForm("Return", AtomicType("T"), AtomicType("T")).
		AddSpecialForm("GEP", AtomicType("T"), AtomicType("T")).
		AddSpecialForm("First", AtomicType("T"), AtomicType("T")).
		AddSpecialForm("Second", AtomicType("T"), AtomicType("T")).
		AddSpecialForm("Third", AtomicType("T"), AtomicType("T")).
		AddSpecialForm("Print", AtomicType("T"), AtomicType("Void")).
		AddSpecialForm("Println", AtomicType("T"), AtomicType("Void")).
		AddSpecialForm("Sum", AtomicType("T"), AtomicType("T")).
		AddSpecialForm("Product", AtomicType("T"), AtomicType("T")).
		AddSpecialForm("Len", VectorType(AtomicType("T")), AtomicType("Int")).
		AddSpecialForm("Cap", VectorType(AtomicType("T")), AtomicType("Int")).
		AddSpecialForm("Append", StructType(VectorType(AtomicType("T")), VectorType(AtomicType("T"))), VectorType(AtomicType("T"))).
		AddSpecialForm("Map", StructType(AtomicType("T"), VectorType(AtomicType("T"))), AtomicType("[T]")).
		AddSpecialForm("Foldl", StructType(AtomicType("T"), AtomicType("T"), VectorType(AtomicType("T"))), AtomicType("T")).
		AddSpecialForm("Panic", AtomicType("Int"), AtomicType("Void")).
		AddSpecialForm("Equals", StructType(AtomicType("Int"), AtomicType("Int")), AtomicType("Bool")).
		CollectDirectives(program).
		ForkStatements(program).
		CollectFunctions(program).
		GetFunction(&Ident{Namespace: state.PackageIdent, Ident: state.EntryIdent})

	if entry == nil {
		oversight.Warn("Define program entry-point with directive: " + oversight.Yellow("@Entry <fn>") + ".").Exit()
	} else {
		state.EntryFunction = entry
	}

	return state
}

func (state *State) BuildParserEnv() *State {
	state.Functions = make(map[IdentKey]*Function)
	state.TypeDefs = make(map[IdentKey]*Type)
	state.NounDefs = make(map[IdentKey]*Atom)
	state.PopulateTypes()

	return state
}

func (state *State) CollectDirectives(p *lexer.State) *State {
	for _, statement := range p.Statements {
		if statement.Directive != nil {
			state.AnalyseDirective(statement.Directive)
		}
	}

	return state
}

func (state *State) ForkStatements(p *lexer.State) *State {
	// Add types, nouns, and function DECLARATIONS to the state before
	// parsing function bodies to allow referencing before declaration
	for _, statement := range p.Statements {
		if statement.TypeDecl != nil {
			state.AnalyseTypeDecl(statement.TypeDecl)
		} else if statement.NounDecl != nil {
			state.AnalyseNounDecl(statement.NounDecl)
		} else if statement.FnSig != nil {
			state.AnalyseFnDecl(statement.FnSig)
		}
	}

	return state
}

func (state *State) CollectFunctions(p *lexer.State) *State {
	for _, statement := range p.Statements {
		if statement.FnDef != nil {
			state.AnalyseFnDef(statement.FnDef)
		}
	}

	return state
}