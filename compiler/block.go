package compiler

import "sundown/sunday/parse"

func (state *State) CompileBlock(body *parse.Expression) {
	for _, stmt := range body.Block {
		state.CompileExpression(stmt)
	}
}
