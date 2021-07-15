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
	Functions     map[Ident]*Function
	TypeDefs      map[Ident]*TypeDef
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

func (state *State) GetFunction(namespace *string, ident *string) *Function {
	if namespace == nil {
		stdfn := state.Functions[Ident{Namespace: "_", Ident: *ident}]
		if stdfn == nil {
			return nil
		} else {

			return stdfn
		}
	} else {

		pkgfn := state.Functions[Ident{Namespace: *namespace, Ident: *ident}]
		if pkgfn == nil {
			return nil
		} else {
			return pkgfn
		}
	}
}

func (state *State) Analyse(program *parser.Program) {
	state.Directives = make(map[*string]*Directive)
	state.Functions = make(map[Ident]*Function)
	state.TypeDefs = make(map[Ident]*TypeDef)

	retid := Ident{Namespace: "_", Ident: "Return"}
	state.Functions[retid] = &Function{Ident: &retid, Takes: &Type{Atomic: "T"}, Gives: &Type{Atomic: "T"}, Body: nil}

	for _, statement := range program.Statements {
		if statement.Directive != nil {
			dir := state.AnalyseDirective(statement.Directive)

			if state.Directives[dir.Class] == nil {
				state.Directives[dir.Class] = dir
			} else {
				panic("Directive already present")
			}

		} else {
			// TODO: do a first pass, ingesting only the declarationss
			def := state.AnalyseStatement(statement.FnDecl)
			if state.Functions[*def.Ident] == nil {
				state.Functions[*def.Ident] = def
			} else {
				panic("Function already defined in package")
			}
		}
	}

	found := state.GetFunction(state.PackageIdent, state.EntryIdent)

	if found == nil {
		panic("Entry function not defined")
	} else {
		state.EntryFunction = found

	}

	state.PrintFunctions()
}
