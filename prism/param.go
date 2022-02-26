package prism

// Type property for interface
func (a Alpha) Type() Type {
	return a.TypeOf
}

// Type property for interface
func (a Omega) Type() Type {
	return a.TypeOf
}

// String function for interface
func (a Alpha) String() string {
	return "α"
}

// String function for interface
func (o Omega) String() string {
	return "ω"
}

type Alpha struct {
	TypeOf Type
}
type Omega struct {
	TypeOf Type
}
