package ir

type Function struct {
	Ident *Ident
	Takes *Type
	Gives *Type
	Body  *Expression
}

func (f *Function) String() string {
	if f.Ident.Namespace != nil {
		return *f.Ident.Namespace + "::" + *f.Ident.Ident
	} else {
		return *f.Ident.Ident
	}
}
