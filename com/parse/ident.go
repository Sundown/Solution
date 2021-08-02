package parse

import (
	"sundown/solution/lex"
)

type Ident struct {
	Namespace *string
	Ident     *string
}

type IdentKey struct {
	Namespace string
	Ident     string
}

// Dereferences strings so they AREN'T unique, used for key within maps
func (i *Ident) AsKey() IdentKey {
	var n string
	if i.Namespace == nil {
		n = "_"
	} else {
		n = *i.Namespace
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

// Transform lex identifier to parse identifier
func IRIdent(i *lex.Ident) *Ident {
	return &Ident{
		Namespace: i.Namespace,
		Ident:     i.Ident,
	}
}
