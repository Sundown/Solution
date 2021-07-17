package lex

import (
	"os"

	"github.com/alecthomas/participle/v2"
)

var Parser = participle.MustBuild(&State{}, participle.UseLookahead(4), participle.Unquote())

func (prog *State) Lex(args []string) *State {
	file, err := os.Open(args[1])
	if err != nil {
		panic(err)
	}

	defer file.Close()

	err = Parser.Parse(args[1], file, prog)

	if err != nil {
		panic(err)
	}

	return prog
}
