package analysis

type Function struct {
	Name  string
	Takes *Type
	Gives *Type
	Body  *Expression
}
