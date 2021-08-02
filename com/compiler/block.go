package compiler

import "sundown/solution/parse"

func (state *State) CompileBlock(body *parse.Expression) {
	// Block is just an expression[]
	for _, stmt := range body.Block {
		state.CompileExpression(stmt)
	}
}
