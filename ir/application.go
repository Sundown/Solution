package ir

import "sundown/sunday/parser"

type Application struct {
	TypeOf   *Type
	Function *Function
	Argument *Expression
}

func (a *Application) String() string {
	if a.Function.Ident.Namespace != nil {
		return *a.Function.Ident.Namespace + "::" + *a.Function.Ident.Ident + " " + a.Argument.String()
	} else {
		return *a.Function.Ident.Ident + " " + a.Argument.String()
	}
}

func AnalyseApplication(application *parser.Application) (s *Application) {
	s = &Application{
		Function: AnalyseFunction(application.Function),
		Argument: AnalyseExpression(application.Parameter),
	}

	s.TypeOf = s.Function.Gives

	return s
}
