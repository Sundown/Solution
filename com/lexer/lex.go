package lexer

import (
	"os"
	"sundown/solution/util"

	"github.com/alecthomas/participle/v2"
)

var Parser = participle.MustBuild(&State{}, participle.UseLookahead(4), participle.Unquote())

func (prog *State) Lex(rt *util.Runtime) *State {
	util.Verbose("Init lexer")
	r, err := os.Open(rt.File)
	defer r.Close()

	if err != nil {
		util.Error(err.Error()).Exit()
	}

	err = Parser.Parse(rt.File, r, prog)

	if err != nil {
		panic(err)
	}

	return prog
}
