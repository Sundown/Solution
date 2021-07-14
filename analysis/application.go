package analysis

import "sundown/sunday/parser"

type Application struct {
	TypeOf   *Type
	Function *Function
	Argument *Expression
}

func AnalyseApplication(application *parser.Application) (s *Application) {
	s = &Application{
		Function: AnalyseFunction(application.Function),
		Argument: AnalyseExpression(application.Parameter),
	}

	s.TypeOf = s.Function.Gives

	return s
}

func AnalyseBlock(block []*parser.Expression) (b *Expression) {
	var body []*Expression
	for index, elm := range block {
		body[index] = AnalyseExpression(elm)
	}

	/* TODO: need some way to calculate typeof */
	b = &Expression{Block: body}
	return b
}

func AnalyseStatement(statement *parser.FnDecl) (s *Function) {
	takes, gives := AnalyseType(statement.Takes), AnalyseType(statement.Gives)
	e := Expression{TypeOf: gives}
	for _, expr := range statement.Expressions {
		e.Block = append(e.Block, AnalyseExpression(expr))
	}

	return &Function{
		Name:  *statement.Ident,
		Takes: takes,
		Gives: gives,
		Body:  &e,
	}
}

func AnalyseFunction(function *string) (f *Function) {
	f = &Function{
		Name:  *function,
		Takes: &Type{Atomic: "Int"},
		Gives: &Type{Atomic: "Int"}}
	return f
}
