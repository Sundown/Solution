package ir

import "sundown/sunday/parser"

type Application struct {
	TypeOf   *Type
	Function *Function
	Argument *Expression
}

func (a *Application) String() string {
	return a.Function.Ident.Namespace + "::" + a.Function.Ident.Ident + " " + a.Argument.String()
}

func (state *State) AnalyseApplication(application *parser.Application) (s *Application) {
	s = &Application{
		Function: state.AnalyseFunction(application.Function),
		Argument: state.AnalyseExpression(application.Parameter),
	}

	s.TypeOf = s.Function.Gives

	return s
}
