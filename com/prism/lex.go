package prism

import (
	"os"
	"sundown/solution/palisade"

	"github.com/alecthomas/participle/v2"
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
		participle.UseLookahead(4),
		participle.Unquote(),
	).Parse(env.File, r, &res)

	if err != nil {
		panic(err)
	}

	env.LexResult = &res

	return env
}
