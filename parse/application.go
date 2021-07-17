package parse

import "sundown/sunday/lex"

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
		// TODO: change function to expression type for currying purposes in the future
		Function: state.AnalyseFunction(application.Function),
		Argument: state.AnalyseExpression(application.Parameter),
	}

	s.TypeOf = s.Function.Gives

	return s
}
