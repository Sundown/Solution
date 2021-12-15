package subtle

import (
	"sundown/solution/oversight"
	"sundown/solution/palisade"
)

// Tries to find noun in order (defined_namespace or foundation) then package
// order may change in future such that foundation is last
func (state *State) GetNoun(key *palisade.Ident) *Morpheme {
	k := IRIdent(key)
	noun := state.NounDefs[k.AsKey()]

	if noun == nil {
		noun = state.NounDefs[IdentKey{
			Namespace: *state.PackageIdent,
			Ident:     *key.Ident,
		}]

		if noun == nil {
			fn := state.GetFunction(k)
			if fn != nil {
				noun = &Morpheme{TypeOf: fn.Gives, Function: fn}

			} else {
				oversight.Error("Identifier \"" + oversight.Yellow(k.String()) + "\" is not defined in current scope or Foundation.\n" + key.Pos.String()).Exit()
			}
		}
	}

	return noun
}

func (state *State) AnalyseNounDecl(noun *palisade.NounDecl) {
	if IsReserved(*noun.Ident) {
		oversight.Error("Identifier \"" + oversight.Yellow(*noun.Ident) + "\" is reserved by the compiler.\n" + noun.Pos.String()).Exit()
	}

	var temp *Morpheme

	if noun.Value.Noun != nil {
		temp = state.GetNoun(noun.Value.Noun)
	} else if noun.Value.ParamAlpha != nil {
		// ... why
		oversight.Error("Cannot use \"" + oversight.Yellow("@") + "\" (parameter figurative) as R-value in definition.\n" + noun.Pos.String()).Exit()
	} else {
		temp = state.AnalyseMorpheme(noun.Value)
	}

	key := IdentKey{Namespace: *state.PackageIdent, Ident: *noun.Ident}
	if state.NounDefs[key] == nil {
		state.NounDefs[key] = temp
	} else {
		oversight.Error("Noun \"" + oversight.Yellow(*noun.Ident) + "\" is already defined as " + oversight.Yellow(state.NounDefs[key].String()) + ".\n" + noun.Pos.String()).Exit()
	}
}
