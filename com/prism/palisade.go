package prism

import (
	"sundown/solution/palisade"

	"github.com/alecthomas/participle/v2"
)

func Intern(i palisade.Ident) (p Ident) {
	if i.Namespace == nil {
		p.Package = "_"
	} else {
		p.Package = *i.Namespace
	}

	p.Name = *i.Ident
	return
}

func ParseIdent(s string) (p Ident) {
	var t palisade.Ident
	err := participle.MustBuild(
		&palisade.Ident{},
		participle.UseLookahead(4),
		participle.Unquote()).
		ParseString("", s, &t)

	if err != nil {
		panic(err)
	}

	return Intern(t)
}

func (env Environment) GetDyadicFunction(i Ident) *DyadicFunction {
	if f, ok := env.DyadicFunctions[i]; ok {
		return f
	}

	return nil
}

func (env Environment) GetMonadicFunction(i Ident) *MonadicFunction {
	if f, ok := env.MonadicFunctions[i]; ok {
		return f
	}

	return nil
}
