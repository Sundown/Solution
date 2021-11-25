package parse

import (
	"sundown/solution/lexer"
	"sundown/solution/util"
)

type Application struct {
	TypeOf   *Type
	Function *Function
	Argument *Expression
}

func (a *Application) String() string {
	var sig string

	if a.Function.Ident.IsFoundational() {
		sig = *a.Function.Ident.Ident
	} else {
		sig = *a.Function.Ident.Namespace + "::" + *a.Function.Ident.Ident
	}

	return sig + " " + a.Argument.String()
}

func (state *State) AnalyseApplication(application *lexer.Application) (s *Application) {
	s = &Application{
		// TODO: change function to expression type for currying purposes in the (far) future
		Function: state.AnalyseFunction(application.Function),
		Argument: state.AnalyseExpression(application.Parameter),
	}

	if !s.Argument.TypeOf.Equals(s.Function.Takes) {
		util.Error("Trying to call",
			util.Yellow(s.Function.SigString()), "with",
			util.Yellow(s.Argument.TypeOf.String())+".\n"+
				application.Parameter.Pos.String()).Exit()
	}

	s.TypeOf = s.Function.Gives

	switch *s.Function.Ident.Ident {
	case "Return":
		if !s.Argument.TypeOf.Equals(state.CurrentFunction.Gives) {
			util.Error("Value of type",
				s.Argument.TypeOf.String(),
				"does not match function's return type:",
				state.CurrentFunction.Gives.String()).Exit()
		}
	case "Map":
		if s.Argument.TypeOf.Tuple == nil ||
			s.Argument.Atom.Tuple == nil || s.Argument.Atom.Tuple[0].Atom == nil ||
			s.Argument.Atom.Tuple[0].Atom.Function == nil ||
			s.Argument.Atom.Tuple[1].Atom.Vector == nil {
			util.Error("Malformed call to " + util.Yellow("Map") + ".\n" + application.Pos.String()).Exit()
		}

		s.TypeOf = s.Argument.Atom.Tuple[0].Atom.Function.Gives.AsVector()
		s.Function.Gives = s.TypeOf
	case "GEP":
		s.TypeOf = s.Argument.Atom.Tuple[0].TypeOf.Vector
		s.Function.Gives = s.TypeOf
	case "Foldl":
		if s.Argument.TypeOf.Tuple == nil ||
			s.Argument.Atom.Tuple == nil || s.Argument.Atom.Tuple[0].Atom == nil ||
			s.Argument.Atom.Tuple[0].Atom.Function == nil ||
			s.Argument.Atom.Tuple[0].Atom.Function.Takes.Tuple == nil ||
			s.Argument.Atom.Tuple[1].Atom == nil ||
			s.Argument.Atom.Tuple[2].Atom.Vector == nil {
			util.Error("Malformed call to " + util.Yellow("Foldl") + ".\n" + application.Pos.String()).Exit()
		}

		fn_t := s.Argument.Atom.Tuple[0].Atom.Function.Takes.Tuple
		fn_g := s.Argument.Atom.Tuple[0].Atom.Function.Gives
		id_i := s.Argument.Atom.Tuple[1].Atom.TypeOf
		vect := s.Argument.Atom.Tuple[2].Atom.TypeOf

		if !fn_t[0].Equals(id_i) {
			util.Error("Mapping function cannot accept identity in " + util.Yellow("Foldl") + ".\n" + application.Pos.String()).Exit()
		}

		if !fn_t[0].Equals(fn_g) {
			util.Error("Mapping function does not return the identity type " + util.Yellow("Foldl") + ".\n" + application.Pos.String()).Exit()
		}

		if !fn_t[1].Equals(vect) {
			util.Error("Mapping function cannot accept vector element type in " + util.Yellow("Foldl") + ".\n" + application.Pos.String()).Exit()
		}

		s.TypeOf = fn_g
		s.Function.Gives = s.TypeOf
	case "Sum":
		s.TypeOf = s.Argument.Atom.TypeOf.Vector
		s.Function.Gives = s.TypeOf
	case "Equals":
		s.TypeOf = AtomicType("Bool")
		s.Function.Gives = s.TypeOf
	case "Product":
		//s.TypeOf = s.Argument.Atom.TypeOf.Vector
		s.TypeOf = &IntType
		s.Function.Gives = s.TypeOf
	case "Append":
		s.TypeOf = s.Argument.Atom.Tuple[0].TypeOf
		s.Function.Gives = s.TypeOf
	case "First":
		s.TypeOf = s.Argument.Atom.Tuple[0].TypeOf
		s.Function.Gives = s.TypeOf
	case "Second":
		s.TypeOf = s.Argument.Atom.Tuple[1].TypeOf
		s.Function.Gives = s.TypeOf
	case "Third":
		s.TypeOf = s.Argument.Atom.Tuple[2].TypeOf
		s.Function.Gives = s.TypeOf
	}

	return s
}
