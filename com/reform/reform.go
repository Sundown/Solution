package reform

import (
	"sundown/solution/lexer"

	"github.com/alecthomas/repr"
)

type State struct{}

func (state *State) Init(lex *lexer.State) {
	for _, s := range lex.Expressions {
		for _, v := range s.Singletons {
			repr.Println(v)
		}
	}
}
