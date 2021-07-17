package parse

import (
	"sundown/sunday/lex"
)

func (state *State) GetNoun(key *Ident) *Atom {
	if key.Namespace == nil {
		noun := state.NounDefs[IdentKey{Namespace: "_", Ident: *key.Ident}]
		if noun == nil {
			return state.NounDefs[IdentKey{Namespace: *state.PackageIdent, Ident: *key.Ident}]
		} else {
			return noun
		}
	} else {
		return state.NounDefs[IdentKey{Namespace: *key.Namespace, Ident: *key.Ident}]
	}
}

func (state *State) AnalyseNounDecl(noun *lex.NounDecl) {
	if IsReserved(*noun.Ident) {
		panic("Trying to assign noun to a reserved name")
	}

	var temp *Atom

	if noun.Value.Noun != nil {
		temp = state.GetNoun(IRIdent(noun.Value.Noun))
	} else if noun.Value.Param != nil {
		panic("Trying to define noun to param")
	} else {
		temp = state.AnalyseAtom(noun.Value)
	}

	if state.NounDefs[IdentKey{Namespace: *state.PackageIdent, Ident: *noun.Ident}] == nil {
		state.NounDefs[IdentKey{Namespace: *state.PackageIdent, Ident: *noun.Ident}] = temp
	} else {
		panic("Noun already defined")
	}
}
