package parse

import (
	"sundown/solution/lex"
)

// Tries to find noun in order (defined_namespace or foundation) then package
// order may change in future such that foundation is last
func (state *State) GetNoun(key *Ident) *Atom {
	noun := state.NounDefs[key.AsKey()]

	if noun == nil {
		noun = state.NounDefs[IdentKey{
			Namespace: *state.PackageIdent,
			Ident:     *key.Ident,
		}]
		if noun == nil {
			panic("Noun not defined")
		}
	}

	return noun
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

	key := IdentKey{Namespace: *state.PackageIdent, Ident: *noun.Ident}
	if state.NounDefs[key] == nil {
		state.NounDefs[key] = temp
	} else {
		panic("Noun already defined")
	}
}
