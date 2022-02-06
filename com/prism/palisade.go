package prism

import (
	"os"
	"sundown/solution/palisade"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer/stateful"
)

func Lex(env *Environment) *Environment {
	Verbose("Init palisade")
	res := palisade.PalisadeResult{}
	r, err := os.Open(env.File)
	defer r.Close()

	if err != nil {
		Error(err.Error()).Exit()
	}

	err = participle.MustBuild(
		&palisade.PalisadeResult{},
		participle.UseLookahead(40000), // vectors don't seem to work if this is low
		participle.Lexer(basicLexer),
		participle.Unquote(),
	).Parse(env.File, r, &res)

	if err != nil {
		panic(err)
	}

	env.LexResult = &res

	return env
}

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
	{Name: "whitespace", Pattern: `[ \s]+`, Action: nil}, // THIS IS LOWERCASE FOR A REASON
	{Name: "EOL", Pattern: `[\n\r]+`, Action: nil},
	{Name: "String", Pattern: `"(\\"|[^"])*"`, Action: nil},
	{Name: "Int", Pattern: `\d+`, Action: nil},
	{Name: "Float", Pattern: `(\d*\.)?\d+`, Action: nil},
	{Name: "Ident", Pattern: `([\w]+|[-*+/÷])`, Action: nil},
	{Name: "Char", Pattern: `\'.\'`, Action: nil},
	{Name: "Punct", Pattern: `[-[!@#$%^&*()+_={}\|:;"'<,>.?Δ∇→]|]`, Action: nil},
	{Name: "Alpha", Pattern: "α", Action: nil},
	{Name: "Omega", Pattern: "ω", Action: nil},
})

func ParseIdent(s string) (p Ident) {
	var t palisade.Ident
	err := participle.MustBuild(
		&palisade.Ident{},
		participle.Lexer(basicLexer),
		participle.UseLookahead(2),
		participle.Elide("whitespace"),
		participle.Elide("EOL"),
		participle.Elide("Punct"),
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
