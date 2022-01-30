package prism

import (
	"sundown/solution/palisade"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer/stateful"
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

var basicLexer = stateful.MustSimple([]stateful.Rule{
	{"whitespace", `[ \s]+`, nil}, // THIS IS LOWERCASE FOR A REASON
	{"EOL", `[\n\r]+`, nil},
	{"String", `"(\\"|[^"])*"`, nil},
	{"Int", `\d+`, nil},
	{"Float", `(\d*\.)?\d+`, nil},
	{"Ident", `[\w/+]+`, nil},
	{"Char", `\'.\'`, nil},
	{"Punct", `[-[!@#$%^&*()+_={}\|:;"'<,>.?Δ∇→]|]`, nil},
	{"Alpha", "α", nil},
	{"Omega", "ω", nil},
})

func ParseIdent(s string) (p Ident) {
	var t palisade.Ident
	err := participle.MustBuild(
		&palisade.Ident{},
		participle.Lexer(basicLexer),
		participle.UseLookahead(4),
		//participle.Elide("Whitespace"),
		//	participle.Elide("EOL"),
		participle.Unquote("String")).
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
