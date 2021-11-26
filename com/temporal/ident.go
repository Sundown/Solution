package temporal

import (
	"sundown/solution/lexer"
)

type Ident struct {
	Namespace *string
	Ident     *string
}

type IdentKey struct {
	Namespace string
	Ident     string
}

func (i *Ident) String() string {
	if i.Namespace == nil {
		return *i.Ident
	} else {
		return *i.Namespace + "::" + *i.Ident
	}
}

// Dereferences strings so they AREN'T unique, used for key within maps
func (i *Ident) AsKey() IdentKey {
	var n string
	if i.Namespace == nil {
		n = "_"
	} else {
		n = *i.Namespace
	}

	if i.Ident == nil {
		panic("Unreachable")
	}

	return IdentKey{
		Namespace: n,
		Ident:     *i.Ident,
	}
}

// Is this ident referring to Foundation?
func (i *Ident) IsFoundational() bool {
	return *i.Namespace == "_" || *i.Namespace == "foundation" || *i.Namespace == "se"
}

// Transform lexer identifier to temporal identifier
func IRIdent(i *lexer.Ident) *Ident {
	return &Ident{
		Namespace: i.Namespace,
		Ident:     i.Ident,
	}
}
