package parse

import (
	"sundown/solution/lex"
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

func (state *State) AnalyseApplication(application *lex.Application) (s *Application) {
	s = &Application{
		// TODO: change function to expression type for currying purposes in the (far) future
		Function: state.AnalyseFunction(application.Function),
		Argument: state.AnalyseExpression(application.Parameter),
	}

	if !s.Argument.TypeOf.Equals(s.Function.Takes) {
		util.Error("Trying to call", s.Function.SigString(), "with", s.Argument.TypeOf.String()).Exit()
	}

	s.TypeOf = s.Function.Gives

	if *s.Function.Ident.Ident == "Return" {
		if !s.Argument.TypeOf.Equals(state.CurrentFunction.Gives) {
			util.Error("Value of type",
				s.Argument.TypeOf.String(),
				"does not match function's return type:",
				state.CurrentFunction.Gives.String()).Exit()
		}
	}

	return s
}
