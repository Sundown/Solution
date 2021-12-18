package subtle

import (
	"sundown/solution/palisade"
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

	s.TypeOf = s.Function.Gives

	return s
}
