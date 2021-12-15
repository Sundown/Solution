package subtle

import (
	"sundown/solution/oversight"
	"sundown/solution/palisade"

	"github.com/alecthomas/repr"
)

type Application struct {
	TypeOf        *Type
	Function      *Function
	ArgumentAlpha *Expression
	ArgumentOmega *Expression
}

func (a *Application) String() string {
	var sig string

	if a.Function.Ident.IsFoundational() {
		sig = *a.Function.Ident.Ident
	} else {
		sig = *a.Function.Ident.Namespace + "::" + *a.Function.Ident.Ident
	}

	return sig + " <" + a.ArgumentAlpha.String() + ", " + a.ArgumentOmega.String() + ">"
}

func (state *State) AnalyseApplication(application *palisade.Application) (s *Application) {
	s = &Application{
		// TODO: change function to expression type for currying purposes in the (far) future
		Function:      state.AnalyseFunction(application.Function),
		ArgumentAlpha: state.AnalyseExpression(application.ParameterAlpha),
		ArgumentOmega: state.AnalyseExpression(application.ParameterOmega),
	}

	/* if !s.ArgumentAlpha.TypeOf.Equals(s.Function.TakesAlpha) || !s.ArgumentOmega.TypeOf.Equals(s.Function.TakesOmega) {
		oversight.Error("(Warning this err msg is probably broken) Trying to call",
			oversight.Yellow(s.Function.SigString()), "with",
			oversight.Yellow(s.ArgumentAlpha.TypeOf.String())+".\n"+
				application.ParameterAlpha.Pos.String()).Exit() // TODO <<<<<-----
	} */

	s.TypeOf = s.Function.Gives

	switch *s.Function.Ident.Ident {
	case "Return":
		if !s.ArgumentAlpha.TypeOf.Equals(state.CurrentFunction.Gives) {
			oversight.Error("Value of type",
				s.ArgumentAlpha.TypeOf.String(),
				"does not match function's return type:",
				state.CurrentFunction.Gives.String()).Exit()
		}
	case "Map":
		if s.ArgumentAlpha.TypeOf.Tuple == nil ||
			s.ArgumentAlpha.TypeOf.Tuple[0] == nil ||
			s.ArgumentAlpha.TypeOf.Tuple[1].Vector == nil {
			oversight.Error("Malformed call to " + oversight.Yellow("Map") + ".\n" + application.Pos.String()).Exit()
		}

		s.TypeOf = s.ArgumentAlpha.Morpheme.Tuple[0].Morpheme.Function.Gives.AsVector()
		s.Function.Gives = s.TypeOf
	case "GEP":
		s.TypeOf = s.ArgumentAlpha.Morpheme.Tuple[0].TypeOf.Vector
		s.Function.Gives = s.TypeOf
	/*case "Foldl":
	if s.Argument.TypeOf.Tuple == nil ||
		s.Argument.Morpheme.Tuple == nil || s.Argument.Morpheme.Tuple[0].Morpheme == nil ||
		s.Argument.Morpheme.Tuple[0].Morpheme.Function == nil ||
		s.Argument.Morpheme.Tuple[0].Morpheme.Function.Takes.Tuple == nil ||
		s.Argument.Morpheme.Tuple[1].Morpheme == nil ||
		s.Argument.Morpheme.Tuple[2].Morpheme.Vector == nil {
		oversight.Error("Malformed call to " + oversight.Yellow("Foldl") + ".\n" + application.Pos.String()).Exit()
	}

	fn_t := s.Argument.Morpheme.Tuple[0].Morpheme.Function.Takes.Tuple
	fn_g := s.Argument.Morpheme.Tuple[0].Morpheme.Function.Gives
	id_i := s.Argument.Morpheme.Tuple[1].Morpheme.TypeOf
	vect := s.Argument.Morpheme.Tuple[2].Morpheme.TypeOf

	if !fn_t[0].Equals(id_i) {
		oversight.Error("Mapping function cannot accept identity in " + oversight.Yellow("Foldl") + ".\n" + application.Pos.String()).Exit()
	}

	if !fn_t[0].Equals(fn_g) {
		oversight.Error("Mapping function does not return the identity type " + oversight.Yellow("Foldl") + ".\n" + application.Pos.String()).Exit()
	}

	if !fn_t[1].Equals(vect) {
		oversight.Error("Mapping function cannot accept vector element type in " + oversight.Yellow("Foldl") + ".\n" + application.Pos.String()).Exit()
	}

	s.TypeOf = fn_g
	s.Function.Gives = s.TypeOf*/
	case "Sum":
		s.TypeOf = s.ArgumentAlpha.Morpheme.TypeOf.Vector
		s.Function.Gives = s.TypeOf
	case "Println":
		s.TypeOf = &VoidType
		s.Function.Gives = s.TypeOf
	case "Equals":
		s.TypeOf = &BoolType
		s.Function.TakesAlpha = &IntType
		s.Function.TakesOmega = &IntType
		s.Function.Gives = s.TypeOf
	case "Product":
		//s.TypeOf = s.Argument.Morpheme.TypeOf.Vector
		s.TypeOf = &IntType
		s.Function.Gives = s.TypeOf
	case "Append":
		s.TypeOf = s.ArgumentAlpha.Morpheme.Tuple[0].TypeOf
		s.Function.Gives = s.TypeOf
	case "First":
		s.TypeOf = s.ArgumentAlpha.Morpheme.Tuple[0].TypeOf
		s.Function.Gives = s.TypeOf
	case "Second":
		s.TypeOf = s.ArgumentAlpha.Morpheme.Tuple[1].TypeOf
		s.Function.Gives = s.TypeOf
	case "Third":
		s.TypeOf = s.ArgumentAlpha.Morpheme.Tuple[2].TypeOf
		s.Function.Gives = s.TypeOf
	}

	repr.Println(s)

	return s
}
