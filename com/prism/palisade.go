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

func (env Environment) GetDFunction(i Ident) *DFunction {
	if f, ok := env.DFunctions[i]; ok {
		return f
	}

	return nil
}

func (env Environment) GetMFunction(i Ident) *MFunction {
	if f, ok := env.MFunctions[i]; ok {
		return f
	}

	return nil
}
