package weave

import (
	"sundown/solution/palisade"
	"sundown/solution/prism"
)

type State struct {
	TypeDefinitions map[prism.Ident]prism.Type
}

// Dangle type is for LHS arguments to dyadic operators
// which are left unbound to the RHS until a later analysis pass.

type SubState struct {
	Env            *prism.Environment
	mutablePos     int
	mutableSubPos  int
	Subexpressions []*palisade.Subexpression
}
