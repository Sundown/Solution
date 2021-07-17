package ir

import (
	"sundown/sunday/parser"
)

type Ident struct {
	Namespace *string
	Ident     *string
}

type IdentKey struct {
	Namespace string
	Ident     string
}

func (i *Ident) AsKey() IdentKey {
	return IdentKey{
		Namespace: *i.Namespace,
		Ident:     *i.Ident,
	}
}

func (i *Ident) IsFoundational() bool {
	return *i.Namespace == "_" || *i.Namespace == "foundational" || *i.Namespace == "se"
}

func IRIdent(i *parser.Ident) *Ident {
	return &Ident{
		Namespace: i.Namespace,
		Ident:     i.Ident,
	}
}
