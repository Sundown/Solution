package reform

import "fmt"

func (b Application) String() string {
	return b.Operator + " (" + b.Operand.String() + ")"
}

func (u Dangle) String() string {
	return "(" + u.Outer.String() + " " + u.Inner.String() + ")"
}

func (i Int) String() string {
	return fmt.Sprintf("%d", i.Value)
}

func (i Ident) String() string {
	return i.Value
}

func (e EOF) String() string {
	return "EOF"
}
