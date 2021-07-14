package ir

type Function struct {
	Name  string
	Takes *Type
	Gives *Type
	Body  *Expression
}

func (f *Function) String() string {
	return f.Name + "<" + f.Takes.String() + " -> " + f.Gives.String() + ">\n" + f.Body.String()
}
